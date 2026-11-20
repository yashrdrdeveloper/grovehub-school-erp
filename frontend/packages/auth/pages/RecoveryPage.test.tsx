import type { ReactNode } from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { ApiError } from '@Business/api';

import { RecoveryPage } from './RecoveryPage';

const validateInvitation = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[string]>();
const exchangeInvitation = jest.fn<Promise<{verified:boolean;purpose:string;expiresIn:number}>,[string]>();
const setup = jest.fn<Promise<unknown>,[string]>();
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
  },
}));
jest.mock('@Business/accounts',()=>({
  accountsApi:{setup:(password:string)=>setup(password)},
}));
jest.mock('./LoginLayout',()=>({
  LoginLayout:({children,title}:{children:ReactNode;title:string})=><main><h1>{title}</h1>{children}</main>,
}));

const renderPage = (purpose:'invitation'|'recovery',proof?:string) => {
  window.history.replaceState({},'',`/${purpose==='invitation'?'set':'reset'}-password${proof?`#proof=${proof}`:''}`);
  return render(<MemoryRouter><RecoveryPage purpose={purpose}/></MemoryRouter>);
};

describe('password link flows',()=>{
  beforeEach(()=>{
    for(const mock of [validateInvitation,exchangeInvitation,setup,validateRecovery,exchangeRecovery,resetPassword]) mock.mockReset();
    validateInvitation.mockResolvedValue(ok);
    exchangeInvitation.mockResolvedValue(ok);
    validateRecovery.mockResolvedValue({...ok,purpose:'PASSWORD_RESET'});
    exchangeRecovery.mockResolvedValue({...ok,purpose:'PASSWORD_RESET'});
  });

  it('shows a missing-link illustration instead of a disabled password form',async()=>{
    renderPage('invitation');
    expect(await screen.findByText('This link looks incomplete')).toBeTruthy();
    expect(screen.queryByLabelText('New password')).toBeNull();
    expect(screen.getByText(/contact your RDR School administrator/i)).toBeTruthy();
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
    fireEvent.click(screen.getByRole('button',{name:'Activate account'}));
    expect(setup).not.toHaveBeenCalled();
    fireEvent.change(screen.getByLabelText('New password'),{target:{value:'a-strong-password-123'}});
    fireEvent.change(screen.getByLabelText('Confirm password'),{target:{value:'a-strong-password-123'}});
    fireEvent.click(screen.getByRole('button',{name:'Activate account'}));
    expect(await screen.findByText('Welcome to RDR School!')).toBeTruthy();
    expect(screen.getByRole('link',{name:/Go to login/i})).toBeTruthy();
  });

  it('replaces the form with an illustration after two failed activation calls',async()=>{
    setup.mockRejectedValue(new Error('temporary failure'));
    renderPage('invitation','proof');
    await screen.findByLabelText('New password');
    for(let attempt=0;attempt<2;attempt++){
      fireEvent.change(screen.getByLabelText('New password'),{target:{value:'a-strong-password-123'}});
      fireEvent.change(screen.getByLabelText('Confirm password'),{target:{value:'a-strong-password-123'}});
      fireEvent.click(screen.getByRole('button',{name:'Activate account'}));
      if(attempt===0) await screen.findByText(/one more attempt/i);
    }
    expect(await screen.findByText('We couldn’t finish your setup')).toBeTruthy();
    expect(screen.queryByRole('button',{name:'Activate account'})).toBeNull();
    expect(setup).toHaveBeenCalledTimes(2);
  });

  it('shows an expired-link illustration on a rejected validation',async()=>{
    validateInvitation.mockRejectedValue(new ApiError(401,'IDENTITY_TOKEN_INVALID','Token is invalid or expired'));
    renderPage('invitation','expired');
    expect(await screen.findByText('This link is no longer valid')).toBeTruthy();
    expect(screen.queryByLabelText('New password')).toBeNull();
  });
});
