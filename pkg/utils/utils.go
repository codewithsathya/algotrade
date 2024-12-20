package utils

func ChunkSlice(slice []string, batchSize int) [][]string {
	var batches [][]string
    for batchSize < len(slice) {
        slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
    }
    return append(batches, slice)
}