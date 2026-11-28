import type { InitI18nStatus } from '../constants/useInitI18n.status';

/** Routes mount only when the initializer reports READY. */
export type UseInitI18n = () => InitI18nStatus;
