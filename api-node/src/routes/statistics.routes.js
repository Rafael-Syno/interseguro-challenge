'use strict';

const { Router } = require('express');
const { analyze } = require('../services/statistics.service');

// POST /api/v1/statistics con { "matrices": [{ name, data }] }.
function statisticsRouter() {
  const router = Router();

  router.post('/', (req, res, next) => {
    try {
      const body = req.body || {};
      const matrices = body.matrices ?? body.data ?? body;
      const epsilon = Number(body.epsilon) > 0 ? Number(body.epsilon) : undefined;

      res.json({
        source: 'stats-api-node',
        statistics: analyze(matrices, { epsilon })
      });
    } catch (err) {
      next(err);
    }
  });

  return router;
}

module.exports = { statisticsRouter };
