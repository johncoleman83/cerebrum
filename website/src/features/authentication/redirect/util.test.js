import { redirectTo } from '../../../util/general';
import getEnvironment from '../../../util/environment';
import { redirectToLogin } from './util';

jest.mock('../../../util/general.js');
jest.mock('../../../util/environment.js');
jest.mock('../../../config/app.js');

describe('redirectToLogin', () => {
  beforeEach(() => {
    jest.resetAllMocks();
    getEnvironment.mockReturnValue('production');
  });

  it('should call redirectTo, passing in correct url', () => {
    redirectToLogin();
    expect(redirectTo).toHaveBeenCalledWith('/login');
  });
});
