import { useRef, useState } from 'react';
import { _map } from '@Base/utilities/map';
import type { UserRoleAssignment } from '@Business/domain/users';

import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Form } from '@Components/components/Form';
import { Input } from '@Components/components/Input';
import { Sheet } from '@Components/components/Sheet';

import { EditProfile } from './components/editProfile';
import { UserRoleFields } from '../userRoleFields';
import { isRoleAssignmentValid } from '../../helpers/userSetup.isRoleAssignmentValid';
import { USER_SETUP_TEST_IDS } from '../../constants/userSetup.testIds';
import styles from './UserForm.module.css';
import type { UserFormProps } from './types/userForm.props';

export const UserForm = ({ row, roles, rolesLoading, rolesError, retryRoles, profile, profileLoading, profileError, retryProfile, onClose, onCreate, onEdit, onEditProfile, submitting, error, profileSubmitting, profileSubmissionError, profileSaved }: UserFormProps) => {
  const [firstName, setFirstName] = useState('');
  const [middleName, setMiddleName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState(row?.email ?? '');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [username, setUsername] = useState(row?.username ?? '');
  const [assignment, setAssignment] = useState<UserRoleAssignment>({
    primaryRoleId: row?.primaryRoleId ?? '',
    secondaryRoleIds: row?.secondaryRoleIds ?? [],
  });
  const roleOptions = _map(roles, (role) => ({ value: role.id, label: role.name }));
  const rolesValid = isRoleAssignmentValid(assignment, roles);
  const attempt = useRef<{ key: string; hash: string } | null>(null);
  const submit = () => {
    if (!rolesValid) return;
    if (row) {
      const changed = username.trim();
      if (changed && changed !== row.username) {
        onEdit(row, { roleAssignment: assignment, username: changed, version: row.version });
        return;
      }
      onEdit(row, { roleAssignment: assignment, version: row.version });
      return;
    }
    const input = {
      person: { firstName: firstName.trim(), middleName: middleName.trim() || null, lastName: lastName.trim() },
      email: email.trim(),
      phoneNumber: phoneNumber.trim() || undefined,
      roleAssignment: assignment,
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
    if (row) return __('Save account and roles');
    return __('Create user');
  };
  return (
    <Sheet open onOpenChange={(open) => { if (!open) onClose(); }} title={getTitle()} dataTestId={getTestId()}>
      <Box className={styles.sections}>
        {row && profileLoading && <Box role="status">{__('Loading personal information…')}</Box>}
        {row && profileError && <Box role="alert">{__('Unable to load personal information.')} <Button variant={BUTTON_VARIANT.OUTLINE} onClick={retryProfile}>{__('Retry')}</Button></Box>}
        {row && profile && (
          <EditProfile profile={profile} onSave={(input) => onEditProfile(row, input)} submitting={profileSubmitting} error={profileSubmissionError} saved={profileSaved} />
        )}
      <Form onSubmit={(event) => { event.preventDefault(); submit(); }}>
        {row && <Box className={styles.sectionTitle}>{__('Account and roles')}</Box>}
        {!row && (
          <>
            <Input label={__('First name')} required value={firstName} onChange={(event) => setFirstName(event.currentTarget.value)} />
            <Input label={__('Middle name')} value={middleName} onChange={(event) => setMiddleName(event.currentTarget.value)} />
            <Input label={__('Last name')} required value={lastName} onChange={(event) => setLastName(event.currentTarget.value)} />
          </>
        )}
        <Input label={__('Email address')} type="email" required readOnly={Boolean(row)} value={email} onChange={(event) => setEmail(event.currentTarget.value)} />
        {!row && <Input label={__('Phone number')} type="tel" value={phoneNumber} onChange={(event) => setPhoneNumber(event.currentTarget.value)} />}
        {row && <Input label={__('Username')} value={username} onChange={(event) => setUsername(event.currentTarget.value)} />}
        {!row && <Box className={styles.notice}>{__('An invitation email will be sent to set the password.')}</Box>}
        {rolesLoading && <Box role="status">{__('Loading assignable roles…')}</Box>}
        {rolesError && <Box role="alert">{__('Unable to load assignable roles.')} <Button variant={BUTTON_VARIANT.OUTLINE} onClick={retryRoles}>{__('Retry')}</Button></Box>}
        {!rolesLoading && !rolesError && <UserRoleFields options={roleOptions} value={assignment} onChange={setAssignment} disabled={submitting} />}
        <Box className={styles.notice}>{__('School job roles do not grant administrator access.')}</Box>
        {error && <Box role="alert">{__('Unable to save the account. Please try again.')}</Box>}
        <Box className={styles.actions}>
          <Button variant={BUTTON_VARIANT.OUTLINE} onClick={onClose} disabled={submitting || profileSubmitting}>{__('Cancel')}</Button>
          <Button type="submit" loading={submitting} disabled={rolesLoading || rolesError || !rolesValid || profileSubmitting || (!row && (!firstName.trim() || !lastName.trim() || !email.trim()))}>{getLabel()}</Button>
        </Box>
      </Form>
      </Box>
    </Sheet>
  );
};
