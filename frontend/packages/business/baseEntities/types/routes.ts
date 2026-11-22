import type { RouteType } from '../constants/route.type';

export interface RouteEntity<Id extends string = string> {
  readonly id: Id;
  readonly name: string;
  readonly path: string;
  readonly type: RouteType;
}

export type RouteConfig<Id extends string = string> = readonly RouteEntity<Id>[];
