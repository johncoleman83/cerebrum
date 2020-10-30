import { getMe } from './util';

export const FETCH_ME_REQUEST = 'fetchMeActionRequest';
export const FETCH_ME_SUCCESS = 'fetchMeActionSuccess';
export const FETCH_ME_ERROR = 'fetchMeActionError';

function fetchMeRequest() {
  return {
    type: FETCH_ME_REQUEST,
  };
}

function fetchMeSuccess(user) {
  return {
    type: FETCH_ME_SUCCESS,
    payload: user,
  };
}

function fetchMeError(error) {
  return {
    type: FETCH_ME_ERROR,
    payload: error,
  };
}

function fetchMeAction(authToken) {
  return async (dispatch) => {
    try {
      dispatch(fetchMeRequest());

      const { data: user } = await getMe(authToken);

      dispatch(fetchMeSuccess(user));
    } catch (error) {
      console.error(error);
      dispatch(fetchMeError(error));
    }
  };
}

export { fetchMeAction };
