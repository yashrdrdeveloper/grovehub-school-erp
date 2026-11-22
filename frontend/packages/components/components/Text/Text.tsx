import { Typography, TYPOGRAPHY_TAG } from '../Typography';
import type { TextComponent } from './types/text.props';
export const Text: TextComponent = ({ as = TYPOGRAPHY_TAG.SPAN, dataTestId, ...props }) => (
  <Typography {...props} as={as} dataTestId={dataTestId} />
);
