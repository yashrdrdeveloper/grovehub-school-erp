export interface DateRangeValue {
  readonly startDate: string | null;
  readonly endDate: string | null;
  readonly presetId?: string;
}

export interface DateRangeProps {
  readonly value: DateRangeValue;
  readonly onChange: (value: DateRangeValue) => void;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
}
