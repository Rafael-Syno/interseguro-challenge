'use strict';

const { analyze, isDiagonal } = require('../src/services/statistics.service');

describe('statistics.service', () => {
  test('calcula max, min, suma y promedio sobre todas las matrices', () => {
    const { overall } = analyze([
      { name: 'Q', data: [[1, 2], [3, 4]] },
      { name: 'R', data: [[5, -6], [0, 7]] }
    ]);

    expect(overall.max).toBe(7);
    expect(overall.min).toBe(-6);
    expect(overall.sum).toBe(16);
    expect(overall.values).toBe(8);
    expect(overall.average).toBe(2);
    expect(overall.matrices).toBe(2);
  });

  test('detecta matrices diagonales y las reporta por nombre', () => {
    const { overall, perMatrix } = analyze([
      { name: 'Q', data: [[2, 0], [0, 3]] },
      { name: 'R', data: [[1, 2], [3, 4]] }
    ]);

    expect(overall.anyDiagonal).toBe(true);
    expect(overall.diagonalMatrices).toEqual(['Q']);
    expect(perMatrix[0].isDiagonal).toBe(true);
    expect(perMatrix[1].isDiagonal).toBe(false);
  });

  test('una matriz no cuadrada nunca es diagonal', () => {
    expect(isDiagonal([[1, 0, 0], [0, 2, 0]])).toBe(false);
  });

  test('tolera el ruido de punto flotante de la factorización', () => {
    expect(isDiagonal([[1, 1e-15], [-2e-16, 3]])).toBe(true);
    expect(isDiagonal([[1, 0.001], [0, 3]])).toBe(false);
  });

  test('acepta una lista simple de matrices sin nombre', () => {
    const { perMatrix } = analyze([[[1, 2], [3, 4]]]);
    expect(perMatrix[0].name).toBe('M1');
  });

  test('rechaza matrices no rectangulares', () => {
    expect(() => analyze([{ name: 'A', data: [[1, 2], [3]] }])).toThrow(/no es rectangular/);
  });

  test('rechaza valores no numéricos', () => {
    expect(() => analyze([{ name: 'A', data: [[1, 'x']] }])).toThrow(/no numérico/);
  });

  test('rechaza una entrada vacía', () => {
    expect(() => analyze([])).toThrow(/al menos una matriz/);
  });
});
