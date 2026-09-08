'use strict';

const { createApp } = require('./app');

const port = Number(process.env.PORT || 3000);
const app = createApp();

const server = app.listen(port, () => {
  console.log(`stats-api-node escuchando en :${port}`);
});

// Apagado ordenado en contenedores.
for (const signal of ['SIGINT', 'SIGTERM']) {
  process.on(signal, () => {
    console.log('cerrando el servidor...');
    server.close(() => process.exit(0));
  });
}
