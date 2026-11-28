import { getStoryPanelPhotoIndex } from '../storyPanel.getPhotoIndex';

describe('getStoryPanelPhotoIndex', () => {
  it('preserves historical photos for the corresponding story', () => {
    expect(getStoryPanelPhotoIndex('history')).toBe(4);
    expect(getStoryPanelPhotoIndex('philosophy')).toBe(2);
    expect(getStoryPanelPhotoIndex('principal')).toBe(0);
  });
});
