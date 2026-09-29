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

func TestSumGreater(t *testing.T) {
	tests := []struct {
		name      string
		slice     []float32
		threshold float32
		want      float32
	}{
		{
			name:      "等しい数は足さない",
			slice:     []float32{1, 4, 3, 8, 5},
			threshold: 4,
			want:      13,
		},
		{
			name:      "nilのスライス",
			slice:     nil,
			threshold: 4,
			want:      0,
		},
		{
			name:      "端数の存在するケース",
			slice:     []float32{1, 2, 3, 4, 5},
			threshold: 3,
			want:      9,
		},
	}

	impls := []struct {
		name string
		fn   func([]float32, float32) float32
	}{
		{
			name: "Scalar",
			fn:   SumGreaterScalar,
		},
		{
			name: "Simd",
			fn:   SumGreaterSIMD,
		},
	}

	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					f := impl.fn
					got := f(tt.slice, tt.threshold)

					if tt.want != got {
						t.Errorf("want: %g, got: %g", tt.want, got)
					}
				})
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
