package server

import (
	"bytes"
	"colaco-api/internal/api/v1"
	"colaco-api/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func s2p(s string) *string {
	return &s
}

func i2p(i int) *int {
	return &i
}

func f322p(f float32) *float32 {
	return &f
}

func TestAuthLoginSuccess(t *testing.T) {
	e := echo.New()
	reqBody := `{"username":"admin","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	vm := NewVendingMachine()

	if assert.NoError(t, vm.AuthLogin(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestPostPurchaseSuccess(t *testing.T) {
	e := echo.New()
	reqBody := `{"name":"Coke","payment":2.00}`
	req := httptest.NewRequest(http.MethodPost, "/purchase", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	vm := NewVendingMachine(
		WithStorage(storage.NewMemoryStorage()),
		WithStartingSodas([]v1.VendingSlot{
			{
				OccupiedSoda: &v1.Soda{Name: s2ptr("Coke")},
				Cost:         f322p(1.5),
				Quantity:     i2p(10),
			},
		}))

	if assert.NoError(t, vm.PostPurchase(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		body := rec.Body.Bytes()
		var p v1.PurchaseSodaResponse

		if err := json.Unmarshal(body, &p); assert.NoError(t, err) {
			expectedChange := float32(2.0) - 1.5
			assert.Equal(t, expectedChange, *p.Change, "The change returned is not correct")
		}
	}
}

func TestRestockSodaSuccess(t *testing.T) {
	e := echo.New()
	reqBody := `{"name":"Coke","quantity":5}`
	req := httptest.NewRequest(http.MethodPost, "/restock", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	vm := NewVendingMachine(
		WithStorage(storage.NewMemoryStorage()),
		WithStartingSodas([]v1.VendingSlot{
			{
				OccupiedSoda: &v1.Soda{Name: s2ptr("Coke")},
				Quantity:     i2p(5),
				MaxQuantity:  i2p(15),
			},
		}))

	if assert.NoError(t, vm.RestockSoda(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestVendingMachine_UpdatePrice(t *testing.T) {
	e := echo.New()

	// Mock setup or in-memory setup for SlotStorage
	mockStorage := storage.NewMemoryStorage()         // Assuming you have a function to create a mock or in-memory storage
	vm := NewVendingMachine(WithStorage(mockStorage)) // Assuming NewVendingMachine accepts a SlotStorage parameter

	// Initial setup with a soda
	sodaName := "Coke"
	initialPrice := float32(1.0)
	mockStorage.UpsertSlot(strings.ToLower(sodaName), v1.VendingSlot{
		OccupiedSoda: &v1.Soda{Name: &sodaName},
		Cost:         &initialPrice,
	})

	// Preparing the request to update the price
	newPrice := float32(1.5)
	updatePriceBody := v1.UpdatePriceBody{
		Name:     sodaName,
		NewPrice: newPrice,
	}
	bodyBytes, _ := json.Marshal(updatePriceBody) // Error handling omitted for brevity

	req := httptest.NewRequest(http.MethodPost, "/update-price", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := vm.UpdatePrice(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	slot, found, _ := mockStorage.GetSlot(strings.ToLower(sodaName)) // Adjust based on your actual mock/storage implementation
	if assert.True(t, found) {
		assert.Equal(t, newPrice, *slot.Cost)
	}
}

func TestDeleteVendingSuccess(t *testing.T) {
	e := echo.New()
	mockStorage := storage.NewMemoryStorage()
	vm := NewVendingMachine(WithStorage(mockStorage))
	sodaName := "Coke"
	vm.SlotStorage.UpsertSlot(sodaName, v1.VendingSlot{
		OccupiedSoda: &v1.Soda{Name: &sodaName},
	})
	_, existsBeforeDelete, _ := vm.SlotStorage.GetSlot(strings.ToLower(sodaName))
	assert.True(t, existsBeforeDelete, "Soda should exist before deletion")

	deleteRequestBody := v1.VendingSlotRequestBody{Name: sodaName}
	bodyBytes, err := json.Marshal(deleteRequestBody)
	if err != nil {
		t.Fatalf("Marshalling request body failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/delete-vending", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, vm.DeleteVending(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		_, existsAfterDelete, _ := vm.SlotStorage.GetSlot(strings.ToLower(sodaName))
		assert.False(t, existsAfterDelete, "Soda should not exist after deletion")
	}
}

func TestGetVendingSuccess(t *testing.T) {
	e := echo.New()

	// Setup: Initialize the vending machine with some sodas
	calories := 150
	description := "Classic Coke"
	name := "Coke"
	originStory := "Invented in the 19th century"
	ounces := float32(12.0)
	cost := float32(1.25)
	maxQuantity := 20
	quantity := 10

	mockStorage := storage.NewMemoryStorage()
	vm := NewVendingMachine(WithStorage(mockStorage))
	vm.SlotStorage.UpsertSlot("coke", v1.VendingSlot{
		Cost:        &cost,
		MaxQuantity: &maxQuantity,
		OccupiedSoda: &v1.Soda{
			Calories:    &calories,
			Description: &description,
			Name:        &name,
			OriginStory: &originStory,
			Ounces:      &ounces,
		},
		Quantity: &quantity,
	})

	req := httptest.NewRequest(http.MethodGet, "/vending", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute the GetVending endpoint
	if assert.NoError(t, vm.GetVending(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response body to check if it contains the expected items
		var vendingSlots []v1.VendingSlot
		if err := json.Unmarshal(rec.Body.Bytes(), &vendingSlots); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		// Verify the response contains the correct number of items
		assert.Len(t, vendingSlots, 1, "Expected 1 soda in the response")

		// Further detailed checks can be performed on the contents of vendingSlots
		// Example: Check the first item's name and calories
		if len(vendingSlots) > 0 {
			assert.Equal(t, name, *vendingSlots[0].OccupiedSoda.Name)
			assert.Equal(t, calories, *vendingSlots[0].OccupiedSoda.Calories)
			// Add more assertions as needed
		}
	}
}

func TestPostNewSodaSuccess(t *testing.T) {
	e := echo.New()

	// Setup: Prepare a new soda to add
	newSoda := v1.Soda{
		Calories:    new(int),
		Description: new(string),
		Name:        new(string),
		OriginStory: new(string),
		Ounces:      new(float32),
	}
	*newSoda.Calories = 200
	*newSoda.Description = "A refreshing lemon-lime soda."
	*newSoda.Name = "Sprite"
	*newSoda.OriginStory = "Invented to compete with similar citrus sodas."
	*newSoda.Ounces = 12.0

	newSlot := v1.VendingSlot{
		Cost:         new(float32),
		MaxQuantity:  new(int),
		OccupiedSoda: &newSoda,
		Quantity:     new(int),
	}
	*newSlot.Cost = 1.5
	*newSlot.MaxQuantity = 20
	*newSlot.Quantity = 15

	reqBodyBytes, err := json.Marshal(newSlot)
	if err != nil {
		t.Fatalf("Failed to marshal new soda request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	vm := NewVendingMachine(WithStorage(storage.NewMemoryStorage()))

	// Execute the PostNew endpoint
	if assert.NoError(t, vm.PostNew(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Verify the new soda has been added to the vending machine
		addedSoda, exists, _ := vm.SlotStorage.GetSlot(strings.ToLower(*newSoda.Name))
		if assert.True(t, exists, "New soda should exist in the vending machine after addition") {
			// Perform detailed checks on the added soda
			assert.Equal(t, *newSoda.Name, *addedSoda.OccupiedSoda.Name)
			assert.Equal(t, *newSlot.Cost, *addedSoda.Cost)
			// Add more assertions as needed to verify the added soda details
		}
	}
}
