import type { ReactNode } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

import { ForgotPasswordContainer as ForgotPasswordPage } from '..';

const mockForgotPassword = jest.fn<Promise<{accepted:boolean;delivery:string}>,[string]>();

jest.mock('@Business/auth',()=>({
  authenticationApi:{forgotPassword:(email:string)=>mockForgotPassword(email)},
}));
jest.mock('../../../components/loginLayout',()=>({
  LoginLayout:({children,title}:{children:ReactNode;title:string})=><div><h1>{title}</h1>{children}</div>,
}));

const renderRecovery=()=>render(<MemoryRouter initialEntries={['/forgot-password']}>
  <Routes>
    <Route path="/forgot-password" element={<ForgotPasswordPage/>}/>
    <Route path="/login" element={<p>Sign in page</p>}/>
  </Routes>
</MemoryRouter>);

describe('ForgotPasswordContainer',()=>{
  beforeEach(()=>{
    mockForgotPassword.mockReset();
    jest.useRealTimers();
  });
  afterEach(()=>jest.useRealTimers());

  it('shows guidance, validates email and preserves its action',()=>{
    renderRecovery();
    expect(screen.getByRole('heading',{name:'Forgot your password?'})).toBeTruthy();
    expect(screen.getByText('Secure password reset')).toBeTruthy();
    fireEvent.click(screen.getByRole('button',{name:'Send reset link'}));
    expect(screen.getByText('Please enter your email address.')).toBeTruthy();
    expect(mockForgotPassword).not.toHaveBeenCalled();
    expect(screen.getAllByRole('link',{name:'Back to sign in'}).length).toBeGreaterThan(0);
  });

  it('shows a fixed countdown, allows editing the address, and permits resend after 30 seconds',async()=>{
    mockForgotPassword.mockResolvedValue({accepted:true,delivery:'email'});
    jest.useFakeTimers();
    renderRecovery();
    const email=screen.getByRole('textbox',{name:'Email address'}) as HTMLInputElement;
    fireEvent.change(email,{target:{value:'first@example.com'}});
    await act(async()=>{
      fireEvent.click(screen.getByRole('button',{name:'Send reset link'}));
      await Promise.resolve();
    });
    expect(screen.getByText('Request received!')).toBeTruthy();
    expect(screen.getByText('If eligible, a password reset link will be sent to first@example.com.')).toBeTruthy();
    expect(mockForgotPassword).toHaveBeenCalledWith('first@example.com');
    expect(screen.getByText('Check your email')).toBeTruthy();
    expect(screen.getAllByText(/a password reset link will be sent to/i).length).toBeGreaterThan(0);
    expect(screen.getAllByRole('link',{name:'Back to sign in'}).length).toBeGreaterThan(0);

    const cooldownButton=screen.getAllByRole('button',{name:'Resend in 30s'})[0];
    expect((cooldownButton as HTMLButtonElement).disabled).toBe(true);
    fireEvent.change(email,{target:{value:'corrected@example.com'}});
    expect(email.value).toBe('corrected@example.com');
    expect(screen.getByText('Email address updated')).toBeTruthy();
    expect((screen.getAllByRole('button',{name:'Resend in 30s'})[0] as HTMLButtonElement).disabled).toBe(true);

    act(()=>{jest.advanceTimersByTime(30_100);});
    const send=screen.getByRole('button',{name:'Send reset link'});
    expect((send as HTMLButtonElement).disabled).toBe(false);
    await act(async()=>{
      fireEvent.click(send);
      await Promise.resolve();
    });
    expect(mockForgotPassword).toHaveBeenCalledTimes(2);
    expect(mockForgotPassword).toHaveBeenLastCalledWith('corrected@example.com');
    expect(screen.getByText('Request received!')).toBeTruthy();
    expect(screen.getByText('If eligible, a password reset link will be sent to corrected@example.com.')).toBeTruthy();
    expect((screen.getAllByRole('button',{name:'Resend in 30s'})[0] as HTMLButtonElement).disabled).toBe(true);
  });

  it('keeps Back to sign in enabled throughout recovery',async()=>{
    mockForgotPassword.mockResolvedValue({accepted:true,delivery:'email'});
    renderRecovery();
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'parent@example.com'}});
    fireEvent.click(screen.getByRole('button',{name:'Send reset link'}));
    await waitFor(()=>expect(screen.getByText('Request received!')).toBeTruthy());
    fireEvent.click(screen.getAllByRole('link',{name:'Back to sign in'})[0]!);
    expect(screen.getByText('Sign in page')).toBeTruthy();
  });
});
