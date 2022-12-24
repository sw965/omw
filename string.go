package omw

func StringIndexAccess(x string, index int) string {
	i := 0
	for _, c := range x {
		if i == index {
			return string(c)
		}
		i += 1
	}
	return ""
}
