import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

import { ForgotPasswordPage } from './ForgotPasswordPage';

const mockForgotPassword = jest.fn<Promise<{accepted:boolean;delivery:string}>,[string]>();

jest.mock('@Business/auth',()=>({
  authenticationApi:{forgotPassword:(email:string)=>mockForgotPassword(email)},
}));
jest.mock('./LoginLayout',()=>({
  LoginLayout:({children,title}:{children:React.ReactNode;title:string})=><div><h1>{title}</h1>{children}</div>,
}));

const renderRecovery=()=>render(<MemoryRouter initialEntries={['/forgot-password']}>
  <Routes>
    <Route path="/forgot-password" element={<ForgotPasswordPage/>}/>
    <Route path="/login" element={<p>Sign in page</p>}/>
  </Routes>
</MemoryRouter>);

describe('forgot password recovery',()=>{
  beforeEach(()=>{
    mockForgotPassword.mockReset();
    jest.useRealTimers();
  });
  afterEach(()=>jest.useRealTimers());

  it('shows guidance, validates email and preserves its action',()=>{
    renderRecovery();
    expect(screen.getByRole('heading',{name:'Forgot password?'})).toBeTruthy();
    expect(screen.getByText('A secure link, just for you')).toBeTruthy();
    fireEvent.click(screen.getByRole('button',{name:'Send reset link'}));
    expect(screen.getByText('Please enter your email address.')).toBeTruthy();
    expect(mockForgotPassword).not.toHaveBeenCalled();
    expect(screen.getByRole('link',{name:'Back to sign in'})).toBeTruthy();
  });

  it('keeps Back to sign in active and makes resend available after 30 seconds',async()=>{
    mockForgotPassword.mockResolvedValue({accepted:true,delivery:'email'});
    renderRecovery();
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'parent@example.com'}});
    fireEvent.click(screen.getByRole('button',{name:'Send reset link'}));
    expect(await screen.findByText('You’re one step closer!')).toBeTruthy();
    expect(mockForgotPassword).toHaveBeenCalledTimes(1);
    expect(screen.getByRole('button',{name:'Back to sign in'})).toBeTruthy();
    expect((screen.getByRole('button',{name:/Resend reset link in/}) as HTMLButtonElement).disabled).toBe(true);
    jest.useFakeTimers();
    act(()=>{jest.advanceTimersByTime(30_100);});
    const resend=screen.getByRole('button',{name:'Resend reset link'});
    expect((resend as HTMLButtonElement).disabled).toBe(false);
    fireEvent.click(resend);
    await act(async()=>{await Promise.resolve();});
    expect(mockForgotPassword).toHaveBeenCalledTimes(2);
  });

  it('navigates back to sign in after a successful request',async()=>{
    mockForgotPassword.mockResolvedValue({accepted:true,delivery:'email'});
    renderRecovery();
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'parent@example.com'}});
    fireEvent.click(screen.getByRole('button',{name:'Send reset link'}));
    await waitFor(()=>expect(screen.getByRole('button',{name:'Back to sign in'})).toBeTruthy());
    fireEvent.click(screen.getByRole('button',{name:'Back to sign in'}));
    expect(screen.getByText('Sign in page')).toBeTruthy();
  });
});
