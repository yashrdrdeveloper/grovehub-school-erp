export interface NumberRangeValue {
  readonly min: number | null;
  readonly max: number | null;
}

export interface NumberRangeProps {
  readonly value: NumberRangeValue;
  readonly onChange: (value: NumberRangeValue) => void;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
}
