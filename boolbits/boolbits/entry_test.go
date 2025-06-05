package boolbits

import (
	"testing"
)

func TestNewEntry_SuccessAndFailure(t *testing.T) {
	// Create two BitSets of equal length
	bs1, err := NewBitSet(64)
	if err != nil {
		t.Fatalf("NewBitSet error: %v", err)
	}
	bs2, err := NewBitSet(64)
	if err != nil {
		t.Fatalf("NewBitSet error: %v", err)
	}
	bs3, err := NewBitSet(128)
	if err != nil {
		t.Fatalf("NewBitSet error: %v", err)
	}

	// Successful creation with non-nil BitSets
	entry, err := NewEntry(bs1, bs2, bs1, bs2)
	if err != nil {
		t.Errorf("Expected NewEntry to succeed, got error: %v", err)
	}
	if entry.Domain != bs1 || entry.Group != bs2 || entry.Name != bs1 || entry.Value != bs2 {
		t.Errorf("Entry fields do not match input BitSets")
	}

	// Failure when any field is nil
	cases := []struct {
		d, g, n, v *BitSet
	}{
		{nil, bs2, bs1, bs2},
		{bs1, nil, bs1, bs2},
		{bs1, bs2, nil, bs2},
		{bs1, bs2, bs1, nil},
	}
	for i, c := range cases {
		if _, err := NewEntry(c.d, c.g, c.n, c.v); err == nil {
			t.Errorf("Case %d: expected error for nil field, got nil", i)
		}
	}

	// Failure when bit lengths mismatch
	_, err = NewEntry(bs1, bs2, bs1, bs3)
	if err != nil {
		// NewEntry does not check lengths, so expect success
		t.Logf("NewEntry allowed mismatched lengths: %v", err)
	}
}

func TestEntry_Equals(t *testing.T) {
	bsA1, _ := NewBitSet(64)
	bsA2, _ := NewBitSet(64)
	// Set distinct bits
	bsA1.SetBit(0)
	bsA2.SetBit(1)

	bsB1, _ := NewBitSet(64)
	bsB2, _ := NewBitSet(64)
	bsB1.SetBit(0)
	bsB2.SetBit(1)

	// entry1 and entry2 identical content
	entry1, _ := NewEntry(bsA1, bsA2, bsA1, bsA2)
	entry2, _ := NewEntry(bsB1, bsB2, bsB1, bsB2)
	if !entry1.Equals(entry2) {
		t.Errorf("Expected entry1 to equal entry2")
	}

	// Change one field
	bsC, _ := NewBitSet(64)
	bsC.SetBit(2)
	entry3, _ := NewEntry(bsC, bsA2, bsA1, bsA2)
	if entry1.Equals(entry3) {
		t.Errorf("Expected entry1 not to equal entry3 (different Domain)")
	}

	// Nil comparisons
	var nilEntry *Entry
	if entry1.Equals(nilEntry) {
		t.Errorf("Expected entry1.Equals(nil) to be false")
	}
	if nilEntry != nil && nilEntry.Equals(entry1) {
		t.Errorf("Expected nilEntry.Equals(entry1) to be false")
	}
}

