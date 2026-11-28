export interface MobileInfiniteEndProps {
  readonly hasNext: boolean;
  readonly loading: boolean;
  readonly error?: boolean;
  readonly onLoadMore?: () => void;
}
