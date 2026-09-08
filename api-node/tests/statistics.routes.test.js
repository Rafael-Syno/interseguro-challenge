'use strict';

const request = require('supertest');
const jwt = require('jsonwebtoken');
const { createApp } = require('../src/app');

const SECRET = 'test-secret';
const app = createApp({ jwtSecret: SECRET, authEnabled: true });
const token = jwt.sign({ sub: 'qr-api-go' }, SECRET, { algorithm: 'HS256', expiresIn: '5m' });

describe('POST /api/v1/statistics', () => {
  test('responde 401 sin token', async () => {
    const res = await request(app).post('/api/v1/statistics').send({ matrices: [] });
    expect(res.status).toBe(401);
  });

  test('responde 401 con token firmado con otro secreto', async () => {
    const bad = jwt.sign({ sub: 'x' }, 'otro-secreto');
    const res = await request(app)
      .post('/api/v1/statistics')
      .set('Authorization', `Bearer ${bad}`)
      .send({ matrices: [{ name: 'A', data: [[1]] }] });
    expect(res.status).toBe(401);
  });

  test('devuelve las estadísticas con un token válido', async () => {
    const res = await request(app)
      .post('/api/v1/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({
        matrices: [
          { name: 'Q', data: [[1, 0], [0, 1]] },
          { name: 'R', data: [[4, 2], [0, 8]] }
        ]
      });

    expect(res.status).toBe(200);
    expect(res.body.statistics.overall).toMatchObject({
      max: 8,
      min: 0,
      sum: 16,
      average: 2,
      anyDiagonal: true
    });
    expect(res.body.statistics.overall.diagonalMatrices).toEqual(['Q']);
  });

  test('responde 422 ante una matriz inválida', async () => {
    const res = await request(app)
      .post('/api/v1/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send({ matrices: [{ name: 'A', data: [[1, 2], [3]] }] });

    expect(res.status).toBe(422);
    expect(res.body.error.message).toMatch(/rectangular/);
  });

  test('/health no requiere autenticación', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });
});