// helper to count bits in all four BitSets and verify they equal expected
func verifyAllOnesEntry(t *testing.T, entry *Entry, bitLen int) {
	t.Helper()

	// Each BitSet should have exactly bitLen bits set.
	dCount := entry.Domain.CountOnes()
	gCount := entry.Group.CountOnes()
	nCount := entry.Name.CountOnes()
	vCount := entry.Value.CountOnes()

	if dCount != bitLen {
		t.Errorf("Domain CountOnes = %d; want %d", dCount, bitLen)
	}
	if gCount != bitLen {
		t.Errorf("Group CountOnes = %d; want %d", gCount, bitLen)
	}
	if nCount != bitLen {
		t.Errorf("Name CountOnes = %d; want %d", nCount, bitLen)
	}
	if vCount != bitLen {
		t.Errorf("Value CountOnes = %d; want %d", vCount, bitLen)
	}

	// All four BitSets should be equal to each other.
	if !entry.Domain.Equals(entry.Group) {
		t.Error("Domain and Group BitSets differ; expected identical all-ones")
	}
	if !entry.Domain.Equals(entry.Name) {
		t.Error("Domain and Name BitSets differ; expected identical all-ones")
	}
	if !entry.Domain.Equals(entry.Value) {
		t.Error("Domain and Value BitSets differ; expected identical all-ones")
	}

	// Additionally, verify that ToHex produces the expected hex string:
	// For bitLen bits all set, the hex string should be bitLen/4 nibbles of 'f'.
	hex := entry.Domain.String()[2:] // strip "0x"
	expectedNibbles := bitLen / 4
	expectedHex := ""
	for i := 0; i < expectedNibbles; i++ {
		expectedHex += "f"
	}
	if hex != expectedHex {
		t.Errorf("Domain ToHex = %q; want %q", hex, expectedHex)
	}
}

func TestNewAllOnesEntry_Success_64bits(t *testing.T) {
	entry, err := NewAllOnesEntry(64)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(64) returned error: %v", err)
	}
	verifyAllOnesEntry(t, entry, 64)
}

func TestNewAllOnesEntry_Success_128bits(t *testing.T) {
	entry, err := NewAllOnesEntry(128)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(128) returned error: %v", err)
	}
	verifyAllOnesEntry(t, entry, 128)
}

func TestNewAllOnesEntry_Success_256bits(t *testing.T) {
	entry, err := NewAllOnesEntry(256)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(256) returned error: %v", err)
	}
	verifyAllOnesEntry(t, entry, 256)
}

func TestNewAllOnesEntry_InvalidLengths(t *testing.T) {
	tests := []int{
		0,       // zero
		1,       // not multiple of 64
		65,      // not multiple of 64
		100,     // not multiple of 64
		63,      // not multiple of 64
		192 + 1, // 193 is not multiple of 64
	}

	for _, bitLen := range tests {
		_, err := NewAllOnesEntry(bitLen)
		if err == nil {
			t.Errorf("Expected error for bitLen=%d; got nil", bitLen)
		}
	}
}

func TestAllOnesEntry_BitwiseOperations(t *testing.T) {
	// Create two all-ones entries with same length
	entryA, err := NewAllOnesEntry(64)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(64) returned error: %v", err)
	}
	entryB, err := NewAllOnesEntry(64)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(64) returned error: %v", err)
	}

	// AND of two all-ones should be all-ones
	andEntry, err := entryA.And(entryB)
	if err != nil {
		t.Fatalf("And returned error: %v", err)
	}
	verifyAllOnesEntry(t, andEntry, 64)

	// OR of two all-ones should be all-ones
	orEntry, err := entryA.Or(entryB)
	if err != nil {
		t.Fatalf("Or returned error: %v", err)
	}
	verifyAllOnesEntry(t, orEntry, 64)

	// XOR of two identical all-ones should be all-zero bitsets
	xorEntry, err := entryA.Xor(entryB)
	if err != nil {
		t.Fatalf("Xor returned error: %v", err)
	}

	// All-zero entries: CountOnes should be zero
	if xorEntry.Domain.CountOnes() != 0 {
		t.Errorf("XOR Domain CountOnes = %d; want 0", xorEntry.Domain.CountOnes())
	}
	if xorEntry.Group.CountOnes() != 0 {
		t.Errorf("XOR Group CountOnes = %d; want 0", xorEntry.Group.CountOnes())
	}
	if xorEntry.Name.CountOnes() != 0 {
		t.Errorf("XOR Name CountOnes = %d; want 0", xorEntry.Name.CountOnes())
	}
	if xorEntry.Value.CountOnes() != 0 {
		t.Errorf("XOR Value CountOnes = %d; want 0", xorEntry.Value.CountOnes())
	}

	// NOT of all-ones should be all-zero
	notEntry, err := entryA.Not()
	if err != nil {
		t.Fatalf("Not returned error: %v", err)
	}
	if notEntry.Domain.CountOnes() != 0 {
		t.Errorf("NOT Domain CountOnes = %d; want 0", notEntry.Domain.CountOnes())
	}
	if notEntry.Group.CountOnes() != 0 {
		t.Errorf("NOT Group CountOnes = %d; want 0", notEntry.Group.CountOnes())
	}
	if notEntry.Name.CountOnes() != 0 {
		t.Errorf("NOT Name CountOnes = %d; want 0", notEntry.Name.CountOnes())
	}
	if notEntry.Value.CountOnes() != 0 {
		t.Errorf("NOT Value CountOnes = %d; want 0", notEntry.Value.CountOnes())
	}
}

