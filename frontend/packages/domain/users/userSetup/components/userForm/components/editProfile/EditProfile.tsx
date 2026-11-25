import { useState } from 'react';

import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Form } from '@Components/components/Form';
import { Input } from '@Components/components/Input';

import styles from '../../UserForm.module.css';
import type { EditProfileProps } from './types/editProfile.props';

export const EditProfile = ({ profile, onSave, submitting, error, saved }: EditProfileProps) => {
  const [firstName, setFirstName] = useState(profile.firstName);
  const [middleName, setMiddleName] = useState(profile.middleName ?? '');
  const [lastName, setLastName] = useState(profile.lastName);

  return (
    <Form onSubmit={(event) => {
      event.preventDefault();
      onSave({ firstName: firstName.trim(), middleName: middleName.trim() || null, lastName: lastName.trim() });
    }}>
      <Box className={styles.sectionTitle}>{__('Personal information')}</Box>
      <Input label={__('First name')} required value={firstName} maxLength={100} onChange={(event) => setFirstName(event.currentTarget.value)} />
      <Input label={__('Middle name')} value={middleName} maxLength={100} onChange={(event) => setMiddleName(event.currentTarget.value)} />
      <Input label={__('Last name')} required value={lastName} maxLength={100} onChange={(event) => setLastName(event.currentTarget.value)} />
      <Input label={__('Phone number')} value={profile.phoneNumber ?? ''} readOnly />
      <Box className={styles.notice}>{__('Contact changes require a separate verified contact workflow.')}</Box>
      {saved && <Box role="status">{__('Personal information saved.')}</Box>}
      {error && <Box role="alert">{__('Unable to save personal information.')}</Box>}
      <Box className={styles.actions}>
        <Button variant={BUTTON_VARIANT.OUTLINE} type="submit" loading={submitting} disabled={!firstName.trim() || !lastName.trim()}>{__('Save personal information')}</Button>
      </Box>
    </Form>
  );
};
