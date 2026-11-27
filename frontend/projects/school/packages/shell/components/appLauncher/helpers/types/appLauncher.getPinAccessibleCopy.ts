export interface AppLauncherPinAccessibleCopy {
  readonly accessibleName: string;
  readonly tooltip: string;
}

export type GetAppLauncherPinAccessibleCopy = (
  isPinned: boolean,
  name: string,
) => AppLauncherPinAccessibleCopy;
