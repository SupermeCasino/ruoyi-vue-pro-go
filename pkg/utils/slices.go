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

// IsEqualList 比较两个切片是否相等（忽略顺序）
func IsEqualList[T comparable](list1, list2 []T) bool {
	if len(list1) != len(list2) {
		return false
	}
	m := make(map[T]int)
	for _, v := range list1 {
		m[v]++
	}
	for _, v := range list2 {
		if m[v] <= 0 {
			return false
		}
		m[v]--
	}
	return true
}
