const { createApp } = require('./app');

const port = Number(process.env.PORT) || 3000;
const jwtSecret = process.env.JWT_SECRET;

if (!jwtSecret || jwtSecret.length < 16) {
  console.error('JWT_SECRET must be set and have at least 16 characters');
  process.exit(1);
}

const app = createApp({
  jwtSecret,
  allowedOrigin: process.env.ALLOWED_ORIGIN || '*',
});

const server = app.listen(port, () => {
  console.log(`matrix-stats-api listening on port ${port}`);
});

process.on('SIGTERM', () => server.close(() => process.exit(0)));
process.on('SIGINT', () => server.close(() => process.exit(0)));
