package bitmapper

/*
bitmapper_test.go containing tests that:
Validate deduplication, correct map lengths, and expected keys.
Check each BitSet has exactly one bit set, no overlaps, and proper bit-length rounding.
Handle empty input slices.
Ensure order of first appearance determines bit-index.
Verify integration with boolbits.BitSet operations (AND/OR).
Show using reflect.DeepEqual for map comparison.
A successful creation of Entry with valid keys.
Failure scenarios for missing keys in each position.
*/

import (
	"Fast_BitFilter_MetaData/boolbits/boolbits"
	"reflect"
	"testing"
)

func TestGenerateBitMaps_DeduplicationAndAssignment(t *testing.T) {
	domains := []string{"domain1", "domain2", "domain1", "domain3"}
	groups := []string{"groupA", "groupB", "groupA"}
	names := []string{"nameX", "nameY", "nameY", "nameZ"}
	values := []string{"val1", "val2", "val1", "val3", "val2"}

	domainMap, groupMap, nameMap, valueMap, err := GenerateBitMaps(domains, groups, names, values)
	if err != nil {
		t.Fatalf("GenerateBitMaps returned unexpected error: %v", err)
	}

	// Expected unique counts
	expectedUniqueDomains := []string{"domain1", "domain2", "domain3"}
	expectedUniqueGroups := []string{"groupA", "groupB"}
	expectedUniqueNames := []string{"nameX", "nameY", "nameZ"}
	expectedUniqueValues := []string{"val1", "val2", "val3"}

	// Check map lengths
	if len(domainMap) != len(expectedUniqueDomains) {
		t.Errorf("Expected %d unique domains, got %d", len(expectedUniqueDomains), len(domainMap))
	}
	if len(groupMap) != len(expectedUniqueGroups) {
		t.Errorf("Expected %d unique groups, got %d", len(expectedUniqueGroups), len(groupMap))
	}
	if len(nameMap) != len(expectedUniqueNames) {
		t.Errorf("Expected %d unique names, got %d", len(expectedUniqueNames), len(nameMap))
	}
	if len(valueMap) != len(expectedUniqueValues) {
		t.Errorf("Expected %d unique values, got %d", len(expectedUniqueValues), len(valueMap))
	}

	// Helper to collect keys and verify presence
	verifyKeys := func(m map[string]*boolbits.BitSet, expected []string, sliceName string) {
		for _, key := range expected {
			if _, ok := m[key]; !ok {
				t.Errorf("Expected key '%s' in %s map", key, sliceName)
			}
		}
	}

	verifyKeys(domainMap, expectedUniqueDomains, "domain")
	verifyKeys(groupMap, expectedUniqueGroups, "group")
	verifyKeys(nameMap, expectedUniqueNames, "name")
	verifyKeys(valueMap, expectedUniqueValues, "value")

	// Check each BitSet has exactly one bit set and no overlaps within same map
	verifySingleBits := func(m map[string]*boolbits.BitSet, expectedCount int, sliceName string) {
		seenBits := make(map[int]struct{})
		for key, bs := range m {
			count := bs.CountOnes()
			if count != 1 {
				t.Errorf("BitSet for '%s' in %s map should have exactly 1 bit set, got %d", key, sliceName, count)
			}
			// Find index of bit set
			bitIndex := -1
			for i := 0; i < bs.NumBits; i++ {
				val, _ := bs.TestBit(i)
				if val {
					bitIndex = i
					break
				}
			}
			if bitIndex < 0 {
				t.Errorf("Could not find set bit for '%s' in %s map", key, sliceName)
			} else {
				if _, exists := seenBits[bitIndex]; exists {
					t.Errorf("Duplicate bit index %d in %s map for key '%s'", bitIndex, sliceName, key)
				} else {
					seenBits[bitIndex] = struct{}{}
				}
			}
		}
		if len(seenBits) != expectedCount {
			t.Errorf("Expected %d distinct bits in %s map, got %d", expectedCount, sliceName, len(seenBits))
		}
	}

	verifySingleBits(domainMap, len(expectedUniqueDomains), "domain")
	verifySingleBits(groupMap, len(expectedUniqueGroups), "group")
	verifySingleBits(nameMap, len(expectedUniqueNames), "name")
	verifySingleBits(valueMap, len(expectedUniqueValues), "value")

	// Verify bit length is smallest multiple of 64
	verifyBitLen := func(m map[string]*boolbits.BitSet, expectedCount int, sliceName string) {
		for _, bs := range m {
			expectedBits := ((expectedCount / 64) + 1) * 64
			if expectedCount%64 == 0 {
				expectedBits = expectedCount
			}
			if bs.NumBits != expectedBits {
				t.Errorf("BitSet bit length for %s should be %d, got %d", sliceName, expectedBits, bs.NumBits)
			}
		}
	}

	verifyBitLen(domainMap, len(expectedUniqueDomains), "domain")
	verifyBitLen(groupMap, len(expectedUniqueGroups), "group")
	verifyBitLen(nameMap, len(expectedUniqueNames), "name")
	verifyBitLen(valueMap, len(expectedUniqueValues), "value")
}

