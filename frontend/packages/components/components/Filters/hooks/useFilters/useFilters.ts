import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { _find } from '@Base/utilities/find';

import { FilterOperator } from '../../constants/filter.operator';
import { isFilterValid } from '../../helpers/filters.isValid';
import type { FilterCondition, FiltersProps } from '../../types/filters.props';

export const useFilters = ({ value, onChange, fields }: FiltersProps) => {
  const [open, setOpen] = useState(false);
  const [draft, setDraft] = useState<readonly FilterCondition[]>(value);
  const [editing, setEditing] = useState<FilterCondition | null>(null);

  const onOpenChange = (next: boolean) => {
    if (next) {
      setDraft(value);
      setEditing(null);
    }
    setOpen(next);
  };
  const addInline = () => setDraft((items) => [...items, { id: crypto.randomUUID(), field: '', operator: FilterOperator.EQUALS, values: [] }]);
  const updateInline = (next: FilterCondition) => setDraft((items) => _map(items, (item) => {
    if (item.id === next.id) return next;
    return item;
  }));
  const add = () => setEditing({ id: crypto.randomUUID(), field: '', operator: FilterOperator.EQUALS, values: [] });
  const edit = (condition: FilterCondition) => setEditing(condition);
  const remove = (id: string) => setDraft((items) => _filter(items, (item) => item.id !== id));
  const selectedField = _find(fields, (field) => field.id === editing?.field);
  const selectField = (id: string | null) => {
    if (!editing || !id) return;
    const field = _find(fields, (entry) => entry.id === id);
    if (!field) return;
    setEditing({ ...editing, field: id, operator: field.operators[0] ?? FilterOperator.EQUALS, values: [] });
  };
  const save = () => {
    if (!editing || !isFilterValid(editing, selectedField)) return;
    if (_find(draft, (item) => item.id === editing.id)) {
      setDraft(_map(draft, (item) => {
        if (item.id === editing.id) return editing;
        return item;
      }));
    } else {
      setDraft([...draft, editing]);
    }
    setEditing(null);
  };
  const apply = () => {
    if (draft.some((condition) => !isFilterValid(condition, _find(fields, (field) => field.id === condition.field)))) return;
    onChange(draft);
    setOpen(false);
  };
  return {
    open, draft, editing, selectedField, onOpenChange, add, addInline, updateInline, edit, remove, selectField, save, apply,
    reset: () => setDraft([]),
    back: () => setEditing(null),
    setEditing,
  };
};
