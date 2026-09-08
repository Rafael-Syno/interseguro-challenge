'use strict';

const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const { statisticsRouter } = require('./routes/statistics.routes');
const { requireJwt } = require('./middleware/auth');

// Factory: los tests la montan sin abrir un puerto.
function createApp(options = {}) {
  const jwtSecret = options.jwtSecret || process.env.JWT_SECRET || 'interChallSecret2000';
  const authEnabled = options.authEnabled ?? process.env.AUTH_ENABLED !== 'false';

  const app = express();
  app.use(express.json({ limit: '5mb' }));
  app.use(cors());
  if (process.env.NODE_ENV !== 'test') app.use(morgan('tiny'));

  app.get('/health', (_req, res) => res.json({ status: 'ok', service: 'stats-api-node' }));

  app.use('/api/v1/statistics', requireJwt(jwtSecret, authEnabled), statisticsRouter());

  app.use((_req, res) => {
    res.status(404).json({ error: { status: 404, message: 'Recurso no encontrado' } });
  });

  // Manejador de errores unificado, con el mismo formato que la API en Go.
  app.use((err, req, res, _next) => {
    const status = err.status || 500;
    if (status >= 500) console.error(err);
    res.status(status).json({
      error: { status, message: err.message || 'Error interno', path: req.path }
    });
  });

  return app;
}

module.exports = { createApp };
