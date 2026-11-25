export interface CalendarProps {
  readonly value: string | null;
  readonly month?: Date;
  readonly onMonthChange?: (month: Date) => void;
  readonly range?: { readonly startDate: string | null; readonly endDate: string | null };
  readonly onDayHover?: (date: string | null) => void;
  readonly compact?: boolean;
  readonly hideFooter?: boolean;
  readonly navigation?: 'both' | 'previous' | 'next';
  readonly onChange: (value: string | null) => void;
  readonly min?: string;
  readonly max?: string;
  readonly disabled?: boolean;
  readonly clearable?: boolean;
  readonly dataTestId?: string;
}
