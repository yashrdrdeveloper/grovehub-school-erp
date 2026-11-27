import { useCallback, useState } from 'react';

import type { DashboardCalendarView } from '../../types';
import type { UseDashboardCalendar } from './types/useDashboardCalendar';

/** Calendar selection state, intentionally independent of real event loading. */
export const useDashboardCalendar: UseDashboardCalendar = () => {
  const [view, setView] = useState<DashboardCalendarView>('today');
  const showToday = useCallback(() => setView('today'), []);
  const showUpcoming = useCallback(() => setView('upcoming'), []);
  return { view, showToday, showUpcoming };
};
