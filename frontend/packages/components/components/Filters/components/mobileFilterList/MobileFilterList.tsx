import { _map } from '@Base/utilities/map';
import { _find } from '@Base/utilities/find';
import { Pencil, Plus, Trash2, Filter } from '@Components/icons';

import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../../../Button';
import { ButtonWithIcon } from '../../../ButtonWithIcon';
import { EmptyState } from '../../../EmptyState';
import { getFormattedRangeDate } from '../../../DateRange';
import { Tooltip } from '../../../Tooltip';

import { FILTER_OPERATOR_LABELS } from '../../constants/filters.copy';
import { FilterOperator } from '../../constants/filter.operator';
import { FilterFieldType } from '../../constants/filter.fieldType';
import type { FilterCondition } from '../../types/filters.props';
import styles from '../../Filters.module.css';
import type { MobileFilterListProps } from './types/mobileFilterList.props';

export const MobileFilterList = ({ fields, controls }: MobileFilterListProps) => {
  const getValueSummary = (condition: FilterCondition) => {
    const field = _find(fields, (entry) => entry.id === condition.field);
    if (
      field?.type === FilterFieldType.DATE_RANGE ||
      (field?.type === FilterFieldType.DATE && condition.operator === FilterOperator.BETWEEN)
    ) {
      return _map(condition.values, (value) => getFormattedRangeDate(value)).join(' – ');
    }
    return condition.values.join(', ');
  };
  return (
    <Box className={styles.panel}>
      <Box className={styles.list}>
        {!controls.draft.length && (
          <EmptyState
            compact
            icon={Filter}
            title={__('No filters added')}
            description={__('Add a filter to narrow down your results.')}
          />
        )}
        {_map(controls.draft, (condition) => (
          <Box className={styles.condition} key={condition.id}>
            <Box className={styles.summary}>
              <Box>
                {_find(fields, (field) => field.id === condition.field)?.label ??
                  condition.field}
              </Box>
              <Box className={styles.muted}>
                {FILTER_OPERATOR_LABELS[condition.operator]} {getValueSummary(condition)}
              </Box>
            </Box>
            <ButtonWithIcon
              icon={Pencil}
              variant={BUTTON_VARIANT.GHOST}
              aria-label={__('Edit filter')}
              onClick={() => controls.edit(condition)}
            />
            <ButtonWithIcon
              icon={Trash2}
              variant={BUTTON_VARIANT.GHOST}
              aria-label={__('Delete filter')}
              onClick={() => controls.remove(condition.id)}
            />
          </Box>
        ))}
      </Box>
      <Box className={styles.footer}>
        <Tooltip label={controls.addReason} focusable={Boolean(controls.addReason)}>
          <ButtonWithIcon
            icon={Plus}
            size={BUTTON_SIZE.SMALL}
            iconSize={14}
            variant={BUTTON_VARIANT.GHOST}
            className={styles.add}
            disabled={!controls.canAdd}
            onClick={controls.add}
          >
            {__('Add filter')}
          </ButtonWithIcon>
        </Tooltip>
        <Button
          size={BUTTON_SIZE.SMALL}
          variant={BUTTON_VARIANT.GHOST}
          onClick={controls.reset}
        >
          {__('Reset')}
        </Button>
        <Tooltip
          label={
            (!controls.isDraftValid &&
              __('Complete all filter fields and values before applying.')) ||
            undefined
          }
          focusable={Boolean(
            (!controls.isDraftValid &&
              __('Complete all filter fields and values before applying.')) ||
              undefined,
          )}
        >
          <Button
            size={BUTTON_SIZE.SMALL}
            disabled={!controls.isDraftValid}
            onClick={controls.apply}
          >
            {__('Apply filters')}
          </Button>
        </Tooltip>
      </Box>
    </Box>
  );
};
