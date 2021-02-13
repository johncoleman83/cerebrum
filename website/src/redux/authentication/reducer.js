import {
  LOGIN_REQUEST,
  LOGIN_SUCCESS,
  LOGIN_ERROR,
  LOGOUT,
} from './actions';

export const INITIAL_STATE = Object.freeze({
  isFetching: false,
  authToken: '',
  error: null,
});

/**
 * @param {any} state react state.
 * @param {any} action react state.
 * @return {any} modified state
 */

export default function loginReducer(
    state = INITIAL_STATE,
    action,
) {
  switch (action.type) {
    case LOGIN_REQUEST:
      return {
        ...state,
        error: null,
        isFetching: true,
        authToken: '',
      };
    case LOGIN_SUCCESS:
      return {
        ...state,
        authToken: action.payload.token,
        isFetching: false,
        error: null,
      };
    case LOGIN_ERROR:
      return {
        ...state,
        error: action.payload,
        isFetching: false,
        authToken: '',
      };
    case LOGOUT:
      return {
        ...state,
        error: null,
        isFetching: false,
        authToken: '',
      };
    default:
      return state;
  }
}
