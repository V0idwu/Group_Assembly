package test

import (
	"testing"
)

// 第一种
func BenchmarkMake(t *testing.B) {
	t.ResetTimer()

	origin := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < t.N; i++ {
		target := make([]int, len(origin))
		for _, item := range origin {
			if item != 6 {
				target = append(target, item)
			}
		}
	}
}

// 第二种
func BenchmarkReuse(t *testing.B) {
	t.ResetTimer()

	origin := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < t.N; i++ {
		target := origin[:0]
		for _, item := range origin {
			if item != 6 {
				target = append(target, item)
			}
		}
	}
}

// 第三种
func BenchmarkEditOne(t *testing.B) {
	t.ResetTimer()

	origin := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < t.N; i++ {
		for i := 0; i < len(origin); i++ {
			if origin[i] == 6 {
				origin = append(origin[:i], origin[i+1:]...)
				i-- // maintain the correct index
			}
		}
	}
}
