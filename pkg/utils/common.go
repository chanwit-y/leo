package utils

func ternary[T any](condition bool, v1 T, v2 T) T {
	if condition {
		return v1
	} else {
		return v2
	}

}

func space(len int) string {
	space := ""
	for i := 0; i <= len; i++ {
		space += " "
	}
	return space
}