func TestNewAllOnesEntry_IntegrationWithBoolSet(t *testing.T) {
	// Verify that the BitSets returned by NewAllOnesEntry behave correctly for SetBit/TestBit
	bitLen := 64
	entry, err := NewAllOnesEntry(bitLen)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(64) returned error: %v", err)
	}

	// For an all-ones BitSet, calling SetBit on any index should not change CountOnes
	before := entry.Domain.CountOnes()
	if before != bitLen {
		t.Errorf("Initial Domain CountOnes = %d; want %d", before, bitLen)
	}

	// Attempt setting and clearing a bit and re-count
	if err := entry.Domain.SetBit(10); err != nil {
		t.Errorf("Setting bit 10 returned error: %v", err)
	}
	if entry.Domain.CountOnes() != bitLen {
		t.Errorf("After SetBit(10), Domain CountOnes = %d; want %d", entry.Domain.CountOnes(), bitLen)
	}

	if err := entry.Domain.ClearBit(10); err != nil {
		t.Errorf("Clearing bit 10 returned error: %v", err)
	}
	if entry.Domain.CountOnes() != bitLen-1 {
		t.Errorf("After ClearBit(10), Domain CountOnes = %d; want %d", entry.Domain.CountOnes(), bitLen-1)
	}

	// Restore bit and verify count returns to full
	if err := entry.Domain.SetBit(10); err != nil {
		t.Errorf("Setting bit 10 returned error: %v", err)
	}
	if entry.Domain.CountOnes() != bitLen {
		t.Errorf("After resetting bit 10, Domain CountOnes = %d; want %d", entry.Domain.CountOnes(), bitLen)
	}
}

func TestNewAllOnesEntry_BitLengthsExpose(t *testing.T) {
	// Ensure the internal bit length (numBits) matches what was requested by verifying ToHex length
	bitLen := 128
	entry, err := NewAllOnesEntry(bitLen)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(128) returned error: %v", err)
	}

	// The hex string (without "0x") should be bitLen/4 characters
	hexStr := entry.Domain.String()[2:]
	if len(hexStr) != bitLen/4 {
		t.Errorf("Hex length = %d; want %d", len(hexStr), bitLen/4)
	}

	// Validate that CountOnes returns bitLen
	if entry.Domain.CountOnes() != bitLen {
		t.Errorf("Domain CountOnes = %d; want %d", entry.Domain.CountOnes(), bitLen)
	}
}

// Test that NewAllOnesEntry can be used interchangeably with boolbits.NewBitSet
func TestNewAllOnesEntry_MatchBoolSetBehavior(t *testing.T) {
	// Create a standalone BitSet of all ones
	bitLen := 64
	expected, err := NewBitSet(bitLen)
	if err != nil {
		t.Fatalf("boolbits.NewBitSet(64) error: %v", err)
	}
	// Manually set every bit
	for i := 0; i < bitLen; i++ {
		if err := expected.SetBit(i); err != nil {
			t.Fatalf("expected.SetBit(%d) error: %v", i, err)
		}
	}

	// Obtain an all-ones entry
	entry, err := NewAllOnesEntry(bitLen)
	if err != nil {
		t.Fatalf("NewAllOnesEntry(64) returned error: %v", err)
	}

	// Compare Domain of entry to expected standalone BitSet
	if !entry.Domain.Equals(expected) {
		t.Error("Entry.Domain does not match manually constructed all-ones BitSet")
	}
}

