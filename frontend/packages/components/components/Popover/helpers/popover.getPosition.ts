import type { PopoverPositionInput } from './types/popover.getPosition';

export const getPopoverPosition = ({ anchor, panelHeight, viewportWidth, viewportHeight, width }: PopoverPositionInput) => {
  const margin = 12;
  const gap = 4;
  const below = Math.max(0, viewportHeight - anchor.bottom - margin - gap);
  const above = Math.max(0, anchor.top - margin - gap);
  const upward = below < panelHeight && above > below;
  const maxHeight = Math.max(0, Math.max(above, below));
  const height = Math.min(panelHeight, maxHeight);
  const left = Math.max(margin, Math.min(anchor.left, viewportWidth - width - margin));
  const top = Math.max(margin, anchor.bottom + gap - Number(upward) * (anchor.height + gap * 2 + height));
  return { left, top, width, maxHeight };
};
