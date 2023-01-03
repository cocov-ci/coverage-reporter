const nonProductionEnvironment = {
  NODE_ENV: process.env.NODE_ENV || 'development',
  PORT: process.env.PORT || 8080,
  LOG_LEVEL: process.env.LOG_LEVEL || 'info',
};

module.exports = process.env.NODE_ENV === 'production' ? process.env : nonProductionEnvironment;