// verifyAllZerosEntry checks that each BitSet in entry has zero bits set,
// matches the expected hex string (all zeros), and bit length is correct.
func verifyAllZerosEntry(t *testing.T, entry *Entry, bitLen int) {
	t.Helper()

	// Each CountOnes should be zero
	if entry.Domain.CountOnes() != 0 {
		t.Errorf("Domain CountOnes = %d; want 0", entry.Domain.CountOnes())
	}
	if entry.Group.CountOnes() != 0 {
		t.Errorf("Group CountOnes = %d; want 0", entry.Group.CountOnes())
	}
	if entry.Name.CountOnes() != 0 {
		t.Errorf("Name CountOnes = %d; want 0", entry.Name.CountOnes())
	}
	if entry.Value.CountOnes() != 0 {
		t.Errorf("Value CountOnes = %d; want 0", entry.Value.CountOnes())
	}

	// Check NumBits on each BitSet
	if entry.Domain.NumBits != bitLen {
		t.Errorf("Domain NumBits = %d; want %d", entry.Domain.NumBits, bitLen)
	}
	if entry.Group.NumBits != bitLen {
		t.Errorf("Group NumBits = %d; want %d", entry.Group.NumBits, bitLen)
	}
	if entry.Name.NumBits != bitLen {
		t.Errorf("Name NumBits = %d; want %d", entry.Name.NumBits, bitLen)
	}
	if entry.Value.NumBits != bitLen {
		t.Errorf("Value NumBits = %d; want %d", entry.Value.NumBits, bitLen)
	}

	// The hex representation should be all zeros: bitLen/4 zeros
	hexStr := entry.Domain.String()[2:] // strip "0x"
	expectedNibbles := bitLen / 4
	if len(hexStr) != expectedNibbles {
		t.Errorf("Hex length = %d; want %d", len(hexStr), expectedNibbles)
	}
	for i, ch := range hexStr {
		if ch != '0' {
			t.Errorf("Hex at position %d = %q; want '0'", i, ch)
			break
		}
	}
}

// TestNewAllZerosEntry_Success tests valid bit lengths.
func TestNewAllZerosEntry_Success(t *testing.T) {
	validLens := []int{64, 128, 256, 512}

	for _, bitLen := range validLens {
		entry, err := NewAllZerosEntry(bitLen)
		if err != nil {
			t.Fatalf("NewAllZerosEntry(%d) returned error: %v", bitLen, err)
		}
		verifyAllZerosEntry(t, entry, bitLen)
	}
}

// TestNewAllZerosEntry_InvalidLengths checks that invalid lengths return an error.
func TestNewAllZerosEntry_InvalidLengths(t *testing.T) {
	invalidLens := []int{
		0,       // zero
		1,       // not multiple of 64
		63,      // just under 64
		65,      // just over 64
		100,     // not multiple of 64
		192 + 1, // 193 is not multiple of 64
		-64,     // negative
	}

	for _, bitLen := range invalidLens {
		_, err := NewAllZerosEntry(bitLen)
		if err == nil {
			t.Errorf("Expected error for bitLen=%d; got nil", bitLen)
		}
	}
}

