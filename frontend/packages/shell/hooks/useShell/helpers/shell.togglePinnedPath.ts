import { _filter } from '@Base/utilities/filter';
import { _includes } from '@Base/utilities/includes';
import type { TogglePinnedPath } from './types/shell.togglePinnedPath';

/** Pure, order-preserving pin toggle without duplicate entries. */
export const togglePinnedPath: TogglePinnedPath = (pinned, path) => {
  if (_includes(pinned, path)) return _filter(pinned, (item) => item !== path);
  return [...pinned, path];
};
