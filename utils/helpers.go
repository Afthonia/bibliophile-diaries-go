package utils

func Map[T any, Z any](data []T, callback func(T) Z) []Z {
	result := []Z{}
	for _, d := range data {
		result = append(result, callback(d))
	}
	return result
}
