package omw

func StringIndexAccess[T ~string](x T, index int) T {
	i := 0
	for _, c := range x {
		if i == index {
			return T(c)
		}
		i += 1
	}
	return ""
}
