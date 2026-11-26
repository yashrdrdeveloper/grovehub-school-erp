import { useRef, useState } from 'react';

import { UserStatus } from '@Business/domain/users';
import { UserAction } from '../../constants/userSetup.actions';
import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Form } from '@Components/components/Form';
import { Input } from '@Components/components/Input';
import { Sheet } from '@Components/components/Sheet';

import { USER_SETUP_TEST_IDS } from '../../constants/userSetup.testIds';
import styles from './UserForm.module.css';
import type { UserFormProps } from './types/userForm.props';

export const UserForm = ({
  row,
  onClose,
  onCreate,
  onEdit,
  onAccountAction,
  submitting,
  error,
}: UserFormProps) => {
  const [firstName, setFirstName] = useState('');
  const [middleName, setMiddleName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState(row?.email ?? '');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [username, setUsername] = useState(row?.username ?? '');
  const attempt = useRef<{ key: string; hash: string } | null>(null);
  const submit = () => {
    if (row) {
      onEdit(row, username.trim() || null);
      return;
    }
    const input = {
      person: {
        firstName: firstName.trim(),
        middleName: middleName.trim() || null,
        lastName: lastName.trim(),
      },
      email: email.trim(),
      phoneNumber: phoneNumber.trim() || undefined,
    };
    const hash = JSON.stringify(input);
    if (attempt.current?.hash !== hash) attempt.current = { key: crypto.randomUUID(), hash };
    onCreate(input, attempt.current.key);
  };
  const getTitle = () => {
    if (row) return __('Edit user');
    return __('Create user');
  };
  const getTestId = () => {
    if (row) return USER_SETUP_TEST_IDS.EDIT;
    return USER_SETUP_TEST_IDS.CREATE;
  };
  const getLabel = () => {
    if (row) return __('Save changes');
    return __('Create user');
  };
  return (
    <Sheet
      open
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
      title={getTitle()}
      dataTestId={getTestId()}
    >
      <Form
        onSubmit={(event) => {
          event.preventDefault();
          submit();
        }}
      >
        {!row && (
          <>
            <Input
              label={__('First name')}
              required
              value={firstName}
              onChange={(event) => setFirstName(event.currentTarget.value)}
            />
            <Input
              label={__('Middle name')}
              value={middleName}
              onChange={(event) => setMiddleName(event.currentTarget.value)}
            />
            <Input
              label={__('Last name')}
              required
              value={lastName}
              onChange={(event) => setLastName(event.currentTarget.value)}
            />
          </>
        )}
        <Input
          label={__('Email address')}
          type="email"
          required
          readOnly={Boolean(row)}
          value={email}
          onChange={(event) => setEmail(event.currentTarget.value)}
        />
        {!row && (
          <Input
            label={__('Phone number')}
            type="tel"
            value={phoneNumber}
            onChange={(event) => setPhoneNumber(event.currentTarget.value)}
          />
        )}
        {row && (
          <Input
            label={__('Username')}
            value={username}
            onChange={(event) => setUsername(event.currentTarget.value)}
          />
        )}
        {!row && (
          <Box className={styles.notice}>
            {__('An invitation email will be sent to set the password.')}
          </Box>
        )}
        <Box className={styles.notice}>
          {__('Role assignment and editing personal details require additional backend support.')}
        </Box>
        {row?.status === UserStatus.SUSPENDED && (
          <Button
            variant={BUTTON_VARIANT.OUTLINE}
            disabled={submitting}
            onClick={() => onAccountAction?.(UserAction.ACTIVATE, row)}
          >
            {__('Activate account')}
          </Button>
        )}
        {row?.status === UserStatus.PENDING_SETUP && (
          <Button
            variant={BUTTON_VARIANT.OUTLINE}
            disabled={submitting}
            onClick={() => onAccountAction?.(UserAction.RESEND, row)}
          >
            {__('Resend invitation')}
          </Button>
        )}
        {error && <Box role="alert">{__('Unable to save the account. Please try again.')}</Box>}
        <Box className={styles.actions}>
          <Button variant={BUTTON_VARIANT.OUTLINE} onClick={onClose} disabled={submitting}>
            {__('Cancel')}
          </Button>
          <Button
            type="submit"
            loading={submitting}
            disabled={!row && (!firstName.trim() || !lastName.trim() || !email.trim())}
          >
            {getLabel()}
          </Button>
        </Box>
      </Form>
    </Sheet>
  );
};
