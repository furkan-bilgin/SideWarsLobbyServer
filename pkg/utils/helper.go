package utils

// https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}
