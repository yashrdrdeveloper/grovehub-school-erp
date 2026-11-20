import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { LoginV2Page } from './LoginV2Page';

const signIn=jest.fn<Promise<void>,[string,string]>();
const setPreference=jest.fn();

jest.mock('@Components',()=>({SchoolMonogram:()=> 'RDR'}));
jest.mock('./loginv2/SoffitGradient',()=>({__esModule:true,default:()=> <canvas className="lv2-gradient" aria-hidden="true"/>}));
jest.mock('../core/providers/AuthProvider',()=>({
  useAuth:()=>({signIn}),
}));
jest.mock('@Components/hooks/theme',()=>({
  useWorkspaceTheme:()=>({preference:'system',setPreference}),
}));

describe('RDR Login V2',()=>{
  beforeEach(()=>{
    signIn.mockReset();
    setPreference.mockReset();
  });

  it('renders the unified login layout and gradient without photographs or registration',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    expect(screen.getByRole('heading',{name:'Welcome back.'})).toBeTruthy();
    expect(document.querySelector('.rdr-login-v2__visual .lv2-gradient')).not.toBeNull();
    expect(document.querySelectorAll('.rdr-login-v2 img')).toHaveLength(0);
    expect(screen.getByRole('button',{name:/^sign in$/i})).toBeTruthy();
    expect(screen.queryByRole('button',{name:/continue with google/i})).toBeNull();
    expect(screen.queryByRole('button',{name:/create account/i})).toBeNull();
    expect(screen.getByRole('link',{name:'Forgot password?'}).getAttribute('href')).toBe('/forgot-password');
  });

  it('shows field validation and lets users toggle password visibility',async()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    fireEvent.click(screen.getByRole('button',{name:/^sign in$/i}));
    expect(screen.getByText('Please enter your email address.')).toBeTruthy();
    expect(screen.getByText('Please enter your password.')).toBeTruthy();
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'bad-email'}});
    fireEvent.click(screen.getByRole('button',{name:/^sign in$/i}));
    expect(screen.getByText('Enter a valid email address.')).toBeTruthy();
    const password=screen.getByLabelText('Password') as HTMLInputElement;
    expect(password.type).toBe('password');
    fireEvent.click(screen.getByRole('button',{name:'Show password'}));
    expect(password.type).toBe('text');
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'parent@example.com'}});
    fireEvent.change(password,{target:{value:'school-password'}});
    signIn.mockResolvedValueOnce();
    fireEvent.click(screen.getByRole('button',{name:/^sign in$/i}));
    await waitFor(()=>expect(signIn).toHaveBeenCalledWith('parent@example.com','school-password'));
  });

  it('keeps one feedback region, with a clear authentication error and the same CTA',async()=>{
    signIn.mockRejectedValueOnce(new Error('Invalid credentials'));
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'parent@example.com'}});
    fireEvent.change(screen.getByLabelText('Password'),{target:{value:'wrong-password'}});
    fireEvent.click(screen.getByRole('button',{name:/^sign in$/i}));
    expect(await screen.findByText('Incorrect email or password. Please try again.')).toBeTruthy();
    expect(screen.getAllByRole('alert')).toHaveLength(1);
    expect(document.querySelectorAll('.rdr-login-v2__feedback-slot--signin')).toHaveLength(1);
    expect(screen.getByRole('button',{name:/^sign in$/i})).toBeTruthy();
  });

  it('uses a local WebGL canvas for the hero, with no external media dependencies',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const visual=document.querySelector('.rdr-login-v2__visual');
    expect(visual).not.toBeNull();
    expect(visual?.querySelectorAll('img,picture,video,source')).toHaveLength(0);
    expect(visual?.querySelectorAll('canvas.lv2-gradient')).toHaveLength(1);
    expect(visual?.textContent).toContain('RDR School');
  });

  it('keeps system, light, and dark theme choices accessible',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const light=screen.getByRole('button',{name:'Light theme'});
    fireEvent.click(light);
    expect(setPreference).toHaveBeenCalledWith('light');
    expect(screen.getByRole('button',{name:'System theme'}).getAttribute('aria-pressed')).toBe('true');
  });
});
