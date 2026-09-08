'use strict';

const jwt = require('jsonwebtoken');

// Valida el JWT propagado por la API en Go (HS256).
function requireJwt(secret, enabled = true) {
  return (req, res, next) => {
    if (!enabled) return next();

    const header = req.headers.authorization || '';
    if (!header.startsWith('Bearer ')) {
      return res.status(401).json({
        error: { status: 401, message: 'Falta el header Authorization: Bearer <token>' }
      });
    }

    try {
      req.auth = jwt.verify(header.slice(7).trim(), secret, { algorithms: ['HS256'] });
      return next();
    } catch (err) {
      return res.status(401).json({
        error: { status: 401, message: 'Token inválido o expirado' }
      });
    }
  };
}

module.exports = { requireJwt };
