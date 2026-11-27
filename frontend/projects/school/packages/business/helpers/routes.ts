import { _find } from '@Base/utilities/find';

import type { GetRoute } from './types/routes';

/** Typed lookup shared by every domain; missing configured routes fail explicitly. */
export const getRoute: GetRoute = (routes, id) => {
  const found = _find(routes, (route) => route.id === id);

  if (!found) {
    throw new Error('Unknown route id: ' + id);
  }

  return found;
};
