/* eslint-disable no-param-reassign */

export const setAxiosDefaults = (axios) => {
  axios.defaults.headers.common.Accept =
    'application/vnd.pagerduty+json;version=2';
};

export default setAxiosDefaults;
