const express = require('express');
const cors = require('cors');
const { requireJwt } = require('./middleware/auth');
const statisticsRouter = require('./routes/statistics');

function createApp({ jwtSecret, allowedOrigin = '*' }) {
  const app = express();

  app.disable('x-powered-by');
  app.use(cors({ origin: allowedOrigin }));
  app.use(express.json({ limit: '1mb' }));

  app.get('/', (req, res) =>
    res.json({
      service: 'matrix-stats-api',
      endpoints: ['GET /health', 'POST /api/statistics'],
    })
  );
  app.get('/health', (req, res) => res.json({ status: 'ok' }));
  app.use('/api/statistics', requireJwt(jwtSecret), statisticsRouter);

  app.use((req, res) => res.status(404).json({ error: 'not found' }));

  app.use((err, req, res, next) => {
    if (err.type === 'entity.parse.failed') {
      return res.status(400).json({ error: 'invalid JSON body' });
    }
    if (err.type === 'entity.too.large') {
      return res.status(413).json({ error: 'request body too large' });
    }
    console.error(err);
    return res.status(500).json({ error: 'internal server error' });
  });

  return app;
}

module.exports = { createApp };
