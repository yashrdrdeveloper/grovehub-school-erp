export interface AppTestIds {
  readonly card: string;
  readonly pin: string;
}
export type GetAppTestIds = (path: string) => AppTestIds;
