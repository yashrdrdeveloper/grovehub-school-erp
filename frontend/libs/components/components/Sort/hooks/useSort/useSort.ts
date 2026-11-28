import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';

import { SortDirection } from '../../constants/sort.direction';
import type { SortProps, SortCondition } from '../../types/sort.props';

export const useSort = ({
  value,
  onChange,
  fields,
  maxFields,
  initialOpen,
  onDismiss,
}: SortProps) => {
  const [open, setOpen] = useState(Boolean(initialOpen));
  const [draft, setDraft] = useState<readonly SortCondition[]>(value);
  const getAvailableFields = (index?: number) =>
    _filter(
      fields,
      (field) => !draft.some((item, position) => position !== index && item.field === field.id),
    );
  const available = getAvailableFields();
  const canAdd = Boolean(available.length) && draft.length < (maxFields ?? fields.length);
  const getAddReason = () => {
    if (canAdd) return undefined;
    if (!available.length) return __('All sort fields are selected.');
    if (maxFields === 1) return __('Only one sort field is allowed.');
    return __('Maximum of {{count}} sort fields reached.', { count: maxFields ?? fields.length });
  };
  const isDraftValid =
    new Set(draft.map((item) => item.field)).size === draft.length &&
    draft.every((item) => fields.some((field) => field.id === item.field));

  const onOpenChange = (next: boolean) => {
    if (next) setDraft(value);
    setOpen(next);
    if (!next) onDismiss?.();
  };
  const add = () => {
    setDraft((items) => {
      if (items.length >= (maxFields ?? fields.length)) return items;
      const field = fields.find((entry) => !items.some((item) => item.field === entry.id));
      if (!field) return items;
      return [...items, { field: field.id, direction: SortDirection.ASC }];
    });
  };
  const update = (index: number, change: Partial<SortCondition>) => {
    setDraft((items) => {
      if (
        change.field !== undefined &&
        (!fields.some((field) => field.id === change.field) ||
          items.some((item, position) => position !== index && item.field === change.field))
      ) {
        return items;
      }
      return _map(items, (item, position) => {
        if (index === position) return { ...item, ...change };
        return item;
      });
    });
  };
  const reorder = (items: readonly SortCondition[]) => setDraft(items);
  return {
    open,
    draft,
    available,
    getAvailableFields,
    canAdd,
    isDraftValid,
    addReason: getAddReason(),
    onOpenChange,
    add,
    update,
    reorder,
    remove: (index: number) => setDraft((items) => _filter(items, (item) => item !== items[index])),
    reset: () => setDraft([]),
    apply: () => {
      if (!isDraftValid) return;
      onChange(draft);
      onOpenChange(false);
    },
  };
};
