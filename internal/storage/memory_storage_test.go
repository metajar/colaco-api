package storage

import (
	v1 "colaco-api/internal/api/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	ms := NewMemoryStorage()
	assert.NotNil(t, ms)
	assert.Empty(t, ms.StorageMap)
}

func TestUpsertAndGetSlot(t *testing.T) {
	ms := NewMemoryStorage()
	slotName := "coke"
	slot := v1.VendingSlot{Cost: new(float32), Quantity: new(int)}
	*slot.Cost = 1.25
	*slot.Quantity = 20

	ms.UpsertSlot(slotName, slot)

	retSlot, found, err := ms.GetSlot(slotName)
	assert.True(t, found)
	assert.Nil(t, err)
	assert.Equal(t, slot, retSlot)
}

func TestGetSlots(t *testing.T) {
	ms := NewMemoryStorage()
	ms.UpsertSlot("coke", v1.VendingSlot{Cost: new(float32), Quantity: new(int)})
	ms.UpsertSlot("pepsi", v1.VendingSlot{Cost: new(float32), Quantity: new(int)})

	slots := ms.GetSlots()
	assert.Len(t, slots, 2)
}

func TestDeleteSlot(t *testing.T) {
	ms := NewMemoryStorage()
	slotName := "coke"
	ms.UpsertSlot(slotName, v1.VendingSlot{Cost: new(float32), Quantity: new(int)})

	success, err := ms.DeleteSlot(slotName)
	assert.True(t, success)
	assert.Nil(t, err)

	_, found, _ := ms.GetSlot(slotName)
	assert.False(t, found)
}

func TestUpdatePrice(t *testing.T) {
	ms := NewMemoryStorage()
	slotName := "coke"
	initialPrice := float32(1.0)
	ms.UpsertSlot(slotName, v1.VendingSlot{Cost: &initialPrice, Quantity: new(int)})

	newPrice := float32(1.5)
	err := ms.UpdatePrice(slotName, newPrice)
	assert.Nil(t, err)

	slot, found, _ := ms.GetSlot(slotName)
	assert.True(t, found)
	assert.Equal(t, newPrice, *slot.Cost)
}
