package simdlab

import (
	"fmt"
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

	b, _ := simd.LoadFloat32sPart(xs[c:])
	acc = acc.Add(b)

	var arr [8]float32
	var r float32
	acc.Store(arr[:])

	for i, _ := range arr[:n] {
		r += arr[i]
	}

	return r
}

func SumGreaterScalar(xs []float32, t float32) float32 {
	var s float32

	for _, v := range xs {
		if v > t {
			s += v
		}
	}

	return s
}

func SumGreaterSIMD(xs []float32, t float32) float32 {
	var acc simd.Float32s
	var c int
	n := acc.Len()

	tv := simd.BroadcastFloat32s(t)

	for i := 0; i+n <= len(xs); i += n {
		a := simd.LoadFloat32s(xs[i:])
		mask := a.Greater(tv)
		acc = acc.Add(a.Masked(mask))
		c += n
	}

	b, _ := simd.LoadFloat32sPart(xs[c:])
	mask := b.Greater(tv)

	acc = acc.Add(b.Masked(mask))

	var arr [8]float32
	var r float32
	acc.Store(arr[:])

	for i, _ := range arr[:n] {
		r += arr[i]
	}

	return r
}

func BitSIMD() {
	s := []byte("hello, world!xyz")
	a := simd.LoadUint8s(s)
	fmt.Println("Len", a.Len())

	var ts byte = 'm'
	t := simd.BroadcastUint8s(ts)
	e := a.Equal(t)

	var vc = make([]int8, 16)

	i := e.ToInt8s()

	i.Store(vc)

	var m = -1
	var st string
	for i, v := range vc {
		if v < 0 {
			m = i
			break
		}
	}

	if 0 > m {
		fmt.Printf("%c は %s の中にありませんでした。", ts, s)
		return
	} else {
		st = string(s[m])
		fmt.Println(m)
		fmt.Println(st)
	}

}
