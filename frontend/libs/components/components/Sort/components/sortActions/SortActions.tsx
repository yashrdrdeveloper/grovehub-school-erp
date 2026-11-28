import { Plus } from '@Components/icons';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { Tooltip } from '../../../Tooltip';

import styles from '../../Sort.module.css';
import type { SortActionsProps } from './types/sortActions.props';

export const SortActions = ({ controls }: SortActionsProps) => {
  const invalidReason =
    (!controls.isDraftValid && __('Choose a unique field for each sort rule.')) || undefined;
  return (
    <Box className={styles.footer}>
      <Tooltip label={controls.addReason} focusable={Boolean(controls.addReason)}>
        <ButtonWithIcon
          icon={Plus}
          iconSize={14}
          size={BUTTON_SIZE.SMALL}
          variant={BUTTON_VARIANT.GHOST}
          className={styles.add}
          disabled={!controls.canAdd}
          onClick={controls.add}
        >
          {__('Add sort')}
        </ButtonWithIcon>
      </Tooltip>
      <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} onClick={controls.reset}>
        {__('Reset')}
      </Button>
      <Tooltip label={invalidReason} focusable={Boolean(invalidReason)}>
        <Button size={BUTTON_SIZE.SMALL} disabled={!controls.isDraftValid} onClick={controls.apply}>
          {__('Apply sort')}
        </Button>
      </Tooltip>
    </Box>
  );
};
