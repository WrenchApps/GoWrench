package cross_validation

import (
	"wrench/app/manifest"
)

func toHasIdSlice[T manifest.HasId](items []T) []manifest.HasId {
	result := make([]manifest.HasId, len(items))
	for i, item := range items {
		result[i] = item
	}
	return result
}

func duplicateIdsValid(ids []manifest.HasId) []string {
	seen := make(map[string]bool)
	dupes := make(map[string]bool)
	for _, e := range ids {
		id := e.GetId()
		if seen[id] {
			dupes[id] = true
		}
		seen[id] = true
	}

	var result []string
	for id := range dupes {
		result = append(result, id)
	}
	return result
}
