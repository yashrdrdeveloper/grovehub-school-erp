import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { rbacApi } from '@Business/domain/rbac/services/rbac';
import type { Role } from '@Business/domain/rbac/types/contracts';

export const useRbacPage = () => {
  const [searchText, setSearchText] = useState('');
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<Role | null>(null);
  const query = useQuery({ queryKey: ['rbac', searchText, page], queryFn: () => rbacApi.search({ searchText, pageInfo: { pageNumber: page, rowSize: 25 } }) });
  const permissions = useQuery({ queryKey: ['role-permissions', selected?.roleId], enabled: Boolean(selected), queryFn: () => rbacApi.permissions(selected?.roleId ?? '') });
  const search = (value: string) => { setSearchText(value); setPage(1); };
  return { query, page, setPage, searchText, search, selected, setSelected, permissions };
};
