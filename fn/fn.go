package fn

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func MapError[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) (Y, error)) (YS, error) {
	ys := make(YS, len(xs))
	for i, x := range xs {
		y, err := f(x)
		if err != nil {
			return ys, err
		}
		ys[i] = y
	}
	return ys, nil
}

func Filter[XS ~[]X, X any](xs XS, f func(X) bool) XS {
	ys := make(XS, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

func Product2[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, X1, X2 any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			ys = append(ys, f(x1, x2))
		}
	}
	return ys
}

func Product3[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, X1, X2, X3 any](xs1 XS1, xs2 XS2, xs3 XS3, f func(X1, X2, X3) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				ys = append(ys, f(x1, x2, x3))
			}
		}
	}
	return ys
}

func Product4[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, XS4 ~[]X4, X1, X2, X3, X4 any](xs1 XS1, xs2 XS2, xs3 XS3, xs4 XS4, f func(X1, X2, X3, X4) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					ys = append(ys, f(x1, x2, x3, x4))
				}
			}
		}
	}
	return ys
}

func Product5[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, XS4 ~[]X4, XS5 ~[]X5, X1, X2, X3, X4, X5 any](xs1 XS1, xs2 XS2, xs3 XS3, xs4 XS4, xs5 XS5, f func(X1, X2, X3, X4, X5) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					for _, x5 := range xs5 {
						ys = append(ys, f(x1, x2, x3, x4, x5))
					}
				}
			}
		}
	}
	return ys
}

func Product6[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, XS4 ~[]X4, XS5 ~[]X5, XS6 ~[]X6, X1, X2, X3, X4, X5, X6 any](xs1 XS1, xs2 XS2, xs3 XS3, xs4 XS4, xs5 XS5, xs6 XS6, f func(X1, X2, X3, X4, X5, X6) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					for _, x5 := range xs5 {
						for _, x6 := range xs6 {
							ys = append(ys, f(x1, x2, x3, x4, x5, x6))
						}
					}
				}
			}
		}
	}
	return ys
}

func Product7[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, XS4 ~[]X4, XS5 ~[]X5, XS6 ~[]X6, XS7 ~[]X7, X1, X2, X3, X4, X5, X6, X7 any](xs1 XS1, xs2 XS2, xs3 XS3, xs4 XS4, xs5 XS5, xs6 XS6, xs7 XS7, f func(X1, X2, X3, X4, X5, X6, X7) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					for _, x5 := range xs5 {
						for _, x6 := range xs6 {
							for _, x7 := range xs7 {
								ys = append(ys, f(x1, x2, x3, x4, x5, x6, x7))
							}
						}
					}
				}
			}
		}
	}
	return ys
}

func Product8[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, XS4 ~[]X4, XS5 ~[]X5, XS6 ~[]X6, XS7 ~[]X7, XS8 ~[]X8, X1, X2, X3, X4, X5, X6, X7, X8 any](xs1 XS1, xs2 XS2, xs3 XS3, xs4 XS4, xs5 XS5, xs6 XS6, xs7 XS7, xs8 XS8, f func(X1, X2, X3, X4, X5, X6, X7, X8) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					for _, x5 := range xs5 {
						for _, x6 := range xs6 {
							for _, x7 := range xs7 {
								for _, x8 := range xs8 {
									ys = append(ys, f(x1, x2, x3, x4, x5, x6, x7, x8))
								}
							}
						}
					}
				}
			}
		}
	}
	return ys
}

func Identity[X any](x X) X {
	return x
}