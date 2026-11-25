export interface CalendarProps {
  readonly value: string | null;
  readonly onChange: (value: string | null) => void;
  readonly min?: string;
  readonly max?: string;
  readonly disabled?: boolean;
  readonly clearable?: boolean;
  readonly dataTestId?: string;
}
