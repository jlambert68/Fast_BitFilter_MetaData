package boolbits

import (
	"fmt"
)

// Entry holds the BitSet for a single combination of domain, group, name, and value.
// Each field is a pointer to a BitSet.
type Entry struct {
	Domain *BitSet
	Group  *BitSet
	Name   *BitSet
	Value  *BitSet
}

// NewEntry constructs an Entry given four BitSet pointers.
// Returns an error if any field is nil.
func NewEntry(domainBS, groupBS, nameBS, valueBS *BitSet) (*Entry, error) {
	if domainBS == nil {
		return nil, fmt.Errorf("domain BitSet is nil")
	}
	if groupBS == nil {
		return nil, fmt.Errorf("group BitSet is nil")
	}
	if nameBS == nil {
		return nil, fmt.Errorf("name BitSet is nil")
	}
	if valueBS == nil {
		return nil, fmt.Errorf("value BitSet is nil")
	}
	return &Entry{
		Domain: domainBS,
		Group:  groupBS,
		Name:   nameBS,
		Value:  valueBS,
	}, nil
}

// Equals compares two Entries. Returns true if all corresponding BitSets are equal.
func (e *Entry) Equals(o *Entry) bool {
	if e == nil || o == nil {
		return false
	}
	return e.Domain.Equals(o.Domain) &&
		e.Group.Equals(o.Group) &&
		e.Name.Equals(o.Name) &&
		e.Value.Equals(o.Value)
}

// And returns a new Entry by performing bitwise AND on corresponding BitSets.
func (e *Entry) And(o *Entry) (*Entry, error) {
	if e == nil || o == nil {
		return nil, fmt.Errorf("cannot AND nil Entry")
	}
	// Ensure bit lengths match for each field
	if e.Domain.NumBits != o.Domain.NumBits {
		return nil, fmt.Errorf("mismatched Domain bit lengths: %d vs %d", e.Domain.NumBits, o.Domain.NumBits)
	}
	if e.Group.NumBits != o.Group.NumBits {
		return nil, fmt.Errorf("mismatched Group bit lengths: %d vs %d", e.Group.NumBits, o.Group.NumBits)
	}
	if e.Name.NumBits != o.Name.NumBits {
		return nil, fmt.Errorf("mismatched Name bit lengths: %d vs %d", e.Name.NumBits, o.Name.NumBits)
	}
	if e.Value.NumBits != o.Value.NumBits {
		return nil, fmt.Errorf("mismatched Value bit lengths: %d vs %d", e.Value.NumBits, o.Value.NumBits)
	}

	domainRes, err := e.Domain.And(o.Domain)
	if err != nil {
		return nil, fmt.Errorf("Domain AND error: %v", err)
	}
	groupRes, err := e.Group.And(o.Group)
	if err != nil {
		return nil, fmt.Errorf("Group AND error: %v", err)
	}
	nameRes, err := e.Name.And(o.Name)
	if err != nil {
		return nil, fmt.Errorf("Name AND error: %v", err)
	}
	valueRes, err := e.Value.And(o.Value)
	if err != nil {
		return nil, fmt.Errorf("Value AND error: %v", err)
	}
	return &Entry{Domain: domainRes, Group: groupRes, Name: nameRes, Value: valueRes}, nil
}

// Or returns a new Entry by performing bitwise OR on corresponding BitSets.
func (e *Entry) Or(o *Entry) (*Entry, error) {
	if e == nil || o == nil {
		return nil, fmt.Errorf("cannot OR nil Entry")
	}
	// Ensure bit lengths match
	if e.Domain.NumBits != o.Domain.NumBits {
		return nil, fmt.Errorf("mismatched Domain bit lengths: %d vs %d", e.Domain.NumBits, o.Domain.NumBits)
	}
	if e.Group.NumBits != o.Group.NumBits {
		return nil, fmt.Errorf("mismatched Group bit lengths: %d vs %d", e.Group.NumBits, o.Group.NumBits)
	}
	if e.Name.NumBits != o.Name.NumBits {
		return nil, fmt.Errorf("mismatched Name bit lengths: %d vs %d", e.Name.NumBits, o.Name.NumBits)
	}
	if e.Value.NumBits != o.Value.NumBits {
		return nil, fmt.Errorf("mismatched Value bit lengths: %d vs %d", e.Value.NumBits, o.Value.NumBits)
	}

	domainRes, err := e.Domain.Or(o.Domain)
	if err != nil {
		return nil, fmt.Errorf("Domain OR error: %v", err)
	}
	groupRes, err := e.Group.Or(o.Group)
	if err != nil {
		return nil, fmt.Errorf("Group OR error: %v", err)
	}
	nameRes, err := e.Name.Or(o.Name)
	if err != nil {
		return nil, fmt.Errorf("Name OR error: %v", err)
	}
	valueRes, err := e.Value.Or(o.Value)
	if err != nil {
		return nil, fmt.Errorf("Value OR error: %v", err)
	}
	return &Entry{Domain: domainRes, Group: groupRes, Name: nameRes, Value: valueRes}, nil
}

