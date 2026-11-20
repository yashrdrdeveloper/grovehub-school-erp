import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { LoginV2Page } from './LoginV2Page';

const signIn=jest.fn<Promise<void>,[string,string]>();
const setPreference=jest.fn();

jest.mock('@Components',()=>({SchoolMonogram:()=> 'RDR'}));
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

  it('renders gradient light-trail school artwork without photographs, social login or registration',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    expect(screen.getByRole('heading',{name:'Welcome back.'})).toBeTruthy();
    expect(document.querySelector('.rdr-login-v2__light-art')).not.toBeNull();
    expect(document.querySelector('.rdr-login-v2__book-art')).not.toBeNull();
    expect(document.querySelectorAll('.rdr-login-v2 img')).toHaveLength(0);
    expect(screen.getByRole('button',{name:/sign in to workspace/i})).toBeTruthy();
    expect(screen.queryByRole('button',{name:/continue with google/i})).toBeNull();
    expect(screen.queryByRole('button',{name:/create account/i})).toBeNull();
    expect(screen.getByRole('link',{name:'Forgot password?'}).getAttribute('href')).toBe('/forgot-password');
  });

  it('shows field validation and lets users toggle password visibility',async()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    fireEvent.click(screen.getByRole('button',{name:/sign in to workspace/i}));
    expect(screen.getByText('Please enter your email address.')).toBeTruthy();
    expect(screen.getByText('Please enter your password.')).toBeTruthy();
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'bad-email'}});
    fireEvent.click(screen.getByRole('button',{name:/sign in to workspace/i}));
    expect(screen.getByText('Enter a valid email address.')).toBeTruthy();
    const password=screen.getByLabelText('Password') as HTMLInputElement;
    expect(password.type).toBe('password');
    fireEvent.click(screen.getByRole('button',{name:'Show password'}));
    expect(password.type).toBe('text');
    fireEvent.change(screen.getByRole('textbox',{name:'Email address'}),{target:{value:'parent@example.com'}});
    fireEvent.change(password,{target:{value:'school-password'}});
    signIn.mockResolvedValueOnce();
    fireEvent.click(screen.getByRole('button',{name:/sign in to workspace/i}));
    await waitFor(()=>expect(signIn).toHaveBeenCalledWith('parent@example.com','school-password'));
  });

  it('uses only local CSS and SVG for the hero, with no external visual dependencies',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const visual=document.querySelector('.rdr-login-v2__visual');
    expect(visual).not.toBeNull();
    expect(visual?.querySelectorAll('img,picture,video,source')).toHaveLength(0);
    expect(visual?.querySelectorAll('svg.rdr-login-v2__light-art')).toHaveLength(1);
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
