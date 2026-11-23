import { ChevronsLeft, ChevronsRight } from '@Components/icons';
import { getSidebarCollapseView } from '../sidebar.collapseView';
describe('getSidebarCollapseView', () => {
  it('resolves collapsed and expanded states', () => {
    expect(getSidebarCollapseView(true).CollapseIcon).toBe(ChevronsRight);
    expect(getSidebarCollapseView(true).collapsedLabel).toBe('Expand sidebar');
    expect(getSidebarCollapseView(false).CollapseIcon).toBe(ChevronsLeft);
    expect(getSidebarCollapseView(false).collapsedLabel).toBe('Collapse sidebar');
  });
});
