import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';

import type { SelectValue } from '../../types/select.props';
import type { UseSelectionProps } from './types/useSelection';

export const useSelection = <T extends SelectValue>({
  value,
  multiple,
  onChange,
  maxSelections,
}: UseSelectionProps<T>) => {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');
  const [draft, setDraft] = useState<readonly T[]>(value);

  const onOpenChange = (next: boolean) => {
    if (next) {
      setDraft(value);
      setSearch('');
    }
    setOpen(next);
  };

  const onToggle = (item: T) => {
    if (!multiple) {
      onChange([item]);
      setOpen(false);
      return;
    }
    if (draft.includes(item)) {
      setDraft(_filter(draft, (entry) => entry !== item));
      return;
    }
    if (maxSelections !== undefined && draft.length >= maxSelections) return;
    setDraft([...draft, item]);
  };

  const onClear = () => {
    if (multiple) {
      setDraft([]);
      return;
    }
    onChange([]);
    setOpen(false);
  };

  const onApply = () => {
    onChange(draft);
    setOpen(false);
  };

  return {
    open,
    search,
    selected: draft,
    onOpenChange,
    onSearch: setSearch,
    onToggle,
    onClear,
    onApply,
  };
};
