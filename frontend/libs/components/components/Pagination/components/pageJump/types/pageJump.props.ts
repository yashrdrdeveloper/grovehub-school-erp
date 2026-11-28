export interface PageJumpProps {
  readonly currentPage: number;
  readonly totalPages: number;
  readonly disabled: boolean;
  readonly onChange: (page: number) => void;
}
