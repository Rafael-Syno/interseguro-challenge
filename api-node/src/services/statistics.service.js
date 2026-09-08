'use strict';

// Los ceros de una factorización numérica no son exactos.
const DEFAULT_EPSILON = 1e-9;

class ValidationError extends Error {
  constructor(message) {
    super(message);
    this.name = 'ValidationError';
    this.status = 422;
  }
}

// Acepta [{ name, data }] o una lista simple de matrices.
function normalize(matrices) {
  if (!Array.isArray(matrices) || matrices.length === 0) {
    throw new ValidationError('Se espera "matrices" como un arreglo con al menos una matriz');
  }

  return matrices.map((item, index) => {
    const name = item && !Array.isArray(item) && item.name ? String(item.name) : `M${index + 1}`;
    const data = Array.isArray(item) ? item : item && item.data;

    if (!Array.isArray(data) || data.length === 0) {
      throw new ValidationError(`La matriz "${name}" no contiene filas`);
    }

    const columns = Array.isArray(data[0]) ? data[0].length : -1;
    if (columns <= 0) {
      throw new ValidationError(`La matriz "${name}" no contiene columnas`);
    }

    data.forEach((row, i) => {
      if (!Array.isArray(row) || row.length !== columns) {
        throw new ValidationError(`La matriz "${name}" no es rectangular: la fila ${i} tiene un largo distinto`);
      }
      row.forEach((value, j) => {
        if (typeof value !== 'number' || !Number.isFinite(value)) {
          throw new ValidationError(`La matriz "${name}" tiene un valor no numérico en (${i},${j})`);
        }
      });
    });

    return { name, data };
  });
}

// Diagonal solo si es cuadrada, comparando con tolerancia.
function isDiagonal(data, epsilon = DEFAULT_EPSILON) {
  const rows = data.length;
  const columns = data[0].length;
  if (rows !== columns) return false;

  for (let i = 0; i < rows; i += 1) {
    for (let j = 0; j < columns; j += 1) {
      if (i !== j && Math.abs(data[i][j]) > epsilon) return false;
    }
  }
  return true;
}

// Estadísticas del conjunto en un solo recorrido.
function analyze(matrices, options = {}) {
  const epsilon = options.epsilon ?? DEFAULT_EPSILON;
  const decimals = options.decimals ?? 10;
  const normalized = normalize(matrices);

  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;

  const perMatrix = normalized.map(({ name, data }) => {
    let mMax = -Infinity;
    let mMin = Infinity;
    let mSum = 0;

    for (const row of data) {
      for (const value of row) {
        if (value > mMax) mMax = value;
        if (value < mMin) mMin = value;
        mSum += value;
      }
    }

    const size = data.length * data[0].length;
    max = Math.max(max, mMax);
    min = Math.min(min, mMin);
    sum += mSum;
    count += size;

    return {
      name,
      rows: data.length,
      columns: data[0].length,
      isSquare: data.length === data[0].length,
      isDiagonal: isDiagonal(data, epsilon),
      max: round(mMax, decimals),
      min: round(mMin, decimals),
      sum: round(mSum, decimals),
      average: round(mSum / size, decimals)
    };
  });

  return {
    overall: {
      matrices: normalized.length,
      values: count,
      max: round(max, decimals),
      min: round(min, decimals),
      sum: round(sum, decimals),
      average: round(sum / count, decimals),
      anyDiagonal: perMatrix.some((m) => m.isDiagonal),
      diagonalMatrices: perMatrix.filter((m) => m.isDiagonal).map((m) => m.name)
    },
    perMatrix
  };
}

function round(value, decimals) {
  const factor = 10 ** decimals;
  const rounded = Math.round(value * factor) / factor;
  return Object.is(rounded, -0) ? 0 : rounded;
}

module.exports = { analyze, isDiagonal, normalize, ValidationError, DEFAULT_EPSILON };
