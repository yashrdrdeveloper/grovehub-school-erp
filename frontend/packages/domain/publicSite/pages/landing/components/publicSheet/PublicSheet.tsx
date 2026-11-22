import type { ReactNode } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { ArrowRight, GraduationCap, Phone, ShieldCheck, X } from '@Components/icons';
import { school } from '../../constants/school.content';
import { panelTitles, overviewPanels, type StoryKey } from './constants/school.sheetCopy';
import { StoryPanel } from './components/storyPanel';
import { AcademicPanel, CampusPanel, GalleryPanel, AdmissionsPanel, ResourcesPanel, ContactPanel } from './components/infoPanels';
import type { SchoolPanel } from './types/schoolPanel';

export type { SchoolPanel } from './types/schoolPanel';

export const PublicSheet = ({
  panel, onClose, onNavigate, enquiry,
}: {
  panel:SchoolPanel|null;
  onClose:()=>void;
  onNavigate:(panel:SchoolPanel)=>void;
  enquiry:ReactNode;
}) => {
  let selectedStory: StoryKey = 'history';
  if (panel === 'director' || panel === 'principal' || panel === 'philosophy' || panel === 'history') selectedStory = panel;
  const story=panel==='about'||panel==='history'||panel==='director'||panel==='principal'||panel==='philosophy';
  const onEnquire=()=>onNavigate('enquiry');
  let title = 'RDR School information';
  let description = '';
  if (panel) { title = panelTitles[panel]; description = overviewPanels[panel]; }
  return <Dialog.Root open={Boolean(panel)} onOpenChange={open=>{if(!open)onClose();}}>
    <Dialog.Portal>
      <Dialog.Overlay className="rdr-sheet-backdrop"/>
      <Dialog.Content className="rdr-sheet" aria-label={title}>
        <div className="rdr-sheet__handle" aria-hidden="true"/>
        <div className="rdr-sheet__header">
          <div><span className="rdr-sheet-kicker">RAO DEENA RAM VIDYA VIHAR</span><Dialog.Title>{title}</Dialog.Title>
            <Dialog.Description>{description}</Dialog.Description>
          </div>
          <Dialog.Close className="rdr-sheet__close" aria-label="Close and return to the page"><X size={22}/></Dialog.Close>
        </div>
        <div className="rdr-sheet__body" key={panel??'empty'}>
          {panel==='enquiry'&&<div className="rdr-sheet-enquiry"><div className="rdr-sheet-enquiry__intro"><GraduationCap size={26}/><strong>Let's talk about your child's future.</strong><p>Share your details. Our school office is only a message away.</p>
            <a href="tel:+919050077550"><Phone size={15}/> Prefer to call? {school.phone}</a></div>{enquiry}</div>}
          {story&&<StoryPanel initial={selectedStory} navigate={p=>onNavigate(p as SchoolPanel)}/>}
          {panel==='academics'&&<AcademicPanel onEnquire={onEnquire}/>}
          {panel==='campus'&&<CampusPanel onEnquire={onEnquire}/>}
          {panel==='gallery'&&<GalleryPanel onEnquire={onEnquire}/>}
          {panel==='admissions'&&<AdmissionsPanel onEnquire={onEnquire}/>}
          {panel==='resources'&&<ResourcesPanel/>}
          {panel==='contact'&&<ContactPanel onEnquire={onEnquire}/>}
        </div>
        {panel!=='enquiry'&&<div className="rdr-sheet__bottom"><span><ShieldCheck size={15}/> RDR School · Haluhera, Rewari</span><button onClick={onEnquire} type="button">Quick enquiry <ArrowRight size={16}/></button></div>}
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
};
