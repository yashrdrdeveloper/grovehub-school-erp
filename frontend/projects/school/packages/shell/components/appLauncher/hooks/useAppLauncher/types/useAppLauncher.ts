import type { AppLauncherModel } from '../../../types/appLauncher';
export type UseAppLauncher = (open: boolean, pinned: readonly string[]) => AppLauncherModel;
