import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { auditApi } from '@Business/domain/platform/audit/services/audit';
import type { AuditActivity } from '@Business/domain/platform/audit/types/contracts';

export const useAuditPage = () => {
  const [searchText, setSearchText] = useState('');
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<AuditActivity | null>(null);
  const query = useQuery({ queryKey: ['audit', searchText, page], queryFn: () => auditApi.search({ searchText, pageInfo: { pageNumber: page, rowSize: 25 }, scope: 'ME' }) });
  const search = (value: string) => { setSearchText(value); setPage(1); };
  return { query, page, setPage, searchText, search, selected, setSelected };
};
