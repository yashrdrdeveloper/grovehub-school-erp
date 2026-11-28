import { Plus } from '@Components/icons';
import { Box } from '../../../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../../../Button';
import { ButtonWithIcon } from '../../../../../ButtonWithIcon';
import { Tooltip } from '../../../../../Tooltip';

import styles from '../../DesktopFilters.module.css';
import type { DesktopFilterActionsProps } from './types/desktopFilterActions.props';

export const DesktopFilterActions = ({ controls }: DesktopFilterActionsProps) => (
  <Box className={styles.footer}>
    <Tooltip label={controls.addReason} focusable={Boolean(controls.addReason)}>
      <ButtonWithIcon
        icon={Plus}
        iconSize={14}
        size={BUTTON_SIZE.SMALL}
        variant={BUTTON_VARIANT.GHOST}
        disabled={!controls.canAdd}
        onClick={controls.addInline}
      >
        {__('Add filter')}
      </ButtonWithIcon>
    </Tooltip>
    <Box className={styles.actions}>
      <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>
        {__('Reset')}
      </Button>
      <Tooltip
        label={
          (!controls.isDraftValid &&
            __('Complete all filter fields and values before applying.')) ||
          undefined
        }
        focusable={Boolean(
          (!controls.isDraftValid &&
            __('Complete all filter fields and values before applying.')) ||
            undefined,
        )}
      >
        <Button
          size={BUTTON_SIZE.SMALL}
          variant={BUTTON_VARIANT.PRIMARY}
          disabled={!controls.isDraftValid}
          onClick={controls.apply}
        >
          {__('Apply filters')}
        </Button>
      </Tooltip>
    </Box>
  </Box>
);
