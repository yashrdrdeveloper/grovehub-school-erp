import { Typography } from '../Typography';
import type { TextComponent } from './types/text.props';
export const Text: TextComponent = ({ as = 'span', dataTestId, ...props }) => (
  <Typography {...props} as={as} dataTestId={dataTestId} />
);
