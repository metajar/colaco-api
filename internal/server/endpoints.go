package server

import (
	"colaco-api/internal/api/v1"
	"colaco-api/internal/jwt"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"net/http"
)

// AuthLogin handles the authentication and login process for the vending
// machine. It first binds the request body to an AuthRequestBody struct. If the
// request is invalid, it returns a JSON response with a "Invalid request" error.
// Next, it authenticates the user using the authenticateUser function. If the
// username and password are invalid, it returns a JSON response with a "Invalid
// username and/or password
func (v *VendingMachine) AuthLogin(ctx echo.Context) error {
	var loginReq v1.AuthRequestBody
	if err := ctx.Bind(&loginReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, genErrorResponse("Invalid request"))
	}
	if !authenticateUser(loginReq.Username, loginReq.Password) {
		return ctx.JSON(http.StatusUnauthorized, genErrorResponse("Invalid username and/or password"))
	}
	authenticator, err := jwt.NewFakeAuthenticator()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, genErrorResponse("Failed to initialize authenticator"))
	}
	tokenBytes, err := authenticator.CreateJWSWithClaims([]string{"user"})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, genErrorResponse("Failed to sign token"))
	}
	// Return the signed token in an authtokenresponse.
	return ctx.JSON(http.StatusOK, v1.AuthTokenResponse{Token: s2ptr(string(tokenBytes))})
}

// PostPurchase handles the process of purchasing a soda from the vending machine.
// It first binds the request body to a PurchaseSodaBody struct. If the request is invalid,
// it returns a JSON response with an error message.
// Then it locks the vending machine to perform the purchase transaction.
// It converts the soda name to lowercase and checks if it exists in the vending machine's Slots.
// If the soda does not exist, it returns a JSON response with an error message.
// If the payment is sufficient to purchase the soda, it calculates the change,
// decrements the quantity of the soda, updates the vending machine's Slots, and returns
// a JSON response with the change amount and the purchased soda.
// If the payment is insufficient, it returns a JSON response with an error message.
func (v *VendingMachine) PostPurchase(ctx echo.Context) error {
	var purchase v1.PurchaseSodaBody
	if err := ctx.Bind(&purchase); err != nil {
		return ctx.JSON(500, genErrorResponse(err.Error()))
	}
	v.m.Lock()
	defer v.m.Unlock()
	var vslot v1.VendingSlot
	val, found, _ := v.SlotStorage.GetSlot(purchase.Name)
	if !found {
		return ctx.JSON(
			404,
			genErrorResponse(fmt.Sprintf("soda with name %v does not exist", purchase.Name)))
	} else {
		vslot = val
	}
	if purchase.Payment >= *vslot.Cost {
		purchaseDecimal := decimal.NewFromFloat32(purchase.Payment)
		costDecimal := decimal.NewFromFloat32(*vslot.Cost)
		change := purchaseDecimal.Sub(costDecimal)
		*vslot.Quantity--
		v.SlotStorage.UpsertSlot(purchase.Name, vslot)
		f, _ := change.Float64()
		c := float32(f)
		return ctx.JSON(200, v1.PurchaseSodaResponse{
			Change: &c,
			Soda:   vslot.OccupiedSoda,
		})

	}
	mess := fmt.Sprintf("insufficient funds. soda costs %v and you only provided %v", *vslot.Cost, purchase.Payment)
	return ctx.JSON(402, genMessageResponse(mess))

}

// RestockSoda restocks the quantity of a specified soda in the vending machine.
// It first binds the request body to a RestockRequestBody struct. If the request
// is invalid, it returns a JSON response with an error message. Next, it locks
// the vending machine to ensure thread safety while accessing and modifying the
// Slots map. It then checks if the specified soda exists in the Slots map. If it
// doesn't, it returns a JSON response with "slot '{soda name}' not found" error.
// If the soda exists, it retrieves the VendingSlot object from the map. If the
// sum of the requested quantity and the current quantity is greater than the
// maximum allowed quantity, it calculates the leftover quantity and sets the
// vending slot quantity to the maximum allowed quantity, then updates the Slots
// map. Finally, it returns a JSON response with the updated RestockResponse,
// including the leftover quantity, new quantity, and old quantity. If the sum of
// the requested quantity and the current quantity is within the maximum allowed
// quantity, it simply adds the requested quantity to the current quantity,
// updates the Slots map, and returns a JSON response with the updated
// RestockResponse, including the leftover quantity (which is 0 in this case),
// new quantity
func (v *VendingMachine) RestockSoda(ctx echo.Context) error {
	var m v1.RestockRequestBody
	if err := ctx.Bind(&m); err != nil {
		return ctx.JSON(500, genErrorResponse(err.Error()))
	}
	v.m.Lock()
	defer v.m.Unlock()
	var vendSlot v1.VendingSlot
	val, found, _ := v.SlotStorage.GetSlot(m.Name)
	if !found {
		return ctx.JSON(404,
			genErrorResponse(fmt.Sprintf("slot '%v' not found", m.Name)))
	} else {
		vendSlot = val
	}
	var leftover int
	oldQty := *vendSlot.Quantity
	if m.Quantity+*vendSlot.Quantity > *vendSlot.MaxQuantity {
		leftover = (m.Quantity + *vendSlot.Quantity) - *vendSlot.MaxQuantity
		vendSlot.Quantity = vendSlot.MaxQuantity
		v.SlotStorage.UpsertSlot(m.Name, vendSlot)
		return ctx.JSON(200, v1.RestockResponse{
			Leftover:    &leftover,
			NewQuantity: vendSlot.MaxQuantity,
			OldQuantity: &oldQty,
		})
	}
	*vendSlot.Quantity += m.Quantity
	v.SlotStorage.UpsertSlot(m.Name, vendSlot)
	return ctx.JSON(200, v1.RestockResponse{
		Leftover:    &leftover,
		NewQuantity: vendSlot.Quantity,
		OldQuantity: &oldQty,
	})
}

