package simdlab

import (
	"testing"
)

func TestSumScalar(t *testing.T) {
	tests := []struct {
		name  string
		slice []float32
		want  float32
	}{
		{
			name:  "通常",
			slice: []float32{1, 2, 3, 4, 5},
			want:  15,
		},
		{
			name:  "空",
			slice: []float32{},
			want:  0,
		},
		{
			name:  "nil",
			slice: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := SumSIMD
			t.Logf("f の型は: %T", f)
			got := f(tt.slice)
			if tt.want != got {
				t.Errorf("want: %g, got: %g", tt.want, got)
			}
		})
	}

}

func BenchmarkSum(b *testing.B) {
	tests := []struct {
		name  string
		slice []float32
	}{
		{
			name:  "1000",
			slice: make([]float32, 1000),
		},
		{
			name:  "10000",
			slice: make([]float32, 10000),
		},
		{
			name:  "100000",
			slice: make([]float32, 100000),
		},
		{
			name:  "1000000",
			slice: make([]float32, 1000000),
		},
		{
			name:  "10000000",
			slice: make([]float32, 10000000),
		},
	}

	impls := []struct {
		name string
		fn   func([]float32) float32
	}{
		{
			name: "scalar",
			fn:   SumScalar,
		},
		{
			name: "simd",
			fn:   SumSIMD,
		},
	}

	for _, tt := range tests {
		for i := range tt.slice {
			tt.slice[i] = float32(i)
		}
	}

	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			for _, tt := range tests {
				b.Run(tt.name, func(b *testing.B) {
					for b.Loop() {
						impl.fn(tt.slice)
					}
				})
			}
		})
	}
}
