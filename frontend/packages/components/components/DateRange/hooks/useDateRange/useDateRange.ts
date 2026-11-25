import { useState } from 'react';

import { formatRangeDate, getRangeMonth, getRangeSelection, shiftRangeMonth } from '../../helpers/dateRange.calendar';
import type { DateRangePreset } from '../../types/dateRange.preset';
import type { DateRangeProps } from '../../types/dateRange.props';

export const useDateRange = ({ value, onChange, disabled }: DateRangeProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState(value);
  const [month, setMonth] = useState(() => getRangeMonth(value.startDate, new Date()));
  const [hover, setHover] = useState<string | null>(null);
  const [selectedPreset, setSelectedPreset] = useState<string | null>(null);
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
      setSelectedPreset(null);
      setHover(null);
    }
    setOpen(next);
  };
  const selectDate = (date: string | null) => {
    if (!date || disabled) return;
    setDraft((current) => getRangeSelection(current, date));
    setSelectedPreset(null);
    setHover(null);
  };
  const selectPreset = (preset: DateRangePreset) => {
    if (disabled) return;
    const today = new Date();
    const start = preset.start(today);
    const end = preset.end(today);
    setDraft({ startDate: formatRangeDate(start), endDate: formatRangeDate(end) });
    setMonth(shiftRangeMonth(start, 0));
    setSelectedPreset(preset.id);
    setHover(null);
  };
  const clear = () => {
    if (disabled) return;
    setDraft({ startDate: null, endDate: null });
    setSelectedPreset(null);
    setHover(null);
  };
  const apply = () => {
    if (!canApply) return;
    onChange(draft);
    setOpen(false);
  };
  return {
    open, draft, month, nextMonth: shiftRangeMonth(month, 1), selectedPreset, preview: getPreview(),
    invalid, incomplete, canApply, onOpenChange, setMonth, setHover, selectDate, selectPreset, clear, apply,
  };
};
