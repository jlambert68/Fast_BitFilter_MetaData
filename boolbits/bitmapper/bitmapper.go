package bitmapper

import (
	"fmt"

	"github.com/jlambert68/Fast_BitFilter_MetaData/boolbits/boolbits"
)

// GenerateBitMaps takes four string slices (domains, metadataGroupNames, metadataNames, metadataValues),
// removes duplicates in each, and assigns each unique value a BitSet with a single bit set.
// The bit length is chosen as the smallest multiple of 64 that can hold all unique values in that slice.
// It returns four maps: one per input slice, mapping each unique value to its BitSet.
func GenerateBitMaps(
	domains []string,
	metadataGroupNames []string,
	metadataNames []string,
	metadataValues []string,
) (
	map[string]*boolbits.BitSet,
	map[string]*boolbits.BitSet,
	map[string]*boolbits.BitSet,
	map[string]*boolbits.BitSet,
	error,
) {
	// Helper function to deduplicate and preserve order
	dedup := func(input []string) []string {
		seen := make(map[string]struct{})
		unique := []string{}
		for _, v := range input {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				unique = append(unique, v)
			}
		}
		return unique
	}

	// Process each slice
	uniqueDomains := dedup(domains)
	uniqueGroupNames := dedup(metadataGroupNames)
	uniqueNames := dedup(metadataNames)
	uniqueValues := dedup(metadataValues)

	// Helper to compute bit length: smallest multiple of 64 >= count
	computeBitLength := func(count int) int {
		if count <= 0 {
			return 64
		}
		// If count is already multiple of 64, use count; else round up
		if count%64 == 0 {
			return count
		}
		return ((count / 64) + 1) * 64
	}

	// Helper to assign BitSet for a list of unique values
	assign := func(uniqueList []string) (map[string]*boolbits.BitSet, error) {
		count := len(uniqueList)
		bitlen := computeBitLength(count)
		bsMap := make(map[string]*boolbits.BitSet, count)

		for idx, val := range uniqueList {
			bs, err := boolbits.NewBitSet(bitlen)
			if err != nil {
				return nil, fmt.Errorf("failed to create BitSet of length %d: %v", bitlen, err)
			}
			// Set the bit corresponding to this index
			if err := bs.SetBit(idx); err != nil {
				return nil, fmt.Errorf("failed to set bit %d for value '%s': %v", idx, val, err)
			}
			bsMap[val] = bs
		}
		return bsMap, nil
	}

	domainMap, err := assign(uniqueDomains)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	groupMap, err := assign(uniqueGroupNames)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	nameMap, err := assign(uniqueNames)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	valueMap, err := assign(uniqueValues)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return domainMap, groupMap, nameMap, valueMap, nil
}
