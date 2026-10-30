const { computeStatistics, isDiagonal } = require('../src/services/statistics');

describe('isDiagonal', () => {
  test('detects square diagonal matrices', () => {
    expect(isDiagonal([[1, 0], [0, 5]])).toBe(true);
    expect(isDiagonal([[1, 2], [0, 5]])).toBe(false);
  });

  test('handles rectangular matrices', () => {
    expect(isDiagonal([[3, 0, 0], [0, 4, 0]])).toBe(true);
    expect(isDiagonal([[3, 0], [0, 4], [0, 1]])).toBe(false);
  });

  test('ignores floating point noise', () => {
    expect(isDiagonal([[1, 1e-14], [-1e-13, 1]])).toBe(true);
  });
});

describe('computeStatistics', () => {
  test('aggregates values across all matrices', () => {
    const result = computeStatistics([
      { name: 'Q', data: [[1, 2], [3, 4]] },
      { name: 'R', data: [[-5, 0], [0, 10]] },
    ]);

    expect(result.max).toBe(10);
    expect(result.min).toBe(-5);
    expect(result.sum).toBe(15);
    expect(result.count).toBe(8);
    expect(result.average).toBeCloseTo(1.875);
    expect(result.anyDiagonal).toBe(true);
    expect(result.matrices).toEqual([
      { name: 'Q', rows: 2, cols: 2, isDiagonal: false },
      { name: 'R', rows: 2, cols: 2, isDiagonal: true },
    ]);
  });

  test('reports no diagonal matrices when none qualify', () => {
    const result = computeStatistics([{ name: 'A', data: [[1, 1]] }]);
    expect(result.anyDiagonal).toBe(false);
  });
});
