const getEnvironment = () => {
  // "development" for `react-scripts start`
  // "production" for `react-scripts build`
  // "test" for `react-scripts test`
  const env = process.env.NODE_ENV;

  return env;
};

export default getEnvironment;
