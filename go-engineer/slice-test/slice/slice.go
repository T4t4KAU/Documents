package slice

func DeleteElement[T any](s []T, index int) []T {
	if index < 0 || index >= len(s) {
		// 如果给定的索引超出范围，直接返回原切片
		return s
	}

	copy(s[index:], s[index+1:])

	// 删除元素后，将切片长度减1
	s = s[:len(s)-1]

	// 检查是否需要进行缩容
	if len(s) <= cap(s)/2 {
		// 缩容时 创建一个新的切片
		newSlice := make([]T, len(s)*2)
		copy(newSlice, s)
		s = newSlice
	}

	return s
}
