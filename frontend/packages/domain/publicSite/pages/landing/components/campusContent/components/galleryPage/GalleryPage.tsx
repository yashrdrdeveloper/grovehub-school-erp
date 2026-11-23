import { _map } from '@Base/utilities/map';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ExternalLink } from '@Components/icons';
import { photographs } from '../../../../constants/school.content';
import { ImageCard, PhotoMotifs, PHOTO_MOTIFS_VARIANT } from '../../../siteMedia';
import { IntroPage } from '../../../aboutContent';
import { getGalleryPhotoClassName } from '../.././helpers/campusContent.getGalleryPhotoClassName';

export const GalleryPage = () => (
  <>
    <IntroPage
      eyebrow="GALLERY"
      title={__('Life at RDR, in moments.')}
      description={__(
        'A photo story of school celebrations, student life and achievements — captured at RDR School.',
      )}
    />
    <Box as={BOX_TAG.SECTION} className="rdr-section rdr-container">
      <Box as={BOX_TAG.DIV} className="rdr-gallery-grid">
        {_map(photographs, (photo, i) => (
          <figure className={getGalleryPhotoClassName(i)} key={photo.src}>
            <Box as={BOX_TAG.DIV} className="rdr-gallery-grid__art">
              <ImageCard src={photo.src} alt={photo.alt} position={photo.position} />
              {i % 3 === 0 && <PhotoMotifs variant={PHOTO_MOTIFS_VARIANT.GALLERY} />}
            </Box>
            <figcaption>
              <Box as={BOX_TAG.SPAN}>{photo.category}</Box>
              <Box as={BOX_TAG.STRONG}>{photo.title}</Box>
            </figcaption>
          </figure>
        ))}
      </Box>
      <Box as={BOX_TAG.P} className="rdr-gallery__archive-note">
        {__('From the RDR School photo archives. ')}
        <ExternalAnchor href="https://rdrschool.com/gallery.html" target="_blank" rel="noreferrer">
          {__('See the original gallery ')}
          <ExternalLink size={13} />
        </ExternalAnchor>
      </Box>
    </Box>
  </>
);
