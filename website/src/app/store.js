import thunk from 'redux-thunk';
import {
  applyMiddleware, compose, createStore, combineReducers,
} from 'redux';
import {
  default as currentUserReducer,
} from 'src/features/current-user/reducer';
import {
  default as authenticationReducer,
} from 'src/features/authentication/reducer';

const createReduxStore = () => {
  const allReducers = combineReducers({
    currentUser: currentUserReducer,
    authentication: authenticationReducer,
  });

  const middleware = [thunk];

  const allStoreEnhancers = compose(
      applyMiddleware(...middleware),
  );

  return createStore(allReducers, undefined, allStoreEnhancers);
};

export default createReduxStore;
