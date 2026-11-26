import { useMobileViewport } from '@Components/hooks/useMobileViewport';
import { Filter } from '@Components/icons';

import { BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Tooltip } from '../Tooltip';
import { Popover } from '../Popover';

import { DesktopFilters } from './templates/desktopFilters';
import { MobileFilterEditor } from './components/mobileFilterEditor';
import { MobileFilterList } from './components/mobileFilterList';
import styles from './Filters.module.css';
import type { FiltersProps } from './types/filters.props';
import type { useFilters } from './hooks/useFilters';

export const FiltersView = ({
  fields,
  value,
  disabled,
  iconOnly,
  dataTestId,
  initialOpen,
  controls,
}: FiltersProps & { readonly controls: ReturnType<typeof useFilters> }) => {
  const mobile = useMobileViewport();
  const { editing } = controls;
  const getTitle = () => {
    if (!mobile) return __('Manage filters');
    if (editing) return __('Filter condition');
    return __('Manage filters');
  };
  const getButton = () => {
    if (value.length) return __('Filters ({{count}})', { count: value.length });
    return __('Filters');
  };
  return (
    <Popover
      open={controls.open}
      onOpenChange={controls.onOpenChange}
      title={getTitle()}
      width={760}
      compact
      desktopScrollLayout
      trigger={
        !initialOpen && (
          <Tooltip label={getButton()}>
            <ButtonWithIcon
              icon={Filter}
              className={styles.trigger}
              size={BUTTON_SIZE.SMALL}
              iconSize={14}
              variant={BUTTON_VARIANT.OUTLINE}
              aria-label={getButton()}
              aria-expanded={controls.open}
              disabled={disabled}
              dataTestId={dataTestId}
              onClick={() => controls.onOpenChange(true)}
            >
              {!iconOnly && getButton()}
              {iconOnly && !!value.length && value.length}
            </ButtonWithIcon>
          </Tooltip>
        )
      }
    >
      {!mobile && <DesktopFilters fields={fields} controls={controls} />}
      {mobile && !editing && <MobileFilterList controls={controls} fields={fields} />}
      {mobile && editing && <MobileFilterEditor controls={controls} editing={editing} />}
    </Popover>
  );
};
