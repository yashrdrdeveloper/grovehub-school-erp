/** Closed bootstrap states; shell routes may render only after READY. */
export enum InitI18nStatus {
  INITIALIZING = 'initializing',
  READY = 'ready',
  FAILED = 'failed',
}
