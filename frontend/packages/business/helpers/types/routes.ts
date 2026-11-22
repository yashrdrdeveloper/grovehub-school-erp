import type { RouteConfig, RouteEntity } from '../../baseEntities/routes';

export type GetRoute = <Id extends string>(routes: RouteConfig<Id>, id: Id) => RouteEntity<Id>;
