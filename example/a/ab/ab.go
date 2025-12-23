package ab

// Sum adds two numbers of a generic type T which can be either int or float64.
func Sum[T int | float64](a, b T) (r T) {
	r = a + b
	return
}

// SumAll adds a variadic number of values of a generic type T which can be either int or float64.
func SumAll[T int | float64](values ...T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}
