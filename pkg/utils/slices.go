package utils

// Intersect 返回两个切片的交集
func Intersect[T comparable](slice1, slice2 []T) []T {
	m := make(map[T]bool)
	for _, v := range slice1 {
		m[v] = true
	}
	var intersect []T
	for _, v := range slice2 {
		if m[v] {
			intersect = append(intersect, v)
		}
	}
	return intersect
}
