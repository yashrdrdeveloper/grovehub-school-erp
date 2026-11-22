import { Box, BOX_TAG } from '@Components/components/Box';
import { FieldLabel } from '@Components/components/FieldLabel';
import { CREDENTIAL_FIELD_LEGACY_CLASSES } from '../constants/credentialField.legacyClasses';

import type { RenderCredentialFieldLabel } from './types/credentialField.label';

/** Pure grouped-label rendering; input semantics remain owned by CredentialField. */
export const renderCredentialFieldLabel: RenderCredentialFieldLabel = ({
  id,
  label,
  groupLabel,
  labelAction,
}) => {
  const labelNode = <FieldLabel htmlFor={id}>{label}</FieldLabel>;

  if (groupLabel) {
    return (
      <Box as={BOX_TAG.DIV} className={CREDENTIAL_FIELD_LEGACY_CLASSES.GROUPED_LABEL}>
        {labelNode}
        {labelAction}
      </Box>
    );
  }

  return labelNode;
};