func TestGenerateBitMaps_EmptySlices(t *testing.T) {
	domainMap, groupMap, nameMap, valueMap, err := GenerateBitMaps([]string{}, []string{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("GenerateBitMaps returned unexpected error on empty input: %v", err)
	}
	// All maps should be empty
	if len(domainMap) != 0 || len(groupMap) != 0 || len(nameMap) != 0 || len(valueMap) != 0 {
		t.Errorf("Expected all maps to be empty for empty input slices")
	}
}

func TestGenerateBitMaps_OrderPreservation(t *testing.T) {
	// Order of first occurrence should determine bit index
	items := []string{"a", "b", "c", "a", "b"}
	domainMap, _, _, _, err := GenerateBitMaps(items, []string{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("GenerateBitMaps error: %v", err)
	}
	// 'a' -> index 0, 'b' -> index 1, 'c' -> index 2
	aaBS := domainMap["a"]
	bbBS := domainMap["b"]
	ccBS := domainMap["c"]

	// Check indexes
	valA, _ := aaBS.TestBit(0)
	valB, _ := bbBS.TestBit(1)
	valC, _ := ccBS.TestBit(2)
	if !valA || !valB || !valC {
		t.Errorf("Order preservation failed: expected bits at indices 0,1,2 for a,b,c")
	}

	// Ensure other bits are zero
	if count := aaBS.CountOnes(); count != 1 {
		t.Errorf("Expected 1 bit for 'a', got %d", count)
	}
	if count := bbBS.CountOnes(); count != 1 {
		t.Errorf("Expected 1 bit for 'b', got %d", count)
	}
	if count := ccBS.CountOnes(); count != 1 {
		t.Errorf("Expected 1 bit for 'c', got %d", count)
	}
}

func TestGenerateBitMaps_IntegrationWithBoolSet(t *testing.T) {
	// Ensure that operations on returned BitSet behave correctly
	domains := []string{"x", "y"}
	domainMap, _, _, _, err := GenerateBitMaps(domains, []string{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("GenerateBitMaps error: %v", err)
	}
	xBS := domainMap["x"]
	yBS := domainMap["y"]
	// xBS AND yBS should be zero BitSet
	andRes, err := xBS.And(yBS)
	if err != nil {
		t.Fatalf("And error: %v", err)
	}
	if !andRes.IsZero() {
		t.Errorf("Expected AND of disjoint BitSets to be zero")
	}
	// OR should have two bits
	orRes, err := xBS.Or(yBS)
	if err != nil {
		t.Fatalf("Or error: %v", err)
	}
	if orRes.CountOnes() != 2 {
		t.Errorf("Expected OR to have 2 bits set, got %d", orRes.CountOnes())
	}
}

// ReflectDeepEqual usage example (for possible future tests)
func TestReflectDeepEqualExample(t *testing.T) {
	// Compare two identical maps of BitSets
	input := []string{"valA", "valB"}
	domainMap1, _, _, _, err := GenerateBitMaps(input, []string{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("GenerateBitMaps error: %v", err)
	}
	domainMap2, _, _, _, err := GenerateBitMaps(input, []string{}, []string{}, []string{})
	if err != nil {
		t.Fatalf("GenerateBitMaps error: %v", err)
	}
	if !reflect.DeepEqual(domainMap1, domainMap2) {
		t.Errorf("Expected identical maps, but got differences")
	}
}

// New tests for NewEntry
func TestNewEntry_SuccessAndFailure(t *testing.T) {
	// Prepare maps via GenerateBitMaps
	domains := []string{"d1", "d2"}
	groups := []string{"g1"}
	names := []string{"n1", "n2", "n3"}
	values := []string{"v1", "v2"}

	domainMap, groupMap, nameMap, valueMap, err := GenerateBitMaps(domains, groups, names, values)
	if err != nil {
		t.Fatalf("GenerateBitMaps error: %v", err)
	}

	// Successful creation
	entry, err := NewEntry("d2", "g1", "n3", "v1", domainMap, groupMap, nameMap, valueMap)
	if err != nil {
		t.Errorf("NewEntry returned unexpected error: %v", err)
	}
	// Validate that the BitSets correspond to right keys
	if entry.Domain != domainMap["d2"] {
		t.Errorf("Expected entry.Domain to be domainMap[\"d2\"], got different BitSet")
	}
	if entry.Group != groupMap["g1"] {
		t.Errorf("Expected entry.Group to be groupMap[\"g1\"], got different BitSet")
	}
	if entry.Name != nameMap["n3"] {
		t.Errorf("Expected entry.Name to be nameMap[\"n3\"], got different BitSet")
	}
	if entry.Value != valueMap["v1"] {
		t.Errorf("Expected entry.Value to be valueMap[\"v1\"], got different BitSet")
	}

	// Failure cases: missing keys
	cases := []struct {
		dKey, gKey, nKey, vKey string
	}{
		{"missing", "g1", "n1", "v1"},
		{"d1", "missing", "n1", "v1"},
		{"d1", "g1", "missing", "v1"},
		{"d1", "g1", "n1", "missing"},
	}
	for _, c := range cases {
		_, err := NewEntry(c.dKey, c.gKey, c.nKey, c.vKey, domainMap, groupMap, nameMap, valueMap)
		if err == nil {
			t.Errorf("Expected error for missing key combination (%s,%s,%s,%s), got nil", c.dKey, c.gKey, c.nKey, c.vKey)
		}
	}
}
