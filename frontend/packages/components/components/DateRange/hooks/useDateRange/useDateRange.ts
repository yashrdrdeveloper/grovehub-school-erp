import { useState } from 'react';

import { formatRangeDate, getRangeMonth, getRangeSelection, shiftRangeMonth } from '../../helpers/dateRange.calendar';
import { CUSTOM_DATE_RANGE_PRESET, DATE_RANGE_PRESETS } from '../../constants/dateRange.presets';
import type { DateRangePreset } from '../../types/dateRange.preset';
import type { DateRangeProps } from '../../types/dateRange.props';

export const useDateRange = ({ value, onChange, disabled, showFooter = true }: DateRangeProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState(value);
  const [month, setMonth] = useState(() => getRangeMonth(value.startDate, new Date()));
  const [hover, setHover] = useState<string | null>(null);
  const getPreset = () => {
    if (!draft.startDate && !draft.endDate) return draft.presetId ?? null;
    if (DATE_RANGE_PRESETS.some((preset) => preset.id === draft.presetId)) return draft.presetId;
    return CUSTOM_DATE_RANGE_PRESET;
  };
  const selectedPreset = getPreset();
  const invalid = Boolean(draft.startDate && draft.endDate && draft.startDate > draft.endDate);
  const incomplete = Boolean(draft.startDate) !== Boolean(draft.endDate);
  const canApply = !disabled && !invalid && !incomplete;
  const getPreview = () => {
    if (draft.startDate && !draft.endDate && hover) return getRangeSelection(draft, hover);
    return draft;
  };
  const onOpenChange = (next: boolean) => {
    if (next && disabled) return;
    if (next) {
      setDraft(value);
      setMonth(getRangeMonth(value.startDate, new Date()));
      setHover(null);
    }
    setOpen(next);
  };
  const selectDate = (date: string | null) => {
    if (!date || disabled || !open) return;
    const next = { ...getRangeSelection(draft, date), presetId: CUSTOM_DATE_RANGE_PRESET };
    setDraft(next);
    setHover(null);
    if (!showFooter && next.startDate && next.endDate) {
      onChange(next);
      setOpen(false);
    }
  };
  const selectPreset = (preset: DateRangePreset) => {
    if (disabled) return;
    const today = new Date();
    const start = preset.start(today);
    const end = preset.end(today);
    const next = { startDate: formatRangeDate(start), endDate: formatRangeDate(end), presetId: preset.id };
    setDraft(next);
    setMonth(shiftRangeMonth(start, 0));
    setHover(null);
    if (!showFooter) {
      onChange(next);
      setOpen(false);
    }
  };
  const selectCustom = () => {
    if (disabled) return;
    setDraft((current) => ({ ...current, presetId: CUSTOM_DATE_RANGE_PRESET }));
    setHover(null);
  };
  const clear = () => {
    if (disabled) return;
    setDraft({ startDate: null, endDate: null });
    setHover(null);
  };
  const apply = () => {
    if (!canApply) return;
    onChange(draft);
    setOpen(false);
  };
  const onOutsideClick = () => {
    if (!open) return;
    const changed = draft.startDate !== value.startDate || draft.endDate !== value.endDate || draft.presetId !== value.presetId;
    if (canApply && changed) onChange(draft);
    setOpen(false);
  };
  return {
    open, draft, month, nextMonth: shiftRangeMonth(month, 1), selectedPreset, preview: getPreview(),
    invalid, incomplete, canApply, onOpenChange, onOutsideClick, setMonth, setHover, selectDate, selectPreset, selectCustom, clear, apply,
  };
};
