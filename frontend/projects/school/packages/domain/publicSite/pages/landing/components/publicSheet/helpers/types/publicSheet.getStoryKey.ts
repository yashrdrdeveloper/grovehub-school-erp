import type { StoryKey } from '../../constants/school.sheetCopy';
import type { SchoolPanel } from '../../types/schoolPanel';
export type GetPublicSheetStoryKey = (panel: SchoolPanel | null) => StoryKey;
