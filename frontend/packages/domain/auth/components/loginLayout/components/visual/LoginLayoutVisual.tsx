import { Box, BOX_TAG } from '@Components/components/Box';
import { SchoolMonogram } from '@Components';

import { BookOpen, CalendarDays, GraduationCap, MessageCircle } from '@Components/icons';

import styles from './LoginLayoutVisual.module.css';
import { LOGIN_LAYOUT_LEGACY_CLASSES } from '../../constants/loginLayout.legacyClasses';
import type { LoginLayoutVisualComponent } from './types/loginLayoutVisual';

/** Shared school identity artwork. No session state, effects or API calls. */
export const LoginLayoutVisual: LoginLayoutVisualComponent = () => (
  <Box
    as={BOX_TAG.SECTION}
    className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL}
    aria-label={__('RDR School branded gradient background')}
  >
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.MOBILE_HEADER}>
      <Box
        as={BOX_TAG.DIV}
        className={LOGIN_LAYOUT_LEGACY_CLASSES.MOBILE_IDENTITY}
        aria-label={__('RDR School')}
      >
        <SchoolMonogram />
        <Box as={BOX_TAG.SPAN} className={LOGIN_LAYOUT_LEGACY_CLASSES.MOBILE_SCHOOL_NAME}>
          {__('RDR School')}
          <Box as={BOX_TAG.SMALL}>{__('Your school workspace')}</Box>
        </Box>
      </Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_BRAND}>
      <SchoolMonogram />
      <Box as={BOX_TAG.SPAN}>
        <Box as={BOX_TAG.STRONG}>{__('RDR School')}</Box>
        <Box as={BOX_TAG.SMALL}>{__('Rao Deena Ram Vidya Vihar')}</Box>
      </Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_COPY}>
      <Box as={BOX_TAG.SPAN} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_KICKER}>
        <Box as={BOX_TAG.SPAN} />
        {__('THE RDR SCHOOL WORKSPACE')}
      </Box>
      <Box as={BOX_TAG.H2}>
        {__('For every school day.')}
        <br />
        <Box as={BOX_TAG.SPAN}>{__('And everyone behind it.')}</Box>
      </Box>
      <Box as={BOX_TAG.P}>{__('A connected space for the people who make learning happen.')}</Box>
      <Box as={BOX_TAG.DIV} className={styles.features} aria-hidden="true">
        <Box as={BOX_TAG.DIV} className={styles.feature}>
          <GraduationCap size={19} />
          <Box as={BOX_TAG.SPAN}>{__('Student management')}</Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles.feature}>
          <CalendarDays size={19} />
          <Box as={BOX_TAG.SPAN}>{__('Attendance & calendar')}</Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles.feature}>
          <BookOpen size={19} />
          <Box as={BOX_TAG.SPAN}>{__('Academic operations')}</Box>
        </Box>
        <Box as={BOX_TAG.DIV} className={styles.feature}>
          <MessageCircle size={19} />
          <Box as={BOX_TAG.SPAN}>{__('School communication')}</Box>
        </Box>
      </Box>
    </Box>
    <Box as={BOX_TAG.DIV} className={LOGIN_LAYOUT_LEGACY_CLASSES.VISUAL_FOOTER}>
        <Box as={BOX_TAG.SPAN}>{__('LEARN')}</Box>
        <Box as={BOX_TAG.I} aria-hidden="true" />
        <Box as={BOX_TAG.SPAN}>{__('ORGANISE')}</Box>
        <Box as={BOX_TAG.I} aria-hidden="true" />
        <Box as={BOX_TAG.SPAN}>{__('GROW')}</Box>
    </Box>
  </Box>
);
