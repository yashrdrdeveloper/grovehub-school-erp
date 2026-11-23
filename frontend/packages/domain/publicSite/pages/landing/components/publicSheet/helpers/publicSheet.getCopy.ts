import { panelTitles, overviewPanels } from '../constants/school.sheetCopy';
import type { SchoolPanel } from '../types/schoolPanel';

/** Pure lookup for the panel's already-translated static copy. */
export const getPublicSheetCopy = (panel: SchoolPanel | null) => {
  if (panel) return { title: panelTitles[panel], description: overviewPanels[panel] };
  return { title: __('RDR School information'), description: '' };
};
