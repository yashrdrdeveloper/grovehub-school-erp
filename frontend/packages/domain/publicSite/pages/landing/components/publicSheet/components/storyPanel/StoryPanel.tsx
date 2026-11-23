import { _map } from '@Base/utilities/map';
import { ExternalAnchor } from '@Components/components/ExternalAnchor';
import { Button } from '@Components/components/Button';
import { Box, BOX_TAG } from '@Components/components/Box';
import { useState } from 'react';
import { ArrowRight, Check, ExternalLink } from '@Components/icons';
import { legacyStories, storyNavigation, type StoryKey } from '../../constants/school.sheetCopy';
import { PanelPhoto } from '../panelPhoto';
import { getStoryPanelPhotoIndex } from './helpers/storyPanel.getPhotoIndex';
import { getStoryPanelTabClassName } from './helpers/storyPanel.getTabClassName';

const StoryPanel = ({
  initial = 'history',
  navigate,
}: {
  initial?: StoryKey;
  navigate: (path: string) => void;
}) => {
  const [selected, setSelected] = useState<StoryKey>(initial);
  const story = legacyStories[selected];
  return (
    <Box as={BOX_TAG.DIV} className="rdr-sheet-story">
      <Box
        as={BOX_TAG.DIV}
        className="rdr-sheet-story__tabs"
        role="tablist"
        aria-label={__('School story topics')}
      >
        {_map(storyNavigation, (item) => {
          const Icon = item.icon;
          return (
            <Button
              unstyled
              type="button"
              role="tab"
              aria-selected={selected === item.id}
              key={item.id}
              className={getStoryPanelTabClassName(selected === item.id)}
              onClick={() => setSelected(item.id)}
            >
              <Icon size={16} />
              {item.label}
            </Button>
          );
        })}
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-sheet-story__layout" role="tabpanel">
        <Box as={BOX_TAG.DIV} className="rdr-sheet-story__copy">
          <Box as={BOX_TAG.SPAN} className="rdr-sheet-kicker">
            {story.eyebrow}
          </Box>
          <Box as={BOX_TAG.H3}>{story.title}</Box>
          <Box as={BOX_TAG.P}>{story.description}</Box>
          <Box as={BOX_TAG.DIV} className="rdr-sheet-story__points">
            {_map(story.highlights, (label) => (
              <Box as={BOX_TAG.SPAN} key={label}>
                <Check size={15} />
                {label}
              </Box>
            ))}
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-sheet-story__footer">
            <ExternalAnchor href={story.link} target="_blank" rel="noreferrer">
              {story.linkLabel}
              <ExternalLink size={15} />
            </ExternalAnchor>
            <Button unstyled type="button" onClick={() => navigate('enquiry')}>
              {__('Ask us a question ')}
              <ArrowRight size={16} />
            </Button>
          </Box>
          {(selected === 'principal' || selected === 'director') && (
            <Box as={BOX_TAG.SMALL}>
              {__(
                "Role and message are from the school's historical website. Please confirm current leadership before official publication.",
              )}
            </Box>
          )}
        </Box>
        <PanelPhoto
          index={getStoryPanelPhotoIndex(selected)}
          alt={__('Students participating in activities at RDR School')}
        />
      </Box>
    </Box>
  );
};

export { StoryPanel };
