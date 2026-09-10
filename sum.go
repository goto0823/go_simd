package simdlab

func SumScalar(xs []float32) float32 {
	var s float32

	for _, i := range xs {
		s += i
	}

	return s
}
