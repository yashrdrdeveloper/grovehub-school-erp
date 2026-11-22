import { getAppTestIds } from '../appLauncher.getAppTestIds';

describe('getAppTestIds', () => {
  it('uses route identities rather than display names or sensitive user input', () => {
    expect(getAppTestIds('/dashboard')).toEqual({card:'catalog-app-dashboard',pin:'catalog-pin-dashboard'});
    expect(getAppTestIds('/school/tasks')).toEqual({card:'catalog-app-school-tasks',pin:'catalog-pin-school-tasks'});
  });
});
