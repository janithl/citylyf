package utils

func ConvertToF64(s []int) []float64 {
	f64 := make([]float64, len(s))
	for i, val := range s {
		f64[i] = float64(val)
	}
	return f64
}
