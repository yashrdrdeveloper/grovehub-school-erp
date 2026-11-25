export interface DateRangePreset {
  readonly id: string;
  readonly label: string;
  readonly start: (today: Date) => Date;
  readonly end: (today: Date) => Date;
}
