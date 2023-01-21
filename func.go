package omw

func Mapping[X, Y any](xs []X, f func(X) Y) []Y {
	result := make([]Y, len(xs))
	for i, x := range xs {
		result[i] = f(x)
	}
	return result
}

func MappingErr[X, Y any](xs []X, f func(X) (Y, error)) ([]Y, error) {
	result := make([]Y, 0, len(xs))
	for _, x := range xs {
		y, err := f(x)
		if err != nil {
			return result, err
		}
		result = append(result, y)
	}
	return result, nil
}

func Filter[X any](xs []X, f func(X) bool) []X {
	result := make([]X, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			result = append(result, x)
		}
	}
	return result
}

func FilterErr[X any](xs []X, f func(X) (bool, error)) ([]X, error) {
	result := make([]X, 0, len(xs))
	for _, x := range xs {
		ok, err := f(x)
		if err != nil {
			return result, err
		}

		if ok {
			result = append(result, x)
		}
	}
	return result, nil
}

func Equals[T comparable] (xs1, xs2 []T) bool {
	for i, x1 := range  xs1 {
		x2 := xs2[i]
		if x1 != x2 {
			return false 
		}
	}
	return true
}

func Contains[T comparable] (xs []T, x T) bool {
	for _, ele := range xs {
		if ele == x {
			return true
		}
	}
	return false
}

func AllContains[T comparable] (xs, eles []T) bool {
	for _, ele := range eles {
		if !Contains(xs, ele) {
			return false
		}
	}
	return true
}

func Count[T comparable] (xs []T, x T) int {
	result := 0
	for _, ele := range xs {
		if ele == x {
			result += 1
		}
	}
	return result
}

func AllCount[T comparable] (xs []T) map[T]int {
	result := map[T]int{}
	for _, x := range xs {
		_, ok := result[x]
		if ok {
			result[x] += 1
		} else {
			result[x] = 1
		}
	}
	return result
}

func Indices[T comparable](xs []T, f func(T) bool) []int {
	result := make([]int, 0, len(xs))
	for i, x := range xs {
		if f(x) {
			result = append(result, i)
		}
	}
	return result
}

func IsUnique[T comparable] (xs []T) bool {
	for _, x := range xs {
		if Count(xs, x) != 1 {
			return false
		}
	}
	return true
}