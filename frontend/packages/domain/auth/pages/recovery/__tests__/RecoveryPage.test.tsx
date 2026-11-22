import type { ReactNode } from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { ApiError } from '@Business/api';

import { RecoveryContainer as RecoveryPage } from '..';
import { clearRecoveryFlow } from '../../../core/session/recoveryFlowSession';

const validateInvitation = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[string]>();
const exchangeInvitation = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[string]>();
const setup = jest.fn<Promise<unknown>,[string]>();
const setupStatus = jest.fn<Promise<{verified:boolean;expiresIn:number}>,[]>();
const resetGrantStatus = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[]>();
const validateRecovery = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[string]>();
const exchangeRecovery = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[string]>();
const resetPassword = jest.fn<Promise<unknown>,[string]>();
const ok = { verified:true, purpose:'ACCOUNT_SETUP', expiresIn:600 };

jest.mock('@Business/auth',()=>({
  authenticationApi:{
    validateInvitation:(proof:string)=>validateInvitation(proof),
    exchangeInvitation:(proof:string)=>exchangeInvitation(proof),
    validateRecovery:(proof:string)=>validateRecovery(proof),
    exchangeRecovery:(proof:string)=>exchangeRecovery(proof),
    resetPassword:(password:string)=>resetPassword(password),
    resetGrantStatus:()=>resetGrantStatus(),
  },
}));
jest.mock('@Business/accounts',()=>({
  accountsApi:{setup:(password:string)=>setup(password),setupStatus:()=>setupStatus()},
}));
jest.mock('../../../components/loginLayout',()=>({
  LoginLayout:({children,title,mobileTitle}:{children:ReactNode;title:string;mobileTitle?:string})=><main><h1>{title}</h1>{mobileTitle && <span data-testid="mobile-recovery-title">{mobileTitle}</span>}{children}</main>,
}));

