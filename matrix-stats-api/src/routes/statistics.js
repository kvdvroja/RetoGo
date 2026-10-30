const { Router } = require('express');
const { computeStatistics } = require('../services/statistics');
const { validateStatisticsRequest } = require('../validation');

const router = Router();

router.post('/', (req, res) => {
  const error = validateStatisticsRequest(req.body);
  if (error) {
    return res.status(400).json({ error });
  }
  return res.json(computeStatistics(req.body.matrices));
});

module.exports = router;
