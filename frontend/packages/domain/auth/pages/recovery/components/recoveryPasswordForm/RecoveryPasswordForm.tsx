import { Box, BOX_TAG } from '@Components/components/Box';
import { AlertCircle, Eye, EyeOff, KeyRound, LockKeyhole, ShieldCheck } from '@Components/icons';
import { RECOVERY_TEST_IDS } from '../../constants/recovery.testIds';
import type { ChangeEventHandler } from 'react';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '@Components/components/Button';
import { Form } from '@Components/components/Form';

import { CredentialField } from '../../../../components/credentialField';
import type { RecoveryPasswordFormComponent } from './types/recoveryPasswordForm';

/** Password-reset and invitation form with no backend calls or flow orchestration. */
export const RecoveryPasswordForm: RecoveryPasswordFormComponent = (props) => {
  const {
    id,
    view,
    password,
    confirmation,
    visible,
    fieldErrors,
    error,
    onSubmit,
    onPasswordChange,
    onConfirmationChange,
    onToggleVisibility,
  } = props;
  const handlePasswordChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onPasswordChange(event.currentTarget.value);
  };
  const handleConfirmationChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onConfirmationChange(event.currentTarget.value);
  };
  let VisibilityIcon = Eye;
  let passwordType: 'text' | 'password' = 'password';
  let visibilityLabel = 'Show passwords';
  let confirmationLabel = 'Show confirmation password';
  if (visible) {
    VisibilityIcon = EyeOff;
    passwordType = 'text';
    visibilityLabel = 'Hide passwords';
    confirmationLabel = 'Hide confirmation password';
  }
  const saving = view === 'saving';
  let buttonLabel = 'Update password';
  if (saving) buttonLabel = 'Securing your account…';

  return (
    <Form
      unstyled
      dataTestId={RECOVERY_TEST_IDS.FORM}
      className="rdr-login-v2__form rdr-recovery__form"
      onSubmit={onSubmit}
      noValidate
    >
      <Box as={BOX_TAG.DIV} className="rdr-recovery__fields">
        <CredentialField
          dataTestId={RECOVERY_TEST_IDS.PASSWORD}
          id={id + '-password'}
          label="New password"
          Icon={LockKeyhole}
          type={passwordType}
          autoComplete="new-password"
          placeholder="New password"
          value={password}
          minLength={12}
          disabled={saving}
          required
          error={fieldErrors.password}
          onChange={handlePasswordChange}
          trailingAction={
            <Button
              unstyled
              type="button"
              className="rdr-login-v2__visibility"
              data-testid={RECOVERY_TEST_IDS.PASSWORD_VISIBILITY}
              disabled={saving}
              aria-label={visibilityLabel}
              onClick={onToggleVisibility}
            >
              <VisibilityIcon size={19} aria-hidden="true" />
            </Button>
          }
        />

        <CredentialField
          dataTestId={RECOVERY_TEST_IDS.CONFIRMATION}
          id={id + '-confirmation'}
          label="Confirm password"
          Icon={KeyRound}
          type={passwordType}
          autoComplete="new-password"
          placeholder="Confirm password"
          value={confirmation}
          disabled={saving}
          required
          error={fieldErrors.confirmation}
          onChange={handleConfirmationChange}
          trailingAction={
            <Button
              unstyled
              type="button"
              className="rdr-login-v2__visibility"
              disabled={saving}
              aria-label={confirmationLabel}
              onClick={onToggleVisibility}
            >
              <VisibilityIcon size={19} aria-hidden="true" />
            </Button>
          }
        />
        <Box as={BOX_TAG.P} className="rdr-recovery__requirement">
          <ShieldCheck size={17} /> Use at least 12 characters to keep your account secure.
        </Box>
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-recovery__actions">
        <Button
          type="submit"
          variant={BUTTON_VARIANT.PRIMARY}
          size={BUTTON_SIZE.LARGE}
          dataTestId={RECOVERY_TEST_IDS.SUBMIT}
          className="rdr-login-v2__submit"
          loading={saving}
          loadingText={buttonLabel}
        >
          {buttonLabel}
        </Button>
        <Box as={BOX_TAG.SPAN} className="rdr-recovery__remaining" aria-hidden="true">
          {'\u00a0'}
        </Box>
        <Box as={BOX_TAG.DIV} className="rdr-recovery__error-slot" aria-live="polite">
          {error && (
            <Box as={BOX_TAG.P} role="alert" className="rdr-recovery__inline-error">
              <AlertCircle size={16} />
              {error}
            </Box>
          )}
        </Box>
      </Box>
    </Form>
  );
};
