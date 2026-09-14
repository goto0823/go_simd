package simdlab

import (
	"simd"
)

func SumScalar(xs []float32) float32 {
	var s float32

	for _, i := range xs {
		s += i
	}

	return s
}

func SumSIMD(xs []float32) float32 {
	var acc simd.Float32s
	var c int
	n := acc.Len()

	for i := 0; i+n <= len(xs); i += n {
		a := simd.LoadFloat32s(xs[i:])
		acc = acc.Add(a)
		c += n
	}

	var arr [8]float32
	var r float32
	acc.Store(arr[:])

	for i, _ := range arr[:n] {
		r += arr[i]
	}

	for _, v := range xs[c:] {
		r += v
	}

	return r
}
