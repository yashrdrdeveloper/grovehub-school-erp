import { useState } from 'react';
import { ChevronDown } from '@Components/icons';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT } from '../Button';
import { ButtonWithIcon, BUTTON_ICON_POSITION } from '../ButtonWithIcon';
import { Popover } from '../Popover';
import { Number } from '../Number';

import styles from './NumberRange.module.css';
import type { NumberRangeProps } from './types/numberRange.props';

export const NumberRange = ({ value, onChange, disabled, dataTestId }: NumberRangeProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState(value);
  const label = __('Number range');
  const invalid = draft.min !== null && draft.max !== null && draft.min > draft.max;
  const summary = () => {
    if (value.min === null && value.max === null) return label;
    return [value.min ?? '—', value.max ?? '—'].join(' – ');
  };
  const onOpenChange = (next: boolean) => {
    if (next) setDraft(value);
    setOpen(next);
  };
  return (
    <Box dataTestId={dataTestId}>
      <Popover
        open={open}
        onOpenChange={onOpenChange}
        title={label}
        fullWidth
        trigger={
          <ButtonWithIcon
            icon={ChevronDown}
            iconPosition={BUTTON_ICON_POSITION.RIGHT}
            variant={BUTTON_VARIANT.OUTLINE}
            fullWidth
            disabled={disabled}
            onClick={() => onOpenChange(true)}
          >
            {summary()}
          </ButtonWithIcon>
        }
      >
        <Box className={styles.range}>
          <Number
            label={__('Minimum')}
            value={draft.min}
            disabled={disabled}
            onChange={(next) => setDraft({ ...draft, min: next })}
          />
          <Number
            label={__('Maximum')}
            value={draft.max}
            disabled={disabled}
            onChange={(next) => setDraft({ ...draft, max: next })}
          />
          {invalid && <Box role="alert">{__('The start must not exceed the end.')}</Box>}
          <Box className={styles.actions}>
            <Button
              variant={BUTTON_VARIANT.OUTLINE}
              onClick={() => setDraft({ min: null, max: null })}
            >
              {__('Clear')}
            </Button>
            <Button
              disabled={invalid}
              onClick={() => {
                onChange(draft);
                setOpen(false);
              }}
            >
              {__('Apply')}
            </Button>
          </Box>
        </Box>
      </Popover>
    </Box>
  );
};
