import type { SchoolPanel } from '../../types/schoolPanel';
export type GetPublicSheetCopy = (panel: SchoolPanel | null) => {
  readonly title: string; readonly description: string;
};
