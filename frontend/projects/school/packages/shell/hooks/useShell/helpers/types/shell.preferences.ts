export type ReadPinnedApps = () => string[];
export type WritePinnedApps = (paths: string[]) => void;
export type ReadSidebarCollapsed = () => boolean;
export type WriteSidebarCollapsed = (collapsed: boolean) => void;
