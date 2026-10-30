const TOLERANCE = 1e-10;

function isDiagonal(matrix) {
  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < matrix[i].length; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > TOLERANCE) {
        return false;
      }
    }
  }
  return true;
}

function computeStatistics(matrices) {
  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;

  const summaries = matrices.map(({ name, data }) => {
    for (const row of data) {
      for (const value of row) {
        if (value > max) max = value;
        if (value < min) min = value;
        sum += value;
        count++;
      }
    }

    return {
      name,
      rows: data.length,
      cols: data[0].length,
      isDiagonal: isDiagonal(data),
    };
  });

  return {
    max,
    min,
    average: sum / count,
    sum,
    count,
    anyDiagonal: summaries.some((m) => m.isDiagonal),
    matrices: summaries,
  };
}

module.exports = { computeStatistics, isDiagonal };
