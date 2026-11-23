import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box, BOX_TAG } from '@Components/components/Box';
import type { ReactNode } from 'react';
import {
  DialogRoot,
  DialogPortal,
  DialogOverlay,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogClose,
} from '@Components/components/Dialog';
import { ArrowRight, GraduationCap, Phone, ShieldCheck, X } from '@Components/icons';
import { school } from '../../constants/school.content';
import { panelTitles, overviewPanels, type StoryKey } from './constants/school.sheetCopy';
import { StoryPanel } from './components/storyPanel';
import {
  AcademicPanel,
  CampusPanel,
  GalleryPanel,
  AdmissionsPanel,
  ResourcesPanel,
  ContactPanel,
} from './components/infoPanels';
import type { SchoolPanel } from './types/schoolPanel';

export type { SchoolPanel } from './types/schoolPanel';

export const PublicSheet = ({
  panel,
  onClose,
  onNavigate,
  enquiry,
}: {
  panel: SchoolPanel | null;
  onClose: () => void;
  onNavigate: (panel: SchoolPanel) => void;
  enquiry: ReactNode;
}) => {
  let selectedStory: StoryKey = 'history';
  if (
    panel === 'director' ||
    panel === 'principal' ||
    panel === 'philosophy' ||
    panel === 'history'
  )
    selectedStory = panel;
  const story =
    panel === 'about' ||
    panel === 'history' ||
    panel === 'director' ||
    panel === 'principal' ||
    panel === 'philosophy';
  const onEnquire = () => onNavigate('enquiry');
  let title = 'RDR School information';
  let description = '';
  if (panel) {
    title = panelTitles[panel];
    description = overviewPanels[panel];
  }
  return (
    <DialogRoot
      open={Boolean(panel)}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <DialogPortal>
        <DialogOverlay className="rdr-sheet-backdrop" />
        <DialogContent className="rdr-sheet" aria-label={title}>
          <Box as={BOX_TAG.DIV} className="rdr-sheet__handle" aria-hidden="true" />
          <Box as={BOX_TAG.DIV} className="rdr-sheet__header">
            <Box as={BOX_TAG.DIV}>
              <Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">
                {__('RAO DEENA RAM VIDYA VIHAR')}
              </Box>
              <DialogTitle>{title}</DialogTitle>
              <DialogDescription>{description}</DialogDescription>
            </Box>
            <DialogClose
              className="rdr-sheet__close"
              aria-label={__('Close and return to the page')}
            >
              <X size={22} />
            </DialogClose>
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-sheet__body" key={panel ?? 'empty'}>
            {panel === 'enquiry' && (
              <Box as={BOX_TAG.DIV} className="rdr-sheet-enquiry">
                <Box as={BOX_TAG.DIV} className="rdr-sheet-enquiry__intro">
                  <GraduationCap size={26} />
                  <Box as={BOX_TAG.STRONG}>{__("Let's talk about your child's future.")}</Box>
                  <Box as={BOX_TAG.P}>
                    {__('Share your details. Our school office is only a message away.')}
                  </Box>
                  <ExternalAnchor href="tel:+919050077550">
                    <Phone size={15} />
                    {__('Prefer to call? ')}
                    {school.phone}
                  </ExternalAnchor>
                </Box>
                {enquiry}
              </Box>
            )}
            {story && (
              <StoryPanel initial={selectedStory} navigate={(p) => onNavigate(p as SchoolPanel)} />
            )}
            {panel === 'academics' && <AcademicPanel onEnquire={onEnquire} />}
            {panel === 'campus' && <CampusPanel onEnquire={onEnquire} />}
            {panel === 'gallery' && <GalleryPanel onEnquire={onEnquire} />}
            {panel === 'admissions' && <AdmissionsPanel onEnquire={onEnquire} />}
            {panel === 'resources' && <ResourcesPanel />}
            {panel === 'contact' && <ContactPanel onEnquire={onEnquire} />}
          </Box>
          {panel !== 'enquiry' && (
            <Box as={BOX_TAG.DIV} className="rdr-sheet__bottom">
              <Box as={BOX_TAG.SPAN}>
                <ShieldCheck size={15} />
                {__('RDR School · Haluhera, Rewari')}
              </Box>
              <Button unstyled onClick={onEnquire} type="button">
                {__('Quick enquiry ')}
                <ArrowRight size={16} />
              </Button>
            </Box>
          )}
        </DialogContent>
      </DialogPortal>
    </DialogRoot>
  );
};
