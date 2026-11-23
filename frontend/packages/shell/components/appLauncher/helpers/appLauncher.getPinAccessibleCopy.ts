import type { GetAppLauncherPinAccessibleCopy } from './types/appLauncher.getPinAccessibleCopy';

/** Each runtime message branch uses a literal translation source and named parameters. */
export const getAppLauncherPinAccessibleCopy: GetAppLauncherPinAccessibleCopy = (isPinned, name) => {
  if (isPinned) {
    return {
      accessibleName: __('Unpin {{name}}', { name }),
      tooltip: __('Unpin {{name}} from sidebar', { name }),
    };
  }

  return {
    accessibleName: __('Pin {{name}}', { name }),
    tooltip: __('Pin {{name}} to sidebar', { name }),
  };
};
