package shared

func Intersection[T comparable](data1 []T, data2 []T) []T {
	intersection := make([]T, 0)
	hash := make(map[T]struct{})
	for _, i := range data1 {
		hash[i] = struct{}{}
	}
	for _, j := range data2 {
		if _, ok := hash[j]; ok {
			intersection = append(intersection, j)
		}
	}
	return intersection
}

// Returns the elements in `data1` that aren't in `data2`
func Difference[T comparable](data1 []T, data2 []T) []T {
	diff := make([]T, 0)
	hash := make(map[T]struct{}, len(data2))
	for _, i := range data2 {
		hash[i] = struct{}{}
	}
	for _, j := range data1 {
		if _, ok := hash[j]; !ok {
			diff = append(diff, j)
		}
	}
	return diff
}

func Chunk[T any](slice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		return nil
	}
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
