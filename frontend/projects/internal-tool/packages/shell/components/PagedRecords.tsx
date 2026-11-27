import { useEffect, useState, type ReactNode } from 'react';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { RefreshCw, ChevronLeft, ChevronRight } from 'lucide-react';
import { __ } from '@Base/i18n';
import { errorText, search } from '../../business/api/client';
import type { Page } from '../../business/api/types';

export interface Column<T> { title: string; cell: (record: T) => ReactNode }
interface Props<T> { resource: string; revision: number; columns: Column<T>[]; rowKey: (record: T) => string | number; onSelect?: (record: T) => void; title: string }
export function PagedRecords<T>({ resource, revision, columns, rowKey, onSelect, title }: Props<T>) {
  const [pageNumber, setPageNumber] = useState(1);
  const [rowSize, setRowSize] = useState(25);
  const [refresh, setRefresh] = useState(0);
  const [page, setPage] = useState<Page<T> | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const controller = new AbortController();
    setLoading(true); setError(''); setPage(null);
    void search<T>(resource, pageNumber, rowSize, controller.signal).then((value) => {
      if (!controller.signal.aborted) setPage(value);
    }).catch((failure) => {
      if (!controller.signal.aborted) setError(errorText(failure));
    }).finally(() => { if (!controller.signal.aborted) setLoading(false); });
    return () => controller.abort();
  }, [resource, pageNumber, rowSize, refresh, revision]);
  return <section className="record-card" data-testid={resource + '-records'}>
    <div className="record-toolbar"><strong>{title}</strong><Button variant={BUTTON_VARIANT.GHOST} onClick={() => setRefresh((n) => n + 1)} aria-label={__('Refresh')} dataTestId={resource + '-refresh'}><RefreshCw size={16} /></Button></div>
    {error && <p className="error" role="alert">{error}</p>}
    {loading && <p className="empty" role="status">{__('Loading records…')}</p>}
    {!loading && page?.records.length === 0 && <p className="empty">{__('No records found.')}</p>}
    {!!page?.records.length && <div className="table-scroll"><table><thead><tr>{columns.map((column) => <th key={column.title}>{column.title}</th>)}{onSelect && <th>{__('Actions')}</th>}</tr></thead><tbody>{page.records.map((record) => <tr key={rowKey(record)}>{columns.map((column) => <td key={column.title}>{column.cell(record)}</td>)}{onSelect && <td><Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => onSelect(record)} dataTestId={resource + '-open-' + rowKey(record)}>{__('Open')}</Button></td>}</tr>)}</tbody></table></div>}
    <footer className="record-footer"><label>{__('Rows per page')} <select value={rowSize} onChange={(event) => { setRowSize(Number(event.target.value)); setPageNumber(1); }} data-testid={resource + '-row-size'}>{[10,25,50,100].map((n) => <option key={n}>{n}</option>)}</select></label>
      <div className="inline"><Button variant={BUTTON_VARIANT.OUTLINE} disabled={pageNumber <= 1 || loading} onClick={() => setPageNumber((n) => n - 1)} aria-label={__('Previous page')} dataTestId={resource + '-prev'}><ChevronLeft size={16}/></Button><span>{__('Page')} {pageNumber}</span><Button variant={BUTTON_VARIANT.OUTLINE} disabled={!page?.hasMore || loading} onClick={() => setPageNumber((n) => n + 1)} aria-label={__('Next page')} dataTestId={resource + '-next'}><ChevronRight size={16}/></Button></div>
    </footer>
  </section>;
}
