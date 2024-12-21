package utils

import "strings"

func ChunkSlice(slice []string, batchSize int) [][]string {
	var batches [][]string
    for batchSize < len(slice) {
        slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
    }
    return append(batches, slice)
}

func GetSymbol(stream string) string {
	parts := strings.Split(stream, "@")
	return strings.ToUpper(parts[0])
}

func ShiftToStart(slice []string, target string) []string {
	var targetIndex int
	found := false
	for i, v := range slice {
		if v == target {
			targetIndex = i
			found = true
			break
		}
	}
	if !found {
		return slice
	}
	return append(slice[targetIndex:], slice[:targetIndex]...)
}

func ShiftMiddle(slice []string) []string {
	if len(slice) <= 2 {
		return slice
	}

	start := slice[0]
	end := slice[len(slice)-1]
	middle := slice[1 : len(slice)-1]

	firstMiddleElement := middle[0]
	middle = append(middle[1:], firstMiddleElement)

	return append(append([]string{start}, middle...), end)
}