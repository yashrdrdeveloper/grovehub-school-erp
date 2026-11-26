export interface PopoverPositionInput {
  readonly anchor: Pick<DOMRect, 'left' | 'top' | 'bottom' | 'height'>;
  readonly panelHeight: number;
  readonly viewportWidth: number;
  readonly viewportHeight: number;
  readonly width: number;
}
