package main

import (
	"Fast_BitFilter_MetaData/boolbits/bitmapper"
	"Fast_BitFilter_MetaData/boolbits/boolbits"
	"fmt"
	"log"
)

func main() {
	// Demonstrate usage of BitSet for various sizes (multiples of 64 bits)
	sizes := []int{64, 128, 256, 512, 1024, 2048}

	for _, size := range sizes {
		fmt.Printf("=== BitSet of size %d bits ===\n", size)

		// Create a new BitSet of the given size
		bs, err := boolbits.NewBitSet(size)
		if err != nil {
			log.Fatalf("Failed to create BitSet of size %d: %v\n", size, err)
		}

		// Set a few bits: 0, middle, last
		middle := size / 2
		last := size - 1

		if err := bs.SetBit(0); err != nil {
			log.Fatalf("SetBit error: %v\n", err)
		}
		if err := bs.SetBit(middle); err != nil {
			log.Fatalf("SetBit error: %v\n", err)
		}
		if err := bs.SetBit(last); err != nil {
			log.Fatalf("SetBit error: %v\n", err)
		}

		// Test bits we set
		b0, _ := bs.TestBit(0)
		bMid, _ := bs.TestBit(middle)
		bLast, _ := bs.TestBit(last)
		fmt.Printf("Bit 0 set? %v\n", b0)
		fmt.Printf("Bit %d set? %v\n", middle, bMid)
		fmt.Printf("Bit %d set? %v\n", last, bLast)

		// Count number of set bits
		count := bs.CountOnes()
		fmt.Printf("Number of bits set: %d (expected 3)\n", count)

		// Show hex representation (truncated for large sizes)
		hexStr := bs.String()
		if size <= 256 {
			fmt.Printf("Hex: %s\n", hexStr)
		} else {
			// For larger bitsets, only print first and last 16 hex characters
			prefix := hexStr[:18] // "0x" + first 16 hex digits
			suffix := hexStr[len(hexStr)-16:]
			fmt.Printf("Hex (truncated): %s...%s\n", prefix, suffix)
		}

		// Clear the middle bit and verify
		if err := bs.ClearBit(middle); err != nil {
			log.Fatalf("ClearBit error: %v\n", err)
		}
		bMidCleared, _ := bs.TestBit(middle)
		fmt.Printf("After clearing, bit %d set? %v (expected false)\n", middle, bMidCleared)

		// Demonstrate NOT operation on a smaller slice (only print first 4 words)
		notBS := bs.Not()
		fmt.Printf("First 64 bits of NOT: 0x%016x\n", notBS.Words[0])
		fmt.Println()
	}

	// Demonstrate bitwise operations between two 256-bit sets
	fmt.Println("=== Bitwise operations on two 256-bit sets ===")
	a256, err := boolbits.NewBitSet(256)
	if err != nil {
		log.Fatalf("Failed to create 256-bit BitSet a: %v\n", err)
	}
	b256, err := boolbits.NewBitSet(256)
	if err != nil {
		log.Fatalf("Failed to create 256-bit BitSet b: %v\n", err)
	}

	// Set bits in 'a': bits 1, 100, 200, 255
	for _, idx := range []int{1, 100, 200, 255} {
		if err := a256.SetBit(idx); err != nil {
			log.Fatalf("SetBit a256 error: %v\n", err)
		}
	}
	// Set bits in 'b': bits 1, 150, 200
	for _, idx := range []int{1, 150, 200} {
		if err := b256.SetBit(idx); err != nil {
			log.Fatalf("SetBit b256 error: %v\n", err)
		}
	}

	fmt.Printf("a256 bits set: %d (expected 4)\n", a256.CountOnes())
	fmt.Printf("b256 bits set: %d (expected 3)\n", b256.CountOnes())

	andRes, err := a256.And(b256)
	if err != nil {
		log.Fatalf("And error: %v\n", err)
	}
	orRes, err := a256.Or(b256)
	if err != nil {
		log.Fatalf("Or error: %v\n", err)
	}
	xorRes, err := a256.Xor(b256)
	if err != nil {
		log.Fatalf("Xor error: %v\n", err)
	}

	fmt.Printf("AND result bits set: %d (expected 2: bits 1 and 200)\n", andRes.CountOnes())
	fmt.Printf("OR result bits set: %d (expected 5: bits 1, 100, 150, 200, 255)\n", orRes.CountOnes())
	fmt.Printf("XOR result bits set: %d (expected 3: bits 100, 150, 255)\n", xorRes.CountOnes())

	fmt.Printf("a256 hex: %s\n", a256.String())
	fmt.Printf("b256 hex: %s\n", b256.String())
	fmt.Printf("AND hex: %s\n", andRes.String())
	fmt.Printf("OR hex:  %s\n", orRes.String())
	fmt.Printf("XOR hex: %s\n", xorRes.String())

	domainMap, groupMap, nameMap, valueMap, err := bitmapper.GenerateBitMaps(
		[]string{"domain1", "domain2", "domain1", "domain3"},
		[]string{"group1", "group2", "group3", "groupA", "groupB"},
		[]string{"nameA", "nameB", "nameA", "nameY", "nameZ"},
		[]string{"valX", "valY", "val2", "val3"},
	)

	if err != nil {
		log.Fatalf("GenerateBitMaps error: %v\n", err)
	}
	fmt.Println("domainMap", domainMap)
	fmt.Println("groupMap", groupMap)
	fmt.Println("nameMap", nameMap)
	fmt.Println("valueMap", valueMap)

	entry, err := bitmapper.NewEntry(
		"domain2", "group1", "nameA", "valY",
		domainMap, groupMap, nameMap, valueMap,
	)
	if err != nil {
		log.Fatalf("NewEntry error: %v", err)
	}

	// Print out the four BitSets inside the entry:
	fmt.Println("\n---- Example Entry ----")
	fmt.Printf("Domain BitSet for %q: %s\n", "domain2", entry.Domain.String())
	fmt.Printf(" Group BitSet for %q: %s\n", "group1", entry.Group.String())
	fmt.Printf("  Name BitSet for %q: %s\n", "nameA", entry.Name.String())
	fmt.Printf(" Value BitSet for %q: %s\n", "valY", entry.Value.String())

	//  Create two Entry objects:
	//    - entryA uses ("domain2", "groupA", "nameY", "val3")
	//    - entryB uses ("domain2", "groupA", "nameY", "val3") → identical to entryA
	//    - entryC uses ("domain1", "groupB", "nameZ", "val2") → different combination
	entryA, err := bitmapper.NewEntry(
		"domain2", "groupA", "nameY", "val3",
		domainMap, groupMap, nameMap, valueMap,
	)
	if err != nil {
		log.Fatalf("NewEntry A error: %v", err)
	}

	entryB, err := bitmapper.NewEntry(
		"domain2", "groupA", "nameY", "val3",
		domainMap, groupMap, nameMap, valueMap,
	)
	if err != nil {
		log.Fatalf("NewEntry B error: %v", err)
	}

	entryC, err := bitmapper.NewEntry(
		"domain1", "groupB", "nameZ", "val2",
		domainMap, groupMap, nameMap, valueMap,
	)
	if err != nil {
		log.Fatalf("NewEntry C error: %v", err)
	}

	// 4) Compare using Equals()
	fmt.Printf("entryA.Equals(entryB)? %v (expected true)\n", entryA.Equals(entryB))
	fmt.Printf("entryA.Equals(entryC)? %v (expected false)\n", entryA.Equals(entryC))

	// 5) Print out the four BitSets inside entryA for reference
	fmt.Println("\n-- entryA BitSets (hex) --")
	fmt.Printf(" Domain (%q): %s\n", "domain2", entryA.Domain.String())
	fmt.Printf("  Group (%q): %s\n", "groupA", entryA.Group.String())
	fmt.Printf("   Name (%q): %s\n", "nameY", entryA.Name.String())
	fmt.Printf("  Value (%q): %s\n", "val3", entryA.Value.String())

}
