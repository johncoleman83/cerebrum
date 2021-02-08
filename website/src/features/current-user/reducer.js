import {
  FETCH_ME_REQUEST,
  FETCH_ME_SUCCESS,
  FETCH_ME_ERROR,
} from './actions';
import { camelizeKeys } from '../../util/general';

export const INITIAL_STATE = Object.freeze({
  user: {},
  isFetching: false,
  error: null,
});

/**
 * @param {any} state react state.
 * @param {any} action react state.
 * @return {any} modified state
 */

export default function currentUserReducer(
    state = INITIAL_STATE,
    action,
) {
  switch (action.type) {
    case FETCH_ME_REQUEST:
      return {
        ...state,
        error: null,
        isFetching: true,
        user: {},
      };
    case FETCH_ME_SUCCESS:
      return {
        ...state,
        isFetching: false,
        user: camelizeKeys(action.payload),
        error: null,
      };
    case FETCH_ME_ERROR:
      return {
        ...state,
        error: action.payload,
        isFetching: false,
        user: {},
      };
    default:
      return state;
  }
}
