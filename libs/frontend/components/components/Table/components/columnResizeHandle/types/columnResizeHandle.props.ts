export interface ColumnResizeHandleProps {
  readonly label: string;
  readonly width: number;
  readonly min: number;
  readonly max: number;
  readonly onResize: (width: number) => void;
}
