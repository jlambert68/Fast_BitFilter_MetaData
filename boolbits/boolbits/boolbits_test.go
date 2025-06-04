package boolbits

/*
boolbits_test.go to cover:
Creation of valid and invalid sizes.
SetBit, TestBit, ClearBit, IsZero, and CountOnes across various lengths.
Equality comparisons (Equals) including mismatched sizes.
Round-trip hex conversions (NewBitSetFromHex and ToHex).
Bitwise operations (And, Or, Xor, Not).
*/

import (
	"testing"
)

func TestNewBitSetInvalidSize(t *testing.T) {
	// Sizes not multiples of 64 should return an error
	invalidSizes := []int{0, 1, 63, 65, -64, 100}
	for _, size := range invalidSizes {
		if _, err := NewBitSet(size); err == nil {
			t.Errorf("Expected error for size %d, got nil", size)
		}
	}
}

func TestSetTestClearAndCount(t *testing.T) {
	sizes := []int{64, 128, 256}
	for _, size := range sizes {
		bs, err := NewBitSet(size)
		if err != nil {
			t.Fatalf("Failed to create BitSet of size %d: %v", size, err)
		}

		// Initially all bits should be zero
		if !bs.IsZero() {
			t.Errorf("BitSet of size %d should be zero after creation", size)
		}
		if count := bs.CountOnes(); count != 0 {
			t.Errorf("Expected 0 ones, got %d", count)
		}

		// Set first, middle, and last bits
		positions := []int{0, size / 2, size - 1}
		for _, pos := range positions {
			if err := bs.SetBit(pos); err != nil {
				t.Errorf("SetBit(%d) error for size %d: %v", pos, size, err)
			}
			if val, err := bs.TestBit(pos); err != nil || !val {
				t.Errorf("TestBit(%d) expected true for size %d, got val=%v err=%v", pos, size, val, err)
			}
		}

		// Now CountOnes should equal number of positions
		if count := bs.CountOnes(); count != len(positions) {
			t.Errorf("CountOnes expected %d, got %d for size %d", len(positions), count, size)
		}

		// Clear middle bit
		mid := size / 2
		if err := bs.ClearBit(mid); err != nil {
			t.Errorf("ClearBit(%d) error for size %d: %v", mid, size, err)
		}
		if val, err := bs.TestBit(mid); err != nil || val {
			t.Errorf("TestBit(%d) expected false after clear for size %d, got val=%v err=%v", mid, size, val, err)
		}
		if count := bs.CountOnes(); count != len(positions)-1 {
			t.Errorf("After clearing, CountOnes expected %d, got %d for size %d", len(positions)-1, count, size)
		}
	}
}

func TestEqualsDifferentSizes(t *testing.T) {
	bs64, _ := NewBitSet(64)
	bs128, _ := NewBitSet(128)
	if bs64.Equals(bs128) {
		t.Error("Equals should return false for different sizes (64 vs 128)")
	}
}

func TestEqualsSameContent(t *testing.T) {
	bsA, _ := NewBitSet(128)
	bsB, _ := NewBitSet(128)
	if !bsA.Equals(bsB) {
		t.Error("Two zeroed BitSets should be equal")
	}
	// Set same bits in both
	bsA.SetBit(5)
	bsA.SetBit(127)
	bsB.SetBit(5)
	bsB.SetBit(127)
	if !bsA.Equals(bsB) {
		t.Error("BitSets with identical bits set should be equal")
	}
	// Change one bit
	bsB.ClearBit(5)
	if bsA.Equals(bsB) {
		t.Error("BitSets should not be equal after one bit differs")
	}
}

func TestNewBitSetFromHexAndToHex(t *testing.T) {
	examples := []struct {
		size   int
		hexStr string
	}{
		{64, "0123456789abcdef"},                                                  // 8 bytes = 16 hex digits
		{128, "0123456789abcdef0123456789abcdef"},                                 // 16 bytes = 32 hex digits
		{256, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"}, // 32 bytes = 64 hex digits
	}
	for _, ex := range examples {
		bs, err := NewBitSetFromHex(ex.size, ex.hexStr)
		if err != nil {
			t.Errorf("NewBitSetFromHex failed for size %d: %v", ex.size, err)
			continue
		}
		outHex := bs.ToHex()
		if outHex != ex.hexStr {
			t.Errorf("ToHex roundtrip failed for size %d: expected %s, got %s", ex.size, ex.hexStr, outHex)
		}
	}

	// Invalid hex length
	if _, err := NewBitSetFromHex(64, "123"); err == nil {
		t.Error("Expected error for invalid hex length, got nil")
	}
}

func TestBitwiseOperations(t *testing.T) {
	size := 256
	a, _ := NewBitSet(size)
	b, _ := NewBitSet(size)
	// Set bits in 'a': 0, 100, 200, 255
	for _, pos := range []int{0, 100, 200, 255} {
		a.SetBit(pos)
	}
	// Set bits in 'b': 0, 150, 200
	for _, pos := range []int{0, 150, 200} {
		b.SetBit(pos)
	}

	andRes, err := a.And(b)
	if err != nil {
		t.Fatalf("And returned error: %v", err)
	}
	// AND should have bits 0 and 200
	val0, err0 := andRes.TestBit(0)
	val200, err200 := andRes.TestBit(200)
	if err0 != nil || err200 != nil || !val0 || !val200 || andRes.CountOnes() != 2 {
		t.Errorf("And result incorrect: got %d ones", andRes.CountOnes())
	}

	orRes, err := a.Or(b)
	if err != nil {
		t.Fatalf("Or returned error: %v", err)
	}
	// OR should have bits 0, 100, 150, 200, 255 => count 5
	if orRes.CountOnes() != 5 {
		t.Errorf("Or result incorrect: got %d ones, expected 5", orRes.CountOnes())
	}

	xorRes, err := a.Xor(b)
	if err != nil {
		t.Fatalf("Xor returned error: %v", err)
	}
	// XOR should have bits 100, 150, 255 => count 3
	if xorRes.CountOnes() != 3 {
		t.Errorf("Xor result incorrect: got %d ones, expected 3", xorRes.CountOnes())
	}

	// Test Not: invert 'a'
	notA := a.Not()
	// Bits set in notA should be all except 0,100,200,255 (i.e., size-4 total ones)
	expectedOnes := size - 4
	if notA.CountOnes() != expectedOnes {
		t.Errorf("Not result incorrect: got %d ones, expected %d", notA.CountOnes(), expectedOnes)
	}
}
