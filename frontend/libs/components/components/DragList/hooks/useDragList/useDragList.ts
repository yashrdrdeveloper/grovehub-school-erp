import { useState } from 'react';
import { reorderDragList } from '../../helpers/dragList.reorder';
import type { DragListProps } from '../../types/dragList.props';

export const useDragList = <T>({
  items,
  onReorder,
  canMove,
  getItemId,
  disabled,
}: DragListProps<T>) => {
  const [active, setActive] = useState<string | null>(null);
  const [target, setTarget] = useState<string | null>(null);
  const move = (id: string, to: number) => {
    if (disabled) return;
    const from = items.findIndex((item) => getItemId(item) === id);
    const next = reorderDragList(items, from, to, canMove);
    if (next !== items) onReorder(next);
  };
  const finish = () => {
    setActive(null);
    setTarget(null);
  };
  return { active, target, setActive, setTarget, move, finish };
};
