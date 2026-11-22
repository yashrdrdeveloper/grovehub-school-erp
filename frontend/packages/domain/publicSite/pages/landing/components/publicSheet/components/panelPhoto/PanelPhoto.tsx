import { Sparkles } from 'lucide-react';
import { photographs } from '../../../../config/content';

const PanelPhoto = ({index,alt}:{index:number;alt:string}) => {
  const photo=photographs[index];
  if (!photo) return null;
  return <div className="rdr-sheet-photo">
    <img loading="lazy" src={photo.src} alt={alt} style={{objectPosition:photo.position}}
      onError={event=>{
        if(event.currentTarget.dataset.retry) {event.currentTarget.hidden=true;return;}
        event.currentTarget.dataset.retry='1';
        event.currentTarget.src=photo.originalSrc;
      }}/>
    <span aria-hidden="true"><Sparkles size={17}/> RDR SCHOOL</span>
  </div>;
};

export { PanelPhoto };
