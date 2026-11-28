import { Typography } from '../Typography';
import { TEXT_TAG } from './constants/text.tag';
import type { TextComponent } from './types/text.props';
export const Text: TextComponent = ({ as = TEXT_TAG.SPAN, dataTestId, ...props }) => (
  <Typography {...props} as={as} dataTestId={dataTestId} />
);
