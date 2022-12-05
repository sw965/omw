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
