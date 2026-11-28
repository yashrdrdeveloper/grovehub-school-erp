import type { MouseEventHandler, ReactNode } from 'react';
import type { SchoolPanel } from '../../../components/publicSheet';

export interface LandingPageModel {
  readonly CurrentPage: () => ReactNode;
  readonly panel: SchoolPanel | null;
  readonly closePanel: () => void;
  readonly selectPanel: (panel: SchoolPanel) => void;
  readonly handleInternalLink: MouseEventHandler<HTMLDivElement>;
  readonly openEnquiry: () => void;
}
export type UseLandingPage = () => LandingPageModel;
