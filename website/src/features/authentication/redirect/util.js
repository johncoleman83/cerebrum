import { redirectTo } from '../../../util/general';

const redirectToLogin = () => {
  // const path = window.location.pathname;
  const loginUrl = '/login';
  redirectTo(loginUrl);
};

export { redirectToLogin };
