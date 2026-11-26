import type { MobileTableToolbarProps } from '../../../types/mobileTableToolbar.props';
import type { MobileTablePanel } from '../../../constants/mobileTableToolbar.panel';

export type MobileTableOptionsProps<T> = MobileTableToolbarProps<T> & {
  readonly panel: MobileTablePanel | null;
  readonly setPanel: (panel: MobileTablePanel) => void;
  readonly close: () => void;
  readonly batchLabel: string;
};
