import thunk from 'redux-thunk';
import {
  applyMiddleware, compose, createStore, combineReducers,
} from 'redux';
import {
  default as currentUserReducer,
} from 'src/features/current-user/reducer';
import {
  default as authorizationReducer,
} from 'src/features/authorization/reducer';

const createReduxStore = () => {
  const allReducers = combineReducers({
    currentUser: currentUserReducer,
    authorization: authorizationReducer,
  });

  const middleware = [thunk];

  const allStoreEnhancers = compose(
      applyMiddleware(...middleware),
  );

  return createStore(allReducers, undefined, allStoreEnhancers);
};

export default createReduxStore;
