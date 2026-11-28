import { DebouncedInput } from '../DebouncedInput';
import { Box } from '../Box';
import type { SearchInputProps } from './types/searchInput.props';

export const SearchInput = ({ dataTestId, ...props }: SearchInputProps) => (
  <Box dataTestId={dataTestId}>
    <DebouncedInput {...props} type="search" dataTestId={dataTestId && `${dataTestId}-input`} />
  </Box>
);
