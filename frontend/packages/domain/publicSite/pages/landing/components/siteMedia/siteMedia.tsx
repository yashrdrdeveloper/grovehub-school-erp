import type { ReactNode } from 'react';
import { BookOpen, Bus, Computer, FlaskConical, HeartPulse, Pencil, Presentation, ShieldCheck, Sparkles, Star, Trophy } from '@Components/icons';
import { photographs } from '../../constants/school.content';

const facilityIcons = [Computer, FlaskConical, Presentation, BookOpen, Trophy, Bus, HeartPulse, ShieldCheck];
const ImageCard = ({
  src,alt,position,className='',loading='lazy',children,
}: {
  src:string;alt:string;position?:string;className?:string;loading?:'lazy'|'eager';children?:ReactNode;
}) => (
  <div className={`rdr-photo ${className}`}>
    <img src={photographs.find(photo => photo.src === src)?.originalSrc ?? src}
      alt={alt} loading={loading} decoding="async"
      style={{objectPosition:position ?? 'center'}} referrerPolicy="no-referrer"
      onError={event => {
        const image = event.currentTarget;
        const localPhoto = photographs.find(photo => photo.originalSrc === image.src);
        if (localPhoto && !image.dataset.fallbackAttempted) {
          image.dataset.fallbackAttempted = 'true';
          image.src = localPhoto.src;
          return;
        }
        image.style.display = 'none';
        image.parentElement?.classList.add('rdr-photo--unavailable');
      }} />
    {children}
  </div>
);
const Heading = ({eyebrow,title,description}:{eyebrow:string;title:string;description?:string}) => (
  <div className="rdr-section-heading"><span className="rdr-eyebrow">{eyebrow}</span><h2>{title}</h2>{description&&<p>{description}</p>}</div>
);

/** Purely decorative school motifs. Keep them outside important photo faces
 * and out of the accessibility tree; the original photographs remain intact. */
const PhotoMotifs = ({variant='hero'}:{variant?:'hero'|'gallery'|'banner'}) => (
  <div className={`rdr-photo-motifs rdr-photo-motifs--${variant}`} aria-hidden="true">
    <svg className="rdr-photo-motifs__flight" viewBox="0 0 170 76" fill="none">
      <path d="M4 59C34 26 48 83 83 44c20-21 25-38 51-36" stroke="currentColor" strokeWidth="1.6" strokeDasharray="4 7" strokeLinecap="round"/>
      <path d="m131 8 32-5-18 24-5-12-9-7Zm9 7 23-12" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"/>
    </svg>
    <span className="rdr-photo-motifs__spark"><Sparkles size={27} strokeWidth={1.4}/></span>
    <span className="rdr-photo-motifs__star"><Star size={18} strokeWidth={1.5}/></span>
    <span className="rdr-photo-motifs__book"><BookOpen size={25} strokeWidth={1.4}/></span>
    <span className="rdr-photo-motifs__pencil"><Pencil size={23} strokeWidth={1.5}/></span>
  </div>
);

export { facilityIcons };
export { ImageCard };
export { Heading };
export { PhotoMotifs };
