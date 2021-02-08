import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import * as actions from './actions';
import mockData from './mock';
import { getMe as getMeMock } from './util';

jest.mock('./util.js');

const middlewares = [thunk];
const mockStore = configureMockStore(middlewares);

describe('fetchMeAction', () => {
  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('should handle api error', () => {
    const errorPayload = { error: 'fail!' };
    global.console = { error: jest.fn() };
    getMeMock.mockRejectedValue(errorPayload);

    const expectedActions = [
      { type: actions.FETCH_ME_REQUEST },
      { type: actions.FETCH_ME_ERROR, payload: errorPayload },
    ];
    const store = mockStore({});

    return store.dispatch(actions.fetchMeAction('FakeToken')).then(() => {
      expect(store.getActions()).toEqual(expectedActions);
      expect(console.error).toHaveBeenCalledTimes(1);
      expect(console.error).toHaveBeenCalledWith(errorPayload);
    });
  });

  it('should handle expected response', () => {
    getMeMock.mockResolvedValue({
      data: mockData,
    });

    const expectedActions = [
      { type: actions.FETCH_ME_REQUEST },
      {
        type: actions.FETCH_ME_SUCCESS,
        payload: mockData,
      },
    ];
    const store = mockStore({});

    return store.dispatch(actions.fetchMeAction('FakeToken')).then(() => {
      expect(store.getActions()).toEqual(expectedActions);
    });
  });
});
