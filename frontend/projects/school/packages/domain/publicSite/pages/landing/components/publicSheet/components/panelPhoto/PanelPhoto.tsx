import { Box, BOX_TAG } from '@Components/components/Box';
import { Sparkles } from '@Components/icons';
import styles from './PanelPhoto.module.css';
import { photographs } from '../../../../constants/school.content';

const PanelPhoto = ({ index, alt }: { index: number; alt: string }) => {
  const photo = photographs[index];
  if (!photo) return null;
  return (
    <Box as={BOX_TAG.DIV} className={styles.photo}>
      <img
        loading="lazy"
        src={photo.src}
        alt={alt}
        style={{ objectPosition: photo.position }}
        onError={(event) => {
          if (event.currentTarget.dataset.retry) {
            event.currentTarget.hidden = true;
            return;
          }
          event.currentTarget.dataset.retry = '1';
          event.currentTarget.src = photo.originalSrc;
        }}
      />
      <Box as={BOX_TAG.SPAN} aria-hidden="true">
        <Sparkles size={17} />
        {__('RDR SCHOOL')}
      </Box>
    </Box>
  );
};

export { PanelPhoto };