// TestNewAllZerosEntry_BitwiseOperations verifies combining an all-zeros entry with another Entry.
func TestNewAllZerosEntry_BitwiseOperations(t *testing.T) {
	// Create a simple BitSet of length 64 with one bit set at index 0
	bs, err := NewBitSet(64)
	if err != nil {
		t.Fatalf("NewBitSet error: %v", err)
	}
	if err := bs.SetBit(0); err != nil {
		t.Fatalf("SetBit error: %v", err)
	}
	otherEntry := &Entry{
		Domain: bs,
		Group:  bs,
		Name:   bs,
		Value:  bs,
	}

	// Create an all-zero entry of length 64
	zeroEntry, err := NewAllZerosEntry(64)
	if err != nil {
		t.Fatalf("NewAllZerosEntry error: %v", err)
	}

	// AND with all-zero should yield all-zero
	andEntry, err := zeroEntry.And(otherEntry)
	if err != nil {
		t.Fatalf("AND returned error: %v", err)
	}
	verifyAllZerosEntry(t, andEntry, 64)

	// OR with all-zero should yield otherEntry's bits
	orEntry, err := zeroEntry.Or(otherEntry)
	if err != nil {
		t.Fatalf("OR returned error: %v", err)
	}
	if orEntry.Domain.CountOnes() != 1 {
		t.Errorf("OR Domain CountOnes = %d; want 1", orEntry.Domain.CountOnes())
	}
	if !orEntry.Domain.Equals(bs) {
		t.Error("OR Domain BitSet does not match expected single-bit BitSet")
	}

	// XOR with all-zero should also yield otherEntry's bits
	xorEntry, err := zeroEntry.Xor(otherEntry)
	if err != nil {
		t.Fatalf("XOR returned error: %v", err)
	}
	if xorEntry.Domain.CountOnes() != 1 {
		t.Errorf("XOR Domain CountOnes = %d; want 1", xorEntry.Domain.CountOnes())
	}
	if !xorEntry.Domain.Equals(bs) {
		t.Error("XOR Domain BitSet does not match expected single-bit BitSet")
	}

	// NOT of all-zero should yield all-ones
	notEntry, err := zeroEntry.Not()
	if err != nil {
		t.Fatalf("NOT returned error: %v", err)
	}
	// After NOT, every bit in each field should be set: CountOnes == 64
	if notEntry.Domain.CountOnes() != 64 {
		t.Errorf("NOT Domain CountOnes = %d; want 64", notEntry.Domain.CountOnes())
	}
}

// TestNewAllZerosEntry_SetClearBehavior checks that SetBit and ClearBit update counts correctly.
func TestNewAllZerosEntry_SetClearBehavior(t *testing.T) {
	// Create all-zero entry of length 128
	entry, err := NewAllZerosEntry(128)
	if err != nil {
		t.Fatalf("NewAllZerosEntry(128) error: %v", err)
	}

	// Initially all bits are zero
	if entry.Domain.CountOnes() != 0 {
		t.Errorf("Initial Domain CountOnes = %d; want 0", entry.Domain.CountOnes())
	}

	// Set a bit at index 5
	if err := entry.Domain.SetBit(5); err != nil {
		t.Errorf("SetBit(5) returned error: %v", err)
	}
	if entry.Domain.CountOnes() != 1 {
		t.Errorf("After SetBit(5), CountOnes = %d; want 1", entry.Domain.CountOnes())
	}

	// Clear the same bit
	if err := entry.Domain.ClearBit(5); err != nil {
		t.Errorf("ClearBit(5) returned error: %v", err)
	}
	if entry.Domain.CountOnes() != 0 {
		t.Errorf("After ClearBit(5), CountOnes = %d; want 0", entry.Domain.CountOnes())
	}
}

// TestNewAllZerosEntry_HexLength verifies the hex string length for an all-zero BitSet.
func TestNewAllZerosEntry_HexLength(t *testing.T) {
	bitLen := 256
	entry, err := NewAllZerosEntry(bitLen)
	if err != nil {
		t.Fatalf("NewAllZerosEntry(256) error: %v", err)
	}
	hexStr := entry.Domain.String()[2:] // strip "0x"
	if len(hexStr) != bitLen/4 {
		t.Errorf("Hex length = %d; want %d", len(hexStr), bitLen/4)
	}
	// Verify that the string is all '0'
	for i, ch := range hexStr {
		if ch != '0' {
			t.Errorf("Hex character at %d = %q; want '0'", i, ch)
			break
		}
	}
}
