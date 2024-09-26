package utils

func Map[T any, U any](input []T, mapFunc func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = mapFunc(v)
	}
	return result
}
