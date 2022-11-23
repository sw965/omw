package omw

func All(x ...bool) bool {
	for _, ele := range x {
		if !ele {
			return false
		}
	}
	return true
}

func Any(x ...bool) bool {
	for _, ele := range x {
		if ele {
			return true
		}
	}
	return false
}

func IndexAccessString(str string, index int) string {
	i := 0
	for _, c := range str {
		if i == index {
			return string(c)
		}
		i += 1
	}
	return ""
}