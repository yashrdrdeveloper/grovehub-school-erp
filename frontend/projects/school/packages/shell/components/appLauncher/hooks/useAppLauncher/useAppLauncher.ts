import { useCallback, useEffect, useState, type ChangeEvent } from 'react';
import { filterCatalog } from '../../helpers/appLauncher.filterCatalog';
import type { UseAppLauncher } from './types/useAppLauncher';

/** The only owner of search query and open/reset behavior. */
export const useAppLauncher: UseAppLauncher = (open, pinned) => {
  const [filter, setFilter] = useState('');
  useEffect(() => {
    if (open) setFilter('');
  }, [open]);
  const onFilterChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setFilter(event.currentTarget.value);
  }, []);
  return { filter, ...filterCatalog(filter, pinned), onFilterChange };
};
