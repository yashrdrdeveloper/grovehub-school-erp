import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';

import { SortDirection } from '../../constants/sort.direction';
import type { SortProps, SortCondition } from '../../types/sort.props';

export const useSort = ({ value, onChange, fields }: SortProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState<readonly SortCondition[]>(value);
  const available = _filter(fields, (field) => !draft.some((item) => item.field === field.id));
  const onOpenChange = (next: boolean) => {
    if (next) setDraft(value);
    setOpen(next);
  };
  const add = () => {
    const field = available[0];
    if (field) setDraft([...draft, { field: field.id, direction: SortDirection.ASC }]);
  };
  const update = (index: number, change: Partial<SortCondition>) => {
    setDraft(_map(draft, (item, position) => {
      if (index === position) return { ...item, ...change };
      return item;
    }));
  };
  const move = (index: number, offset: number) => {
    const target = index + offset;
    const current = draft[index];
    const other = draft[target];
    if (!current || !other) return;
    setDraft(_map(draft, (item, position) => {
      if (position === index) return other;
      if (position === target) return current;
      return item;
    }));
  };
  return {
    open, draft, available, onOpenChange, add, update, move,
    remove: (index: number) => setDraft(_filter(draft, (item) => item !== draft[index])),
    reset: () => setDraft([]),
    apply: () => {
      onChange(draft);
      setOpen(false);
    },
  };
};
