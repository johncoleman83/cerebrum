import { user, userRole, isUserValid } from './selectors';
import mock from './mock';
import { camelizeKeys } from '../../util/general';

const state = {
  currentUser: {
    user: camelizeKeys(mock),
  },
};

describe('user data selectors tests', () => {
  describe('general access', () => {
    it('gets user data', () => {
      expect(user(state)).toEqual(camelizeKeys(mock));
    });
  });

  describe('roles', () => {
    it('gets user role', () => {
      expect(userRole(state)).toBe('SOME_ROLE');
    });
  });

  describe('other checks', () => {
    it('is true if the user data is valid', () => {
      expect(isUserValid(state)).toBeTruthy();
    });

    it('is false if the user data is not valid', () => {
      expect(isUserValid({ currentUser: { user: {} } })).toBeFalsy();
    });
  });
});
