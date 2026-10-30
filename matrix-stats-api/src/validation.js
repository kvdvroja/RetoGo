const MAX_MATRICES = 10;
const MAX_DIMENSION = 100;

function validateMatrix(data, label) {
  if (!Array.isArray(data) || data.length === 0) {
    return `${label} must be a non-empty array of rows`;
  }
  if (data.length > MAX_DIMENSION) {
    return `${label} exceeds ${MAX_DIMENSION} rows`;
  }

  const cols = Array.isArray(data[0]) ? data[0].length : 0;
  if (cols === 0) {
    return `${label} rows must be non-empty arrays`;
  }
  if (cols > MAX_DIMENSION) {
    return `${label} exceeds ${MAX_DIMENSION} columns`;
  }

  for (const row of data) {
    if (!Array.isArray(row) || row.length !== cols) {
      return `${label} must be rectangular`;
    }
    if (!row.every((v) => typeof v === 'number' && Number.isFinite(v))) {
      return `${label} must contain only finite numbers`;
    }
  }
  return null;
}

function validateStatisticsRequest(body) {
  const matrices = body && body.matrices;
  if (!Array.isArray(matrices) || matrices.length === 0) {
    return '"matrices" must be a non-empty array';
  }
  if (matrices.length > MAX_MATRICES) {
    return `no more than ${MAX_MATRICES} matrices are allowed`;
  }

  for (let i = 0; i < matrices.length; i++) {
    const item = matrices[i];
    if (!item || typeof item !== 'object') {
      return `matrices[${i}] must be an object`;
    }
    const label = typeof item.name === 'string' && item.name ? item.name : `matrices[${i}]`;
    const error = validateMatrix(item.data, label);
    if (error) return error;
  }
  return null;
}

module.exports = { validateStatisticsRequest };
