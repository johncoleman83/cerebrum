import axios from 'axios';
import config from 'src/config/app';

const getMe = async (authToken) => axios({
  method: 'GET',
  baseURL: config.BASE_PATH,
  url: 'me',
  headers: {
    'Authorization': `Bearer ${authToken}`,
  },
});

// TODO: consider how to use refresh later
// const refresh = async (refreshToken, authToken) => axios({
//   method: 'POST',
//   baseURL: config.BASE_PATH,
//   url: `refresh/${refreshToken}`,
//   headers: {
//     'Authorization': `Bearer ${authToken}`,
//   },
// });

export { getMe };
