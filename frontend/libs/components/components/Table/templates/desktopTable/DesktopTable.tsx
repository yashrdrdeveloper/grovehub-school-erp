import { Fragment, useEffect, useId, useRef, useState } from 'react';
import { _map } from '@Base/utilities/map';

import { Box } from '../../../Box';

import { TableScrollbar } from '../../components/tableScrollbar';
import { TABLE_DIMENSIONS } from '../../constants/table.dimensions';
import { getTableFillWidth } from '../../helpers/table.getFillWidth';
import { getTableStretchLayout } from '../../helpers/table.getStretchLayout';
import { TablePin, TableType } from '../../constants/table.type';
import { getVirtualWindow } from '../../helpers/table.getVirtualWindow';

import { DesktopTableHead } from './components/desktopTableHead';
import { DesktopTableBody } from './components/desktopTableBody';
import type { DesktopTableProps } from './types/desktopTable.props';
import styles from '../../Table.module.css';

export const DesktopTable = <T,>(props: DesktopTableProps<T>) => {
  const { controls } = props;
  const [scrollTop, setScrollTop] = useState(0);
  const frame = useRef<HTMLElement>(null);
  const viewportId = useId();
  const [available, setAvailable] = useState(0);
  useEffect(() => {
    const element = frame.current;
    if (!element) return;
    const measure = () => setAvailable(element.clientWidth);
    measure();
    if (typeof ResizeObserver === 'undefined') return;
    const observer = new ResizeObserver(measure);
    observer.observe(element);
    return () => observer.disconnect();
  }, []);
  const virtual =
    props.type === TableType.VIRTUALIZED ||
    ('virtualizationProps' in props && Boolean(props.virtualizationProps));
  const settings = ('virtualizationProps' in props && props.virtualizationProps) || {};
  const height = Math.max(100, settings.height ?? 480);
  const rowHeight = Math.max(44, settings.rowHeight ?? 64);
  const range = getVirtualWindow(
    controls.entries.length,
    scrollTop,
    height,
    rowHeight,
    settings.overscan ?? 6,
  );
  const getEntries = () => {
    if (virtual) return controls.entries.slice(range.start, range.end);
    return controls.entries;
  };
  const selection = Boolean(props.selectionProps);
  const rowActions = Boolean(props.rowActions?.length);
  const fixedWidth =
    Number(selection) * TABLE_DIMENSIONS.selection + Number(rowActions) * TABLE_DIMENSIONS.actions;
  const layout = getTableStretchLayout(
    controls.layout,
    available,
    fixedWidth,
    controls.preferences.widths,
  );
  const columnsWidth = layout.reduce((sum, item) => sum + item.width, 0);
  const fill = getTableFillWidth(available, columnsWidth, fixedWidth);
  const filler = fill > 0;
  const firstRight = layout.findIndex((entry) => entry.pin === TablePin.RIGHT);
  const fillBefore = (index: number) => filler && firstRight === index;
  const fillAtEnd = filler && firstRight < 0;

  const span = Math.max(1, layout.length + Number(selection) + Number(rowActions) + Number(filler));
  const tableWidth = columnsWidth + fixedWidth + fill;
  const frameStyle = () => {
    if (virtual) return { height };
    return undefined;
  };
  const rowStyle = () => {
    if (virtual) return { height: rowHeight };
    return undefined;
  };
  return (
    <Box className={styles.scrollFrame}>
      <Box
        id={viewportId}
        ref={frame}
        className={styles.scroll}
        style={frameStyle()}
        onScroll={(event) => setScrollTop(event.currentTarget.scrollTop)}
      >
        <table
          className={styles.table}
          data-density={controls.preferences.density}
          data-virtual={virtual}
          style={{ width: tableWidth }}
          aria-busy={props.refreshing}
          aria-label={props.headerProps?.title ?? __('Table')}
        >
          <colgroup>
            {selection && <col style={{ width: TABLE_DIMENSIONS.selection }} />}
            {_map(layout, (entry, index) => (
              <Fragment key={entry.column.id}>
                {fillBefore(index) && <col style={{ width: fill }} />}
                <col style={{ width: entry.width }} />
              </Fragment>
            ))}
            {fillAtEnd && <col style={{ width: fill }} />}
            {rowActions && <col style={{ width: TABLE_DIMENSIONS.actions }} />}
          </colgroup>
          <DesktopTableHead
            props={props}
            layout={layout}
            selection={selection}
            rowActions={rowActions}
            fill={fill}
            fillBefore={fillBefore}
            fillAtEnd={fillAtEnd}
          />
          <DesktopTableBody
            props={props}
            entries={getEntries()}
            layout={layout}
            selection={selection}
            rowActions={rowActions}
            fillBefore={fillBefore}
            fillAtEnd={fillAtEnd}
            span={span}
            virtual={virtual}
            range={range}
            rowStyle={rowStyle}
          />
        </table>
      </Box>
      <TableScrollbar viewport={frame} viewportId={viewportId} axis="vertical" />
      <TableScrollbar viewport={frame} viewportId={viewportId} axis="horizontal" />
    </Box>
  );
};