// Xor returns a new Entry by performing bitwise XOR on corresponding BitSets.
func (e *Entry) Xor(o *Entry) (*Entry, error) {
	if e == nil || o == nil {
		return nil, fmt.Errorf("cannot XOR nil Entry")
	}
	// Ensure bit lengths match
	if e.Domain.NumBits != o.Domain.NumBits {
		return nil, fmt.Errorf("mismatched Domain bit lengths: %d vs %d", e.Domain.NumBits, o.Domain.NumBits)
	}
	if e.Group.NumBits != o.Group.NumBits {
		return nil, fmt.Errorf("mismatched Group bit lengths: %d vs %d", e.Group.NumBits, o.Group.NumBits)
	}
	if e.Name.NumBits != o.Name.NumBits {
		return nil, fmt.Errorf("mismatched Name bit lengths: %d vs %d", e.Name.NumBits, o.Name.NumBits)
	}
	if e.Value.NumBits != o.Value.NumBits {
		return nil, fmt.Errorf("mismatched Value bit lengths: %d vs %d", e.Value.NumBits, o.Value.NumBits)
	}

	domainRes, err := e.Domain.Xor(o.Domain)
	if err != nil {
		return nil, fmt.Errorf("Domain XOR error: %v", err)
	}
	groupRes, err := e.Group.Xor(o.Group)
	if err != nil {
		return nil, fmt.Errorf("Group XOR error: %v", err)
	}
	nameRes, err := e.Name.Xor(o.Name)
	if err != nil {
		return nil, fmt.Errorf("Name XOR error: %v", err)
	}
	valueRes, err := e.Value.Xor(o.Value)
	if err != nil {
		return nil, fmt.Errorf("Value XOR error: %v", err)
	}
	return &Entry{Domain: domainRes, Group: groupRes, Name: nameRes, Value: valueRes}, nil
}

// Not returns a new Entry by performing bitwise NOT on each BitSet.
func (e *Entry) Not() (*Entry, error) {
	if e == nil {
		return nil, fmt.Errorf("cannot NOT nil Entry")
	}
	domainRes := e.Domain.Not()
	groupRes := e.Group.Not()
	nameRes := e.Name.Not()
	valueRes := e.Value.Not()
	return &Entry{Domain: domainRes, Group: groupRes, Name: nameRes, Value: valueRes}, nil
}

// NewAllOnesEntry constructs an Entry where each BitSet has all bits set to 1.
// bitLen must be a positive multiple of 64; returns an error otherwise.
func NewAllOnesEntry(bitLen int) (*Entry, error) {
	// Validate bitLen
	if bitLen <= 0 || bitLen%64 != 0 {
		return nil, fmt.Errorf("bit length must be a positive multiple of 64 (got %d)", bitLen)
	}
	// Number of 64-bit words
	numWords := bitLen / 64
	// Create a BitSet and set all bits by filling each word with all ones
	fillAllOnes := func() *BitSet {
		b := make([]uint64, numWords)
		for i := 0; i < numWords; i++ {
			b[i] = ^uint64(0)
		}
		return &BitSet{Words: b, NumBits: bitLen, numWords: numWords}
	}
	domainBS := fillAllOnes()
	groupBS := fillAllOnes()
	nameBS := fillAllOnes()
	valueBS := fillAllOnes()
	return &Entry{Domain: domainBS, Group: groupBS, Name: nameBS, Value: valueBS}, nil
}
