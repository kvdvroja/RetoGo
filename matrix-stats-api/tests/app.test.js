const request = require('supertest');
const jwt = require('jsonwebtoken');
const { createApp } = require('../src/app');

const secret = 'test-secret-with-enough-length';
const app = createApp({ jwtSecret: secret });
const token = jwt.sign({ sub: 'admin' }, secret, { algorithm: 'HS256', expiresIn: '1m' });

const validBody = {
  matrices: [
    { name: 'Q', data: [[0.6, -0.8], [0.8, 0.6]] },
    { name: 'R', data: [[5, 4], [0, 3]] },
  ],
};

describe('POST /api/statistics', () => {
  test('rejects requests without a token', async () => {
    const res = await request(app).post('/api/statistics').send(validBody);
    expect(res.status).toBe(401);
  });

  test('rejects tokens signed with another secret', async () => {
    const forged = jwt.sign({ sub: 'x' }, 'another-secret-value-123');
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${forged}`)
      .send(validBody);
    expect(res.status).toBe(401);
  });

  test('rejects expired tokens', async () => {
    const expired = jwt.sign({ sub: 'x', exp: Math.floor(Date.now() / 1000) - 10 }, secret);
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${expired}`)
      .send(validBody);
    expect(res.status).toBe(401);
  });

  test('returns statistics for valid input', async () => {
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send(validBody);

    expect(res.status).toBe(200);
    expect(res.body.max).toBe(5);
    expect(res.body.min).toBe(-0.8);
    expect(res.body.sum).toBeCloseTo(13.2);
    expect(res.body.average).toBeCloseTo(1.65);
    expect(res.body.anyDiagonal).toBe(false);
    expect(res.body.matrices).toHaveLength(2);
  });

  test.each([
    ['missing matrices', {}],
    ['empty list', { matrices: [] }],
    ['jagged matrix', { matrices: [{ name: 'A', data: [[1, 2], [3]] }] }],
    ['non numeric values', { matrices: [{ name: 'A', data: [['1']] }] }],
    ['empty matrix', { matrices: [{ name: 'A', data: [] }] }],
  ])('returns 400 for %s', async (_, body) => {
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .send(body);
    expect(res.status).toBe(400);
    expect(res.body.error).toBeDefined();
  });

  test('returns 400 for malformed JSON', async () => {
    const res = await request(app)
      .post('/api/statistics')
      .set('Authorization', `Bearer ${token}`)
      .set('Content-Type', 'application/json')
      .send('{"matrices": [');
    expect(res.status).toBe(400);
  });
});

describe('GET /health', () => {
  test('responds ok', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });
});
