package boolbits

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
)

// BitSet represents a bit mask whose size is an arbitrary multiple of 64 bits.
type BitSet struct {
	Words    []uint64 // Underlying Words (1 word = 64 bits)
	NumBits  int      // Total number of bits (must be >0 and divisible by 64)
	numWords int      // Words = NumBits / 64
}

// NewBitSet creates a new BitSet with the specified number of bits.
// numBits must be a positive multiple of 64. Otherwise it returns an error.
func NewBitSet(numBits int) (*BitSet, error) {
	if numBits <= 0 || numBits%64 != 0 {
		return nil, fmt.Errorf("error: numBits must be a positive multiple of 64 (got %d)", numBits)
	}
	numWords := numBits / 64
	return &BitSet{
		Words:    make([]uint64, numWords),
		NumBits:  numBits,
		numWords: numWords,
	}, nil
}

// NewBitSetFromHex initializes a BitSet from a hex string.
// The hex string length must correspond exactly to numBits (numBits/4 hex characters).
// numBits must be a multiple of 64.
func NewBitSetFromHex(numBits int, hexStr string) (*BitSet, error) {
	if numBits <= 0 || numBits%64 != 0 {
		return nil, fmt.Errorf("error: numBits must be a positive multiple of 64 (got %d)", numBits)
	}
	expectedHexLen := numBits / 4 // each hex digit represents 4 bits
	if len(hexStr) != expectedHexLen {
		return nil, fmt.Errorf("error: hex string must be exactly %d characters long (got %d)", expectedHexLen, len(hexStr))
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	expectedBytes := numBits / 8
	if len(data) != expectedBytes {
		return nil, fmt.Errorf("internal error: hex decoding mismatch, expected %d bytes, got %d", expectedBytes, len(data))
	}

	numWords := numBits / 64
	words := make([]uint64, numWords)

	// Assume the hex string is in big-endian order (MSB first).
	for i := 0; i < numWords; i++ {
		offset := i * 8
		var w uint64
		for j := 0; j < 8; j++ {
			w |= uint64(data[offset+j]) << uint((7-j)*8)
		}
		words[i] = w
	}

	return &BitSet{
		Words:    words,
		NumBits:  numBits,
		numWords: numWords,
	}, nil
}

// ToHex returns the bitset as a hex string (without "0x" prefix).
func (b *BitSet) ToHex() string {
	buf := make([]byte, b.numWords*8)
	for i := 0; i < b.numWords; i++ {
		w := b.Words[i]
		offset := i * 8
		buf[offset] = byte(w >> 56)
		buf[offset+1] = byte(w >> 48)
		buf[offset+2] = byte(w >> 40)
		buf[offset+3] = byte(w >> 32)
		buf[offset+4] = byte(w >> 24)
		buf[offset+5] = byte(w >> 16)
		buf[offset+6] = byte(w >> 8)
		buf[offset+7] = byte(w)
	}
	return hex.EncodeToString(buf)
}

// String implements fmt.Stringer and displays the hex representation with "0x" prefix.
func (b *BitSet) String() string {
	return "0x" + b.ToHex()
}

// SetBit sets the bit at index i (0 ≤ i < numBits) to 1.
func (b *BitSet) SetBit(i int) error {
	if i < 0 || i >= b.NumBits {
		return fmt.Errorf("SetBit: index %d out of valid range [0, %d)", i, b.NumBits)
	}
	wordIdx := i / 64
	bitIdx := uint(i % 64)
	b.Words[wordIdx] |= uint64(1) << bitIdx
	return nil
}

// ClearBit clears the bit at index i (0 ≤ i < numBits).
func (b *BitSet) ClearBit(i int) error {
	if i < 0 || i >= b.NumBits {
		return fmt.Errorf("ClearBit: index %d out of valid range [0, %d)", i, b.NumBits)
	}
	wordIdx := i / 64
	bitIdx := uint(i % 64)
	b.Words[wordIdx] &^= uint64(1) << bitIdx
	return nil
}

// TestBit returns true if the bit at index i (0 ≤ i < numBits) is 1.
func (b *BitSet) TestBit(i int) (bool, error) {
	if i < 0 || i >= b.NumBits {
		return false, fmt.Errorf("TestBit: index %d out of valid range [0, %d)", i, b.NumBits)
	}
	wordIdx := i / 64
	bitIdx := uint(i % 64)
	return (b.Words[wordIdx]>>bitIdx)&1 == 1, nil
}

// IsZero returns true if all bits are zero.
func (b *BitSet) IsZero() bool {
	for _, w := range b.Words {
		if w != 0 {
			return false
		}
	}
	return true
}

// CountOnes counts the number of set bits (popcount) in the entire bitset.
func (b *BitSet) CountOnes() int {
	count := 0
	for _, w := range b.Words {
		count += bits.OnesCount64(w)
	}
	return count
}

// ensureSameSize checks that two BitSets have the same numBits.
func ensureSameSize(a, o *BitSet) error {
	if a.NumBits != o.NumBits {
		return errors.New("bitset sizes differ")
	}
	return nil
}

// And performs a bitwise AND (∧) between two BitSets (must have the same numBits).
func (b *BitSet) And(o *BitSet) (*BitSet, error) {
	if err := ensureSameSize(b, o); err != nil {
		return nil, err
	}
	result := make([]uint64, b.numWords)
	for i := 0; i < b.numWords; i++ {
		result[i] = b.Words[i] & o.Words[i]
	}
	return &BitSet{
		Words:    result,
		NumBits:  b.NumBits,
		numWords: b.numWords,
	}, nil
}

// Or performs a bitwise OR (∨) between two BitSets.
func (b *BitSet) Or(o *BitSet) (*BitSet, error) {
	if err := ensureSameSize(b, o); err != nil {
		return nil, err
	}
	result := make([]uint64, b.numWords)
	for i := 0; i < b.numWords; i++ {
		result[i] = b.Words[i] | o.Words[i]
	}
	return &BitSet{
		Words:    result,
		NumBits:  b.NumBits,
		numWords: b.numWords,
	}, nil
}

// Xor performs a bitwise XOR (⊕) between two BitSets.
func (b *BitSet) Xor(o *BitSet) (*BitSet, error) {
	if err := ensureSameSize(b, o); err != nil {
		return nil, err
	}
	result := make([]uint64, b.numWords)
	for i := 0; i < b.numWords; i++ {
		result[i] = b.Words[i] ^ o.Words[i]
	}
	return &BitSet{
		Words:    result,
		NumBits:  b.NumBits,
		numWords: b.numWords,
	}, nil
}

// Not inverts all bits in this BitSet (bitwise NOT).
func (b *BitSet) Not() *BitSet {
	result := make([]uint64, b.numWords)
	for i := 0; i < b.numWords; i++ {
		result[i] = ^b.Words[i]
	}
	return &BitSet{
		Words:    result,
		NumBits:  b.NumBits,
		numWords: b.numWords,
	}
}

// Equals checks if two BitSets are equal. Returns false if numBits differ or any word differs.
func (b *BitSet) Equals(o *BitSet) bool {
	if b.NumBits != o.NumBits {
		return false
	}
	for i := 0; i < b.numWords; i++ {
		if b.Words[i] != o.Words[i] {
			return false
		}
	}
	return true
}
