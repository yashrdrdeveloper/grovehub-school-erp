import type { GetStoryPanelPhotoIndex } from './types/storyPanel.getPhotoIndex';
export const getStoryPanelPhotoIndex: GetStoryPanelPhotoIndex = (selected) => {
  if (selected === 'history') return 4;
  if (selected === 'philosophy') return 2;
  return 0;
};
