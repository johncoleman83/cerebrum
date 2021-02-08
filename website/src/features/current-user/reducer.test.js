import reducer, { INITIAL_STATE } from './reducer';
import * as types from './actions';

describe('current-user reducer', () => {
  it('should return the initial state', () => {
    expect(
        reducer(undefined, {}),
    ).toEqual(INITIAL_STATE);
  });

  it('should handle FETCH_ME_REQUEST', () => {
    expect(
        reducer(INITIAL_STATE, {
          type: types.FETCH_ME_REQUEST,
        }),
    ).toEqual({
      ...INITIAL_STATE,
      user: {},
      isFetching: true,
      error: null,
    });
  });

  it('should handle FETCH_ME_SUCCESS', () => {
    const payload = {
      id: '1234',
      name: 'the user',
    };

    expect(
        reducer(INITIAL_STATE, {
          type: types.FETCH_ME_SUCCESS,
          payload,
        }),
    ).toEqual({
      ...INITIAL_STATE,
      user: {
        id: payload.id,
        name: payload.name,
      },
      isFetching: false,
      error: null,
    });
  });

  it('should handle FETCH_ME_ERROR', () => {
    const payload = 'the error';

    expect(
        reducer(INITIAL_STATE, {
          type: types.FETCH_ME_ERROR,
          payload,
        }),
    ).toEqual({
      ...INITIAL_STATE,
      user: {},
      isFetching: false,
      error: payload,
    });
  });
});