const renderPage = (purpose:'invitation'|'recovery',proof?:string) => {
  window.history.replaceState({},'',`/${purpose==='invitation'?'set':'reset'}-password${proof?`#proof=${proof}`:''}`);
  return render(<MemoryRouter><RecoveryPage purpose={purpose}/></MemoryRouter>);
};

describe('RecoveryContainer',()=>{
  beforeEach(()=>{
    clearRecoveryFlow();
    window.sessionStorage.clear();
    for(const mock of [validateInvitation,exchangeInvitation,setup,setupStatus,validateRecovery,exchangeRecovery,resetPassword,resetGrantStatus]) mock.mockReset();
    validateInvitation.mockResolvedValue(ok);
    exchangeInvitation.mockResolvedValue(ok);
    validateRecovery.mockResolvedValue({...ok,purpose:'PASSWORD_RESET'});
    exchangeRecovery.mockResolvedValue({...ok,purpose:'PASSWORD_RESET'});
    setupStatus.mockRejectedValue(new ApiError(401,'IDENTITY_TOKEN_INVALID','Invalid grant'));
    resetGrantStatus.mockRejectedValue(new ApiError(401,'IDENTITY_TOKEN_INVALID','Invalid grant'));
  });

  it.each([
    ['invitation', 'contact your RDR School administrator'],
    ['recovery', 'request another password reset link'],
  ] as const)('shows missing-link state without calling APIs for %s', async (purpose, helpText) => {
    renderPage(purpose);
    expect(await screen.findAllByText('This link looks incomplete')).toBeTruthy();
    expect(screen.getByTestId('mobile-recovery-title').textContent).toBe(purpose === 'invitation' ? 'Set password' : 'Reset password');
    expect(screen.queryByLabelText('New password')).toBeNull();
    expect(screen.getByText(new RegExp(helpText, 'i'))).toBeTruthy();
    if (purpose === 'recovery') expect(screen.getByRole('link', {name:/Request new link/i}).getAttribute('href')).toBe('/forgot-password');
    for (const mock of [
      validateInvitation, exchangeInvitation, setupStatus,
      validateRecovery, exchangeRecovery, resetGrantStatus,
      setup, resetPassword,
    ]) expect(mock).not.toHaveBeenCalled();
    expect(window.location.hash).toBe('');
  });

  it('does not request grant status without an existing session marker', async () => {
    setupStatus.mockResolvedValue({verified:true,expiresIn:480});
    renderPage('invitation');
    expect(await screen.findAllByText('This link looks incomplete')).toBeTruthy();
    expect(setupStatus).not.toHaveBeenCalled();
  });

  it('resumes an exchanged invitation after refresh using the setup grant cookie', async () => {
    setupStatus.mockResolvedValue({verified:true,expiresIn:480});
    const initial = renderPage('invitation','first-secret');
    await screen.findByLabelText('New password');
    expect(window.location.hash).toBe('');
    initial.unmount();
    renderPage('invitation');
    expect(await screen.findByLabelText('New password')).toBeTruthy();
    expect(setupStatus).toHaveBeenCalledTimes(1);
    expect(validateInvitation).toHaveBeenCalledTimes(1);
    expect(exchangeInvitation).toHaveBeenCalledTimes(1);
  });

  it('resumes a reset-password grant after refresh without replaying the proof', async () => {
    resetGrantStatus.mockResolvedValue({...ok,purpose:'PASSWORD_RESET'});
    const initial = renderPage('recovery','reset-secret');
    await screen.findByLabelText('New password');
    initial.unmount();
    renderPage('recovery');
    expect(await screen.findByLabelText('New password')).toBeTruthy();
    expect(resetGrantStatus).toHaveBeenCalledTimes(1);
    expect(validateRecovery).toHaveBeenCalledTimes(1);
  });

  it('accepts a newly pasted proof in the same mounted set-password page', async () => {
    renderPage('invitation','original');
    await screen.findByLabelText('New password');
    window.history.replaceState({}, '', '/set-password#proof=replaced');
    window.dispatchEvent(new Event('hashchange'));
    await waitFor(() => expect(validateInvitation).toHaveBeenCalledWith('replaced'));
    await waitFor(() => expect(exchangeInvitation).toHaveBeenCalledWith('replaced'));
    expect(window.location.hash).toBe('');
  });

  it('rejects expired session markers without contacting the backend', async () => {
    window.sessionStorage.setItem('groovehub.identity.password-flow.v1',JSON.stringify({
      purpose:'invitation', stage:'grant', expiresAt:Date.now()-1000,
    }));
    renderPage('invitation');
    expect(await screen.findAllByText('This link looks incomplete')).toBeTruthy();
    expect(setupStatus).not.toHaveBeenCalled();
    expect(window.sessionStorage.length).toBe(0);
  });

  it('verifies the proof, strips it from the URL and asks for passwords only',async()=>{
    renderPage('invitation','secret-value');
    await waitFor(()=>expect(screen.getByLabelText('New password')).toBeTruthy());
    expect(validateInvitation).toHaveBeenCalledWith('secret-value');
    expect(exchangeInvitation).toHaveBeenCalledWith('secret-value');
    expect(window.location.hash).toBe('');
    expect(screen.queryByLabelText('Username')).toBeNull();
  });

  it('does not spend an attempt on browser-side validation and celebrates success',async()=>{
    setup.mockResolvedValue({setupCompleted:true});
    renderPage('invitation','proof');
    await screen.findByLabelText('New password');
    fireEvent.click(screen.getByRole('button',{name:'Update password'}));
    expect(setup).not.toHaveBeenCalled();
    fireEvent.change(screen.getByLabelText('New password'),{target:{value:'a-strong-password-123'}});
    fireEvent.change(screen.getByLabelText('Confirm password'),{target:{value:'a-strong-password-123'}});
    fireEvent.click(screen.getByRole('button',{name:'Update password'}));
    expect(await screen.findByText('Welcome to RDR School!')).toBeTruthy();
    expect(screen.getByTestId('mobile-recovery-title').textContent).toBe('Set password');
    expect(screen.getByRole('link',{name:/Go to login/i})).toBeTruthy();
  });

  it('keeps transient activation failures retryable without guessing the budget', async () => {
    setup.mockRejectedValue(new Error('temporary failure'));
    renderPage('invitation', 'proof');
    await screen.findByLabelText('New password');
    for (let attempt = 0; attempt < 2; attempt++) {
      fireEvent.change(screen.getByLabelText('New password'), { target: { value: 'a-strong-password-123' } });
      fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'a-strong-password-123' } });
      fireEvent.click(screen.getByRole('button', { name: 'Update password' }));
      await waitFor(() => expect(setup).toHaveBeenCalledTimes(attempt + 1));
      await screen.findByText('Unable to activate your account right now. Please try again.');
      expect(screen.getByRole('button', { name: 'Update password' })).toBeTruthy();
    }
  });

  it('shows the exhausted state only when reported by the backend', async () => {
    setup.mockRejectedValue(new ApiError(429, 'IDENTITY_SETUP_ATTEMPTS_EXCEEDED', 'Attempts exhausted'));
    renderPage('invitation', 'proof');
    await screen.findByLabelText('New password');
    fireEvent.change(screen.getByLabelText('New password'), { target: { value: 'a-strong-password-123' } });
    fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'a-strong-password-123' } });
    fireEvent.click(screen.getByRole('button', { name: 'Update password' }));
    expect(await screen.findByText('We couldn’t finish your setup')).toBeTruthy();
    expect(screen.queryByRole('button', { name: 'Update password' })).toBeNull();
    expect(setup).toHaveBeenCalledTimes(1);
  });

  it('shows an expired-link illustration on a rejected validation',async()=>{
    validateInvitation.mockRejectedValue(new ApiError(401,'IDENTITY_TOKEN_INVALID','Token is invalid or expired'));
    renderPage('invitation','expired');
    expect(await screen.findByText('This link is no longer valid')).toBeTruthy();
    expect(screen.queryByLabelText('New password')).toBeNull();
  });
});
