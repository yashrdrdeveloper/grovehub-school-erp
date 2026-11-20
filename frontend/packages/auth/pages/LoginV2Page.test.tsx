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

  it('renders a separate school-photo sign-in layout without fake social login or registration',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    expect(screen.getByRole('heading',{name:'Welcome back.'})).toBeTruthy();
    expect(screen.getByRole('img',{name:/RDR School students together/i})).toBeTruthy();
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

  it('falls back to original RDR photo URL if local assets are missing',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const image=screen.getByRole('img',{name:/RDR School students together/i}) as HTMLImageElement;
    expect(image.getAttribute('src')).toBe('/school-images/slide5.jpg');
    fireEvent.error(image);
    expect(image.getAttribute('src')).toBe('https://rdrschool.com/images/slide5.jpg');
    fireEvent.error(image);
    expect(screen.queryByRole('img',{name:/RDR School students together/i})).toBeNull();
  });

  it('keeps system, light, and dark theme choices accessible',()=>{
    render(<MemoryRouter><LoginV2Page/></MemoryRouter>);
    const light=screen.getByRole('button',{name:'Light theme'});
    fireEvent.click(light);
    expect(setPreference).toHaveBeenCalledWith('light');
    expect(screen.getByRole('button',{name:'System theme'}).getAttribute('aria-pressed')).toBe('true');
  });
});
