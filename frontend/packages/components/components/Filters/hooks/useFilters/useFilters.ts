import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { _find } from '@Base/utilities/find';

import { FilterOperator } from '../../constants/filter.operator';
import { isFilterValid } from '../../helpers/filters.isValid';
import type { FilterCondition, FiltersProps } from '../../types/filters.props';

export const useFilters = ({ value, onChange, fields, maxFilters }: FiltersProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState<readonly FilterCondition[]>(value);
  const [editing, setEditing] = useState<FilterCondition | null>(null);

  const getAvailableFields = (id?: string) =>
    _filter(fields, (field) => !draft.some((item) => item.id !== id && item.field === field.id));
  const canAdd = draft.length < Math.min(fields.length, maxFilters ?? fields.length);
  const isDraftValid =
    new Set(draft.map((item) => item.field)).size === draft.length &&
    draft.every((condition) =>
      isFilterValid(condition, _find(fields, (field) => field.id === condition.field)),
    );

  const onOpenChange = (next: boolean) => {
    if (next) {
      setDraft(value);
      setEditing(null);
    }
    setOpen(next);
  };

  const addInline = () => {
    if (!canAdd) return;
    const condition = { id: crypto.randomUUID(), field: '', operator: FilterOperator.EQUALS, values: [] };
    setDraft((items) => {
      if (items.length >= Math.min(fields.length, maxFilters ?? fields.length)) return items;
      return [...items, condition];
    });
  };

  const updateInline = (next: FilterCondition) => {
    setDraft((items) => {
      if (next.field && items.some((item) => item.id !== next.id && item.field === next.field)) {
        return items;
      }
      return _map(items, (item) => {
        if (item.id === next.id) return next;
        return item;
      });
    });
  };

  const add = () => {
    if (!canAdd) return;
    setEditing({ id: crypto.randomUUID(), field: '', operator: FilterOperator.EQUALS, values: [] });
  };
  const edit = (condition: FilterCondition) => setEditing(condition);
  const remove = (id: string) => setDraft((items) => _filter(items, (item) => item.id !== id));
  const selectedField = _find(getAvailableFields(editing?.id), (field) => field.id === editing?.field);
  const selectField = (id: string | null) => {
    if (!editing || !id) return;
    const field = _find(getAvailableFields(editing.id), (entry) => entry.id === id);
    if (!field) return;
    setEditing({ ...editing, field: id, operator: field.operators[0] ?? FilterOperator.EQUALS, values: [] });
  };

  const save = () => {
    if (!editing || !isFilterValid(editing, selectedField)) return;
    const existing = _find(draft, (item) => item.id === editing.id);
    if (!existing && !canAdd) return;
    setDraft((items) => {
      if (items.some((item) => item.id !== editing.id && item.field === editing.field)) return items;
      if (existing) {
        return _map(items, (item) => {
          if (item.id === editing.id) return editing;
          return item;
        });
      }
      if (items.length >= Math.min(fields.length, maxFilters ?? fields.length)) return items;
      return [...items, editing];
    });
    setEditing(null);
  };
  const apply = () => {
    if (!isDraftValid) return;
    onChange(draft);
    setOpen(false);
  };

  return {
    open, draft, editing, selectedField, getAvailableFields, canAdd, isDraftValid,
    onOpenChange, add, addInline, updateInline, edit, remove, selectField, save, apply,
    reset: () => setDraft([]),
    back: () => setEditing(null),
    setEditing,
  };
};
