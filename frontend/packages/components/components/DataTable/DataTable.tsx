import type { ReactNode } from 'react';
import type { DataTableComponent, DataTableProps } from './types/dataTable.props';

import { flexRender, getCoreRowModel, useReactTable } from '@tanstack/react-table';

import { EmptyState } from '@Components/components/EmptyState';

export const DataTable: DataTableComponent = <T,>({
  rows,
  columns,
  emptyTitle = 'Nothing to show',
  dataTestId,
}: DataTableProps<T>) => {
  const table = useReactTable({
    data: rows,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  if (rows.length === 0) {
    return <EmptyState title={emptyTitle} description="Results will appear here when available." />;
  }

  return (
    <div className="table-scroll">
      <table data-testid={dataTestId} className="data-table">
        <thead>
          {table.getHeaderGroups().map((group) => (
            <tr key={group.id}>
              {group.headers.map((header) => {
                if (header.isPlaceholder) return <th key={header.id} />;
                return (
                  <th key={header.id}>
                    {flexRender(header.column.columnDef.header, header.getContext()) as ReactNode}
                  </th>
                );
              })}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext()) as ReactNode}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
