import type { StoryKey } from '../constants/school.sheetCopy';
import type { SchoolPanel } from '../types/schoolPanel';

/** Only story panels use their route-specific key; other panels default safely. */
export const getPublicSheetStoryKey = (panel: SchoolPanel | null): StoryKey => {
  if (
    panel === 'director' ||
    panel === 'principal' ||
    panel === 'philosophy' ||
    panel === 'history'
  )
    return panel;
  return 'history';
};
