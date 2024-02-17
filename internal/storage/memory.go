package storage

import (
	v1 "colaco-api/internal/api/v1"
	"fmt"
	"strings"
	"sync"
)

type MemoryStorage struct {
	StorageMap map[string]v1.VendingSlot
	m          sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	s := make(map[string]v1.VendingSlot)
	return &MemoryStorage{
		StorageMap: s,
		m:          sync.RWMutex{},
	}
}

func (m *MemoryStorage) GetSlot(name string) (v1.VendingSlot, bool, error) {
	m.m.RLock()
	defer m.m.RUnlock()
	if val, ok := m.StorageMap[strings.ToLower(name)]; ok {
		return val, true, nil
	}
	return v1.VendingSlot{}, false, nil
}

func (m *MemoryStorage) UpsertSlot(name string, slot v1.VendingSlot) {
	m.m.Lock()
	defer m.m.Unlock()
	m.StorageMap[strings.ToLower(name)] = slot
}

func (m *MemoryStorage) GetSlots() (slots []v1.VendingSlot) {
	m.m.RLock()
	defer m.m.RUnlock()
	for _, v := range m.StorageMap {
		slots = append(slots, v)
	}
	return slots
}

func (m *MemoryStorage) DeleteSlot(name string) (bool, error) {
	m.m.Lock()
	defer m.m.Unlock()
	delete(m.StorageMap, strings.ToLower(name))
	if _, ok := m.StorageMap[strings.ToLower(name)]; ok {
		return false, fmt.Errorf("still exists")
	}
	return true, nil
}

func (m *MemoryStorage) UpdatePrice(name string, price float32) error {
	m.m.Lock()
	defer m.m.Unlock()
	slot, ok := m.StorageMap[strings.ToLower(name)]
	if !ok {
		return fmt.Errorf("slot not found")
	}
	slot.Cost = &price
	m.StorageMap[strings.ToLower(name)] = slot
	return nil
}

func (m *MemoryStorage) UpdateQuantity(name string, qty int) error {
	m.m.Lock()
	defer m.m.Unlock()
	if val, ok := m.StorageMap[strings.ToLower(name)]; ok {
		val.Quantity = &qty
		m.StorageMap[strings.ToLower(name)] = val
	}
	return fmt.Errorf("slot does not exist")
}

func (m *MemoryStorage) AddSlot(name string, slot v1.VendingSlot) {
	m.m.Lock()
	defer m.m.Unlock()
	m.StorageMap[strings.ToLower(name)] = slot
}
