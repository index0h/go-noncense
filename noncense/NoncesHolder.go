package noncense

import (
	"errors"
)

// Structure of hashmap item
type mapNode struct {
	prev  *mapNode
	next  *mapNode
	value *HString
}

// Structure of ring linked list item
type listNode struct {
	next    *listNode
	mapItem *mapNode
}

// Structure for nonces holder (hashmap + linked list)
type NoncesHolder struct {
	count    uint32
	served   uint32
	sizeMap  uint32
	sizeList uint32

	hashMap []*mapNode
	load    []uint32

	first *listNode
	last  *listNode
}

// Constructor
func NewNoncesHolder(sizeMap uint32, sizeList uint32) *NoncesHolder {
	return &NoncesHolder{
		count:    0,
		sizeList: sizeList,
		sizeMap:  sizeMap,
		hashMap:  make([]*mapNode, sizeMap),
		load:     make([]uint32, sizeMap),
	}
}

// Returns current max load
func (h *NoncesHolder) GetLoad() uint32 {

	if h.count == 0 {
		return 0
	} else if h.count == 1 {
		return 1
	}

	// Calculating max load
	var max uint32
	var i uint32
	for i = 0; i < h.sizeMap; i++ {
		if h.load[i] > max {
			max = h.load[i]
		}
	}

	return max
}

// Returns count of served NONCEs
func (h *NoncesHolder) GetServedCount() uint32 {
	return h.served
}

// Adds new NONCE
func (h *NoncesHolder) Add(value HString) error {

	if h.Has(value) {
		return errors.New("Value already presents")
	}

	// Building map node
	mn := mapNode{value: &value}

	// Placing map node into hashmap
	i := value.trim(h.sizeMap)
	h.load[i]++
	if h.hashMap[i] == nil {
		h.hashMap[i] = &mn
	} else {
		mn.next = h.hashMap[i]
		h.hashMap[i].prev = &mn
		h.hashMap[i] = &mn
	}

	// Building list node
	ln := listNode{mapItem: &mn}

	h.count++
	h.served++
	if h.last == nil {
		// Empty nonces holder
		h.first = &ln
		h.last = &ln
	} else {
		// Moving last to new one
		h.last.next = &ln
		h.last = &ln

		// Truncate check
		if h.count > h.sizeList {

			// Truncate requested
			h.count--

			// Removing from map
			t := h.first.mapItem
			j := t.value.trim(h.sizeMap)
			h.load[j]--
			if t.next != nil && t.prev != nil {
				// Entry inside list
				t.prev.next = t.next
				t.next.prev = t.prev
			} else if t.prev != nil {
				// Entry is last
				t.prev.next = nil
			} else if t.next != nil {
				// Entry is first
				t.next.prev = nil
				h.hashMap[j] = t.next
			} else {
				// Entry is the only one
				h.hashMap[j] = nil
			}

			// Removing from list
			h.first = h.first.next
		}
	}

	return nil
}

// Returns true if map contains provided NONCE
func (h *NoncesHolder) Has(value HString) bool {
	i := value.trim(h.sizeMap)

	var mn *mapNode

	mn = h.hashMap[i]

	for mn != nil {
		if mn.value.HashCode == value.HashCode && mn.value.Value == value.Value {
			return true
		}

		mn = mn.next
	}

	return false
}
