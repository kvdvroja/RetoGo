package matrix

import (
	"errors"
	"math"
)

const epsilon = 1e-10

var (
	ErrEmpty          = errors.New("matrix must not be empty")
	ErrNotRectangular = errors.New("all rows must have the same length")
	ErrNotFinite      = errors.New("matrix values must be finite numbers")
)

type Matrix [][]float64

func Validate(a Matrix) error {
	if len(a) == 0 || len(a[0]) == 0 {
		return ErrEmpty
	}
	cols := len(a[0])
	for _, row := range a {
		if len(row) != cols {
			return ErrNotRectangular
		}
		for _, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return ErrNotFinite
			}
		}
	}
	return nil
}

func QR(a Matrix) (Matrix, Matrix, error) {
	if err := Validate(a); err != nil {
		return nil, nil, err
	}

	m, n := len(a), len(a[0])
	r := clone(a)
	q := identity(m)

	for j := 0; j < n && j < m-1; j++ {
		for i := m - 1; i > j; i-- {
			x, y := r[i-1][j], r[i][j]
			if y == 0 {
				continue
			}
			h := math.Hypot(x, y)
			c, s := x/h, y/h

			for k := j; k < n; k++ {
				top, bottom := r[i-1][k], r[i][k]
				r[i-1][k] = c*top + s*bottom
				r[i][k] = -s*top + c*bottom
			}
			for k := 0; k < m; k++ {
				left, right := q[k][i-1], q[k][i]
				q[k][i-1] = c*left + s*right
				q[k][i] = -s*left + c*right
			}
		}
	}

	return clean(q), clean(r), nil
}

func Multiply(a, b Matrix) Matrix {
	out := zeros(len(a), len(b[0]))
	for i := range a {
		for k := range b {
			if a[i][k] == 0 {
				continue
			}
			for j := range b[0] {
				out[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return out
}

func Transpose(a Matrix) Matrix {
	out := zeros(len(a[0]), len(a))
	for i := range a {
		for j := range a[i] {
			out[j][i] = a[i][j]
		}
	}
	return out
}

func zeros(rows, cols int) Matrix {
	out := make(Matrix, rows)
	for i := range out {
		out[i] = make([]float64, cols)
	}
	return out
}

func identity(size int) Matrix {
	out := zeros(size, size)
	for i := range out {
		out[i][i] = 1
	}
	return out
}

func clone(a Matrix) Matrix {
	out := make(Matrix, len(a))
	for i, row := range a {
		out[i] = append([]float64(nil), row...)
	}
	return out
}

func clean(a Matrix) Matrix {
	for i := range a {
		for j := range a[i] {
			if math.Abs(a[i][j]) < epsilon {
				a[i][j] = 0
			}
		}
	}
	return a
}
