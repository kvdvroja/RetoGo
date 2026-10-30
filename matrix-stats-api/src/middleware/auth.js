const jwt = require('jsonwebtoken');

function requireJwt(secret) {
  return (req, res, next) => {
    const header = req.get('Authorization') || '';
    const [scheme, token] = header.split(' ');

    if (scheme !== 'Bearer' || !token) {
      return res.status(401).json({ error: 'missing bearer token' });
    }

    try {
      req.user = jwt.verify(token, secret, { algorithms: ['HS256'] });
      return next();
    } catch {
      return res.status(401).json({ error: 'invalid or expired token' });
    }
  };
}

module.exports = { requireJwt };
