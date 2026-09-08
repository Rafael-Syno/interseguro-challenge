package matrix

import (
	"math"
	"testing"
)

const tol = 1e-9

func assertReconstructs(t *testing.T, a Matrix) {
	t.Helper()
	q, r := QR(a)

	// A = Q*R
	prod, err := Mul(q, r)
	if err != nil {
		t.Fatalf("no se pudo multiplicar Q*R: %v", err)
	}
	rows, cols := a.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if math.Abs(prod[i][j]-a[i][j]) > tol {
				t.Fatalf("A != Q*R en (%d,%d): %v vs %v", i, j, prod[i][j], a[i][j])
			}
		}
	}

	// Qᵀ*Q = I
	qr, qc := q.Dims()
	if qr != rows || qc != rows {
		t.Fatalf("Q debe ser %dx%d, se obtuvo %dx%d", rows, rows, qr, qc)
	}
	for i := 0; i < qr; i++ {
		for j := 0; j < qr; j++ {
			dot := 0.0
			for k := 0; k < qr; k++ {
				dot += q[k][i] * q[k][j]
			}
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if math.Abs(dot-expected) > tol {
				t.Fatalf("Q no es ortogonal en (%d,%d): %v", i, j, dot)
			}
		}
	}

	// R triangular superior
	for i := 0; i < rows; i++ {
		for j := 0; j < i && j < cols; j++ {
			if math.Abs(r[i][j]) > tol {
				t.Fatalf("R no es triangular superior en (%d,%d): %v", i, j, r[i][j])
			}
		}
	}
}

func TestQRSquare(t *testing.T) {
	assertReconstructs(t, Matrix{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}})
}

func TestQRTallRectangular(t *testing.T) {
	assertReconstructs(t, Matrix{{1, 2}, {3, 4}, {5, 6}})
}

func TestQRWideRectangular(t *testing.T) {
	assertReconstructs(t, Matrix{{1, 2, 3}, {4, 5, 6}})
}

func TestQRSingleElement(t *testing.T) {
	assertReconstructs(t, Matrix{{5}})
}

func TestQRZeroMatrix(t *testing.T) {
	assertReconstructs(t, Matrix{{0, 0}, {0, 0}})
}

func TestQRRankDeficient(t *testing.T) {
	// Columnas linealmente dependientes: el algoritmo debe seguir siendo estable.
	assertReconstructs(t, Matrix{{1, 2}, {2, 4}, {3, 6}})
}

func TestQRKnownResult(t *testing.T) {
	// Caso clásico de Householder: R conocido para esta matriz.
	_, r := QR(Matrix{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}})
	expected := Matrix{{-14, -21, 14}, {0, -175, 70}, {0, 0, -35}}
	rows, cols := expected.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if math.Abs(r[i][j]-expected[i][j]) > 1e-6 {
				t.Fatalf("R inesperado en (%d,%d): %v, se esperaba %v", i, j, r[i][j], expected[i][j])
			}
		}
	}
}

func TestNewValidation(t *testing.T) {
	if _, err := New(nil); err != ErrEmpty {
		t.Fatalf("se esperaba ErrEmpty, se obtuvo %v", err)
	}
	if _, err := New([][]float64{{1, 2}, {3}}); err != ErrRagged {
		t.Fatalf("se esperaba ErrRagged, se obtuvo %v", err)
	}
	if _, err := New([][]float64{{math.NaN()}}); err != ErrNotFinite {
		t.Fatalf("se esperaba ErrNotFinite, se obtuvo %v", err)
	}
	if _, err := New([][]float64{{1, 2}, {3, 4}}); err != nil {
		t.Fatalf("no se esperaba error, se obtuvo %v", err)
	}
}
