// Package matrix: matrices densas y álgebra lineal.
package matrix

import (
	"errors"
	"math"
)

// Matrix es una matriz densa rectangular almacenada por filas.
type Matrix [][]float64

var (
	ErrEmpty     = errors.New("la matriz debe tener al menos una fila y una columna")
	ErrRagged    = errors.New("todas las filas deben tener la misma cantidad de columnas")
	ErrNotFinite = errors.New("la matriz contiene valores no finitos (NaN o Inf)")
)

func New(data [][]float64) (Matrix, error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil, ErrEmpty
	}
	cols := len(data[0])
	for _, row := range data {
		if len(row) != cols {
			return nil, ErrRagged
		}
		for _, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return nil, ErrNotFinite
			}
		}
	}
	return Matrix(data), nil
}

func (m Matrix) Dims() (int, int) {
	if len(m) == 0 {
		return 0, 0
	}
	return len(m), len(m[0])
}

func (m Matrix) Clone() Matrix {
	rows, cols := m.Dims()
	out := Zeros(rows, cols)
	for i := 0; i < rows; i++ {
		copy(out[i], m[i])
	}
	return out
}

func Zeros(rows, cols int) Matrix {
	out := make(Matrix, rows)
	for i := range out {
		out[i] = make([]float64, cols)
	}
	return out
}

func Identity(n int) Matrix {
	out := Zeros(n, n)
	for i := 0; i < n; i++ {
		out[i][i] = 1
	}
	return out
}

func Mul(a, b Matrix) (Matrix, error) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ac != br {
		return nil, errors.New("dimensiones incompatibles para la multiplicación")
	}
	out := Zeros(ar, bc)
	for i := 0; i < ar; i++ {
		for k := 0; k < ac; k++ {
			aik := a[i][k]
			if aik == 0 {
				continue
			}
			for j := 0; j < bc; j++ {
				out[i][j] += aik * b[k][j]
			}
		}
	}
	return out, nil
}

// Round elimina el ruido de punto flotante antes de serializar.
func (m Matrix) Round(decimals int) Matrix {
	factor := math.Pow(10, float64(decimals))
	rows, cols := m.Dims()
	out := Zeros(rows, cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := math.Round(m[i][j]*factor) / factor
			// Evita imprimir "-0".
			if v == 0 {
				v = 0
			}
			out[i][j] = v
		}
	}
	return out
}
