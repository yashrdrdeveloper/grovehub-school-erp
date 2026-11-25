export interface QuickFilterOption {
  readonly id: string;
  readonly label: string;
  readonly count?: number;
  readonly disabled?: boolean;
}

export interface QuickFiltersProps {
  readonly options: readonly QuickFilterOption[];
  readonly value: readonly string[];
  readonly multiple?: boolean;
  readonly onChange: (value: readonly string[]) => void;
  readonly dataTestId?: string;
}
