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
    const forgotLinks=screen.getAllByRole('link',{name:'Forgot password?'});
    expect(forgotLinks).toHaveLength(2);
    expect(forgotLinks.every(link=>link.getAttribute('href')==='/forgot-password')).toBe(true);
    expect(forgotLinks[0].classList.contains('rdr-login-v2__forgot-desktop')).toBe(true);
    expect(forgotLinks[1].classList.contains('rdr-login-v2__forgot-mobile')).toBe(true);
  });

  it('uses minimal phone structure without extra mobile headings or decorative imagery',()=>{
    const {container}=render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const layout=container.querySelector('.rdr-login-v2');
    expect(layout?.classList.contains('rdr-login-v2--mobile-minimal')).toBe(true);
    expect(container.querySelector('.rdr-login-v2__mobile-signin-title')?.textContent).toBe('Sign in');
    expect(container.querySelector('.rdr-login-v2__mobile-header .school-monogram')).not.toBeNull();
    expect(container.querySelector('.rdr-login-v2__mobile-illustration')).toBeNull();
    expect(container.querySelector('.rdr-login-v2__mobile-security')).toBeNull();
    expect(container.querySelector('.rdr-login-v2__mobile-appearance')).not.toBeNull();
    expect(container.querySelector('svg.rdr-login-v2__mobile-school-art')).not.toBeNull();
    expect(container.querySelector('.rdr-login-v2__mobile-school-name')?.textContent).toBe('RDR School');
    expect(container.querySelector('.rdr-login-v2__mobile-description')?.textContent).toContain('workspace');
    expect(container.querySelector('.rdr-login-v2__form--signin > .rdr-login-v2__forgot-mobile')).not.toBeNull();
    expect(screen.getByRole('button',{name:'Light mobile theme'})).toBeTruthy();
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

  it('shows field errors only when present without changing the form structure',()=>{
    const {container}=render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const errors=()=>Array.from(container.querySelectorAll('.rdr-login-v2__field-error--reserved'));
    expect(errors()).toHaveLength(2);
    expect(errors().every(node=>node.classList.contains('is-empty'))).toBe(true);
    fireEvent.click(screen.getByRole('button',{name:/^sign in$/i}));
    expect(errors().every(node=>!node.classList.contains('is-empty'))).toBe(true);
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'student@example.edu'}});
    expect(errors()[0].classList.contains('is-empty')).toBe(true);
    expect(errors()[1].classList.contains('is-empty')).toBe(false);
    expect(container.querySelector('.rdr-login-v2__mobile-signin-title')?.textContent).toBe('Sign in');
    expect(container.querySelector('.rdr-login-v2__mobile-appearance')).not.toBeNull();
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
