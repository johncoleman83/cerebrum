import { loginPostApiCall } from './util';

export const LOGIN_REQUEST = 'loginActionRequest';
export const LOGIN_SUCCESS = 'loginActionSuccess';
export const LOGIN_ERROR = 'loginActionError';
export const LOGOUT = 'logoutAction';

function loginRequest() {
  return {
    type: LOGIN_REQUEST,
  };
}

function loginSuccess(payload) {
  return {
    type: LOGIN_SUCCESS,
    payload: payload,
  };
}

function loginError(error) {
  return {
    type: LOGIN_ERROR,
    payload: error,
  };
}

function logout() {
  return {
    type: LOGOUT,
  };
}

/**
 * @param {string} email something
 * @param {string} password something
 * @return {any} a promise.
 */
function loginAction(email, password) {
  return async (dispatch) => {
    try {
      dispatch(loginRequest());

      const { data: payload } = await loginPostApiCall(email, password);

      dispatch(loginSuccess(payload));
    } catch (error) {
      console.error(error);
      dispatch(loginError({ response: error.response }));
    }
  };
}

/**
 * @return {any}
 */
function logoutAction() {
  return async (dispatch) => {
    dispatch(logout());
  };
}

export { loginAction, logoutAction };
