package matrix

import (
	"errors"
	"math"
	"testing"
)

const tolerance = 1e-9

func assertClose(t *testing.T, got, want Matrix) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("row count mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if len(got[i]) != len(want[i]) {
			t.Fatalf("column count mismatch on row %d", i)
		}
		for j := range want[i] {
			if math.Abs(got[i][j]-want[i][j]) > tolerance {
				t.Fatalf("value mismatch at [%d][%d]: got %v, want %v", i, j, got[i][j], want[i][j])
			}
		}
	}
}

func TestQRReconstructsInput(t *testing.T) {
	cases := map[string]Matrix{
		"square":     {{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}},
		"tall":       {{1, 2}, {3, 4}, {5, 6}},
		"wide":       {{1, 2, 3}, {4, 5, 6}},
		"single row": {{3, 4, 5}},
		"single col": {{3}, {4}},
		"zeros":      {{0, 0}, {0, 0}},
	}

	for name, a := range cases {
		t.Run(name, func(t *testing.T) {
			q, r, err := QR(a)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertClose(t, Multiply(q, r), a)
			assertClose(t, Multiply(Transpose(q), q), identity(len(a)))

			for i := range r {
				for j := 0; j < i && j < len(r[i]); j++ {
					if r[i][j] != 0 {
						t.Fatalf("R is not upper triangular at [%d][%d]: %v", i, j, r[i][j])
					}
				}
			}
		})
	}
}

func TestQRKnownResult(t *testing.T) {
	_, r, err := QR(Matrix{{3, 0}, {4, 5}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertClose(t, r, Matrix{{5, 4}, {0, 3}})
}

func TestQRDoesNotMutateInput(t *testing.T) {
	a := Matrix{{1, 2}, {3, 4}}
	if _, _, err := QR(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertClose(t, a, Matrix{{1, 2}, {3, 4}})
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name  string
		input Matrix
		want  error
	}{
		{"nil", nil, ErrEmpty},
		{"empty row", Matrix{{}}, ErrEmpty},
		{"jagged", Matrix{{1, 2}, {3}}, ErrNotRectangular},
		{"nan", Matrix{{1, math.NaN()}}, ErrNotFinite},
		{"inf", Matrix{{math.Inf(1)}}, ErrNotFinite},
		{"valid", Matrix{{1}}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := Validate(tc.input); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}
