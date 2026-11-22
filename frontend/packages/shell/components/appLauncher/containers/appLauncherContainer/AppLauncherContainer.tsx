import { AppLauncher } from '../../AppLauncher';
import { useAppLauncher } from '../../hooks/useAppLauncher';
import type { AppLauncherProps } from '../../types/appLauncher';

export const AppLauncherContainer = (props: AppLauncherProps) => {
  const model = useAppLauncher(props.open, props.pinned);
  return <AppLauncher {...props} model={model}/>;
};
