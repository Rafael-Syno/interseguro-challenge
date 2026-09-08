package matrix

import "math"

// epsilon es el umbral por debajo del cual una norma se considera nula.
const epsilon = 1e-12

// QR completa por Householder: estable y admite m < n.
func QR(a Matrix) (q Matrix, r Matrix) {
	rows, cols := a.Dims()
	q = Identity(rows)
	r = a.Clone()

	steps := cols
	if rows-1 < steps {
		steps = rows - 1
	}

	for k := 0; k < steps; k++ {
		norm := 0.0
		for i := k; i < rows; i++ {
			norm += r[i][k] * r[i][k]
		}
		norm = math.Sqrt(norm)
		if norm < epsilon {
			continue // la columna ya está anulada, no hace falta reflector
		}

		// alpha con signo opuesto al pivote para evitar cancelación catastrófica.
		alpha := -norm
		if r[k][k] < 0 {
			alpha = norm
		}

		// v = x - alpha*e1, normalizado.
		v := make([]float64, rows-k)
		for i := k; i < rows; i++ {
			v[i-k] = r[i][k]
		}
		v[0] -= alpha

		vNorm := 0.0
		for _, x := range v {
			vNorm += x * x
		}
		vNorm = math.Sqrt(vNorm)
		if vNorm < epsilon {
			continue
		}
		for i := range v {
			v[i] /= vNorm
		}

		// R := H*R, con H = I - 2vvᵀ aplicado solo al bloque activo.
		for j := k; j < cols; j++ {
			dot := 0.0
			for i := k; i < rows; i++ {
				dot += v[i-k] * r[i][j]
			}
			dot *= 2
			for i := k; i < rows; i++ {
				r[i][j] -= dot * v[i-k]
			}
		}

		// Q := Q*H (acumula el producto de reflectores).
		for i := 0; i < rows; i++ {
			dot := 0.0
			for j := k; j < rows; j++ {
				dot += q[i][j] * v[j-k]
			}
			dot *= 2
			for j := k; j < rows; j++ {
				q[i][j] -= dot * v[j-k]
			}
		}
	}

	// Fuerza ceros exactos bajo la diagonal.
	for i := 0; i < rows; i++ {
		for j := 0; j < i && j < cols; j++ {
			r[i][j] = 0
		}
	}

	return q, r
}
