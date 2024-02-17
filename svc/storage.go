package svc

import v1 "colaco-api/internal/api/v1"

type VendingStorageInterface interface {
	GetSlot(name string) (v1.VendingSlot, bool, error)
	UpsertSlot(name string, slot v1.VendingSlot)
	GetSlots() []v1.VendingSlot
	DeleteSlot(name string) (bool, error)
	UpdatePrice(name string, price float32) error
	UpdateQuantity(name string, qty int) error
	AddSlot(name string, slot v1.VendingSlot)
}
