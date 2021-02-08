import axios from 'axios';
import config from 'src/config/app';

const loginPostApiCall = async (username, password) => axios({
  method: 'POST',
  baseURL: config.BASE_PATH,
  url: 'login',
  data: {
    username: username,
    password: password,
  },
});

export { loginPostApiCall };
