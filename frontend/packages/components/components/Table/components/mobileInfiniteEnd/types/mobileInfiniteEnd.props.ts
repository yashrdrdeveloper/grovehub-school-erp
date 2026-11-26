export interface MobileInfiniteEndProps {
  readonly hasNext: boolean;
  readonly loading: boolean;
  readonly error?: boolean;
  readonly loaded: number;
  readonly total?: number;
  readonly onLoadMore?: () => void;
}
