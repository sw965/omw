package conv

func ToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}