import type { GetStoryPanelTabClassName } from './types/storyPanel.getTabClassName';
export const getStoryPanelTabClassName: GetStoryPanelTabClassName = (active) => {
  if (active) return 'active';
  return '';
};
