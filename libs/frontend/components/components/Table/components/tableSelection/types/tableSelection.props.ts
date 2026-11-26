export interface TableSelectionControlProps {
  readonly checked: boolean;
  readonly indeterminate?: boolean;
  readonly disabled?: boolean;
  readonly label: string;
  readonly onChange: () => void;
}
