import { getStoryPanelTabClassName } from '../storyPanel.getTabClassName';

describe('getStoryPanelTabClassName', () => {
  it('marks only the selected topic', () => {
    expect(getStoryPanelTabClassName(true)).toBe('active');
    expect(getStoryPanelTabClassName(false)).toBe('');
  });
});