// UpdatePrice updates the price of a soda in the vending machine. It first binds
// the request body to an UpdatePriceBody struct. If the binding fails, it returns
// a JSON response with an error message. It then locks the vending machine for
// writing, retrieves the slot for the specified soda name, and updates its price
// to the new price provided in the request body. If the slot does not exist, it
// returns a JSON response with an error message. Finally, it responds with a JSON
// response indicating the success of the operation and the updated soda price.
func (v *VendingMachine) UpdatePrice(ctx echo.Context) error {
	var m v1.UpdatePriceBody
	if err := ctx.Bind(&m); err != nil {
		return ctx.JSON(500, genErrorResponse(err.Error()))
	}
	// Locking the vending machine for writing
	v.m.Lock()
	defer v.m.Unlock()

	// Check if the slot exists
	slot, found, _ := v.SlotStorage.GetSlot(m.Name)
	if !found {
		return ctx.JSON(
			404,
			genErrorResponse(fmt.Sprintf("slot '%v' not found", m.Name)))
	}
	// Update the price in the slot
	old := slot.Cost
	slot.Cost = &m.NewPrice
	v.SlotStorage.UpsertSlot(m.Name, slot)

	// Respond with success
	return ctx.JSON(http.StatusOK, v1.UpdatePriceResp{
		NewPrice: &m.NewPrice,
		OldPrice: old,
		SlotName: &m.Name,
	})

}

// DeleteVending deletes a vending slot from the vending machine based on the
// provided name. It first binds the request body to a VendingSlotRequestBody
// struct. If the binding fails, it returns a JSON response with an error
// message.
//
// After binding the request body, it acquires a write lock on the vending
// machine to ensure exclusive access during deletion. It then converts the slot
// name to lowercase for case-insensitive lookup in the Slots map. If the slot
// with the provided name exists, it is deleted from the Slots map. It returns a
// JSON response with a success message if the deletion is successful. If the
// slot does not exist, it returns a JSON response with an error message.
func (v *VendingMachine) DeleteVending(ctx echo.Context) error {
	var m v1.DeleteVendingJSONBody
	if err := ctx.Bind(&m); err != nil {
		return ctx.JSON(500, genErrorResponse(err.Error()))
	}
	v.m.Lock()
	defer v.m.Unlock()
	_, ok, _ := v.SlotStorage.GetSlot(m.Name)
	if ok {
		isDeleted, err := v.SlotStorage.DeleteSlot(m.Name)
		if err != nil {
			return ctx.JSON(500, genErrorResponse(err.Error()))
		}
		if !isDeleted {
			return ctx.JSON(500, genErrorResponse(fmt.Sprintf("%v is not deleted", m.Name)))
		}
		return ctx.JSON(200, genMessageResponse(fmt.Sprintf("soda '%v' deleted successfully", m.Name)))
	}
	return ctx.JSON(404, genMessageResponse(fmt.Sprintf("soda '%v' not found", m.Name)))
}

// GetVending retrieves all the vending slots available in the vending machine.
// It iterates over the slots and appends each slot to the vendingSlots slice.
// If the vendingSlots slice is empty, it returns a JSON response with
// an error message indicating that the vending machine is empty.
// Otherwise, it returns a JSON response with the vendingSlots slice.
func (v *VendingMachine) GetVending(ctx echo.Context) error {
	var vendingSlots []v1.VendingSlot
	vendingSlots = v.SlotStorage.GetSlots()
	count := len(vendingSlots)
	if vendingSlots == nil {
		return ctx.JSON(404, map[string]string{"error": "vending machine is empty"})
	}
	return ctx.JSON(200, v1.VendingMachineResponse{
		Slots: &vendingSlots,
		Total: &count,
	})
}

// PostNew handles the creation of a new vending slot for a soda in the vending machine.
// It first binds the request body to a VendingSlot struct. If the request is invalid,
// it returns a JSON response with an "unacceptable soda" error.
// Next, it checks if a slot with the same soda name already exists in the vending machine.
// If it does, it returns a JSON response with a "soda already exists" error.
// If the slot is unique, it adds the slot to the vending machine and returns a JSON
// response with a success message.
func (v *VendingMachine) PostNew(ctx echo.Context) error {
	var VSlot v1.PostNewJSONRequestBody
	if err := ctx.Bind(&VSlot); err != nil {
		return ctx.JSON(406, genErrorResponse("unacceptable soda"))
	}
	_, found, _ := v.SlotStorage.GetSlot(*VSlot.Slot.OccupiedSoda.Name)

	if found {
		return ctx.JSON(
			409,
			genErrorResponse(fmt.Sprintf("soda already exists for: '%v'", *VSlot.Slot.OccupiedSoda.Name)))
	}
	v.SlotStorage.AddSlot(*VSlot.Slot.OccupiedSoda.Name, VSlot.Slot)
	return ctx.JSON(
		201,
		genMessageResponse(fmt.Sprintf("soda created for: '%v'", *VSlot.Slot.OccupiedSoda.Name)))
}
