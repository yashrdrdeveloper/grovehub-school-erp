import { useCallback } from 'react';
import { _map } from '@Base/utilities/map';
import { Box } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';

import { CUSTOM_DATE_RANGE_PRESET, DATE_RANGE_PRESETS } from '../../constants/dateRange.presets';
import styles from '../../DateRange.module.css';
import type { DateRangePresetsProps } from './types/dateRangePresets.props';

export const DateRangePresets = ({ controls, disabled }: DateRangePresetsProps) => {
  const revealSelectedPreset = useCallback(
    (element: HTMLElement | null) => {
      if (!controls.open) return;
      const selected = element?.querySelector<HTMLElement>('[aria-pressed="true"]');
      if (!element || !selected) return;
      const container = element.getBoundingClientRect();
      const button = selected.getBoundingClientRect();
      if (button.bottom > container.bottom) element.scrollTop += button.bottom - container.bottom;
      if (button.top < container.top) element.scrollTop -= container.top - button.top;
      if (button.right > container.right) element.scrollLeft += button.right - container.right;
      if (button.left < container.left) element.scrollLeft -= container.left - button.left;
    },
    [controls.open],
  );
  return (
    <Box
      ref={revealSelectedPreset}
      className={styles.presets}
      role="group"
      aria-label={__('Date range presets')}
    >
      <Button
        variant={BUTTON_VARIANT.GHOST}
        size={BUTTON_SIZE.SMALL}
        className={styles.preset}
        aria-pressed={controls.selectedPreset === CUSTOM_DATE_RANGE_PRESET}
        disabled={disabled}
        onClick={controls.selectCustom}
      >
        {__('Custom')}
      </Button>
      {_map(DATE_RANGE_PRESETS, (preset) => (
        <Button
          key={preset.id}
          variant={BUTTON_VARIANT.GHOST}
          size={BUTTON_SIZE.SMALL}
          className={styles.preset}
          aria-pressed={controls.selectedPreset === preset.id}
          disabled={disabled}
          onClick={() => controls.selectPreset(preset)}
        >
          {preset.label}
        </Button>
      ))}
    </Box>
  );
};
