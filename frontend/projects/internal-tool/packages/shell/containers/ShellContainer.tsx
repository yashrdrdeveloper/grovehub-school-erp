import { useEffect, useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate, NavLink } from '@Components/router';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { ShieldCheck, Menu, X, LogOut, FileCheck, Workflow, Users, Server, Code2, History, Sun, Moon } from 'lucide-react';
import { ApiError, errorText, request } from '../../business/api/client';
import type { Operator } from '../../business/api/types';
import { LoginPage } from '../../domain/auth/pages/login/LoginPage';
import { RegistryPage } from '../../domain/registry/pages/RegistryPage';
import { OperatorsPage } from '../../domain/operators/pages/OperatorsPage';
import { RequestsPage } from '../../domain/requests/pages/RequestsPage';
import { AuditPage } from '../../domain/audit/pages/AuditPage';

const routes=[{path:'/requests',label:__('My requests'),icon:Workflow,l2:false},{path:'/approvals',label:__('Approvals'),icon:FileCheck,l2:true},{path:'/api-templates',label:__('API Explorer'),icon:Code2,l2:false},{path:'/environments',label:__('Environments'),icon:Server,l2:false},{path:'/operators',label:__('Operators'),icon:Users,l2:true},{path:'/audit',label:__('Workflow history'),icon:History,l2:true}];
function Workspace({operator,onLogout}:{operator:Operator;onLogout:()=>void}){
 const [menu,setMenu]=useState(false);
 const [busy,setBusy]=useState(false);
 const [error,setError]=useState('');
 const [dark,setDark]=useState(()=>window.matchMedia('(prefers-color-scheme: dark)').matches);
 useEffect(()=>{const theme=(dark&&'dark')||'light';document.documentElement.dataset.theme=theme;document.documentElement.style.colorScheme=theme;},[dark]);
 async function logout(){setBusy(true);setError('');try{await request('/internal/protected/auth/logout','POST');onLogout();}catch(failure){setError(errorText(failure));}finally{setBusy(false);}}
 const isL2=operator.persona==='L2';
 return <BrowserRouter><div className="workspace" data-testid="internal-workspace">
  <aside className={'sidebar '+((menu&&'sidebar-open')||'')}><div className="brand"><ShieldCheck size={25}/><div><strong>GrooveHub</strong><span>{__('Internal Tool')}</span></div><Button className="mobile-only" variant={BUTTON_VARIANT.GHOST} aria-label={__('Close menu')} onClick={()=>setMenu(false)} dataTestId="mobile-menu-close"><X size={20}/></Button></div><nav aria-label={__('Workspace navigation')}>{routes.filter((r)=>!r.l2||isL2).map(({path,label,icon:Icon})=><NavLink key={path} to={path} onClick={()=>setMenu(false)} dataTestId={'nav-'+path.slice(1)}><Icon size={17}/><span>{label}</span></NavLink>)}</nav><div className="sidebar-bottom"><span className="status">{operator.persona}</span><span className="operator-email">{operator.email}</span></div></aside>
  <div className="workspace-main"><header className="topbar"><Button className="mobile-only" variant={BUTTON_VARIANT.GHOST} onClick={()=>setMenu(true)} aria-label={__('Open menu')} dataTestId="mobile-menu-open"><Menu size={22}/></Button><span className="topbar-title">{__('Operator workspace')}</span><div className="inline"><Button variant={BUTTON_VARIANT.GHOST} aria-label={__('Toggle theme')} onClick={()=>setDark((v)=>!v)} dataTestId="theme-toggle">{dark&&<Sun size={17}/>} {!dark&&<Moon size={17}/>}</Button><Button variant={BUTTON_VARIANT.OUTLINE} loading={busy} onClick={()=>void logout()} dataTestId="logout"><LogOut size={16}/><span className="desktop-only">{__('Sign out')}</span></Button></div></header>
  <main className="content">{error&&<p className="error" role="alert">{error}</p>}<Routes><Route path="/requests" element={<RequestsPage key="requests" inbox={false} operator={operator}/>}/><Route path="/approvals" element={<RequestsPage key="approvals" inbox operator={operator}/>}/><Route path="/api-templates" element={<RegistryPage key="templates" kind="api-templates" operator={operator}/>}/><Route path="/environments" element={<RegistryPage key="env" kind="environments" operator={operator}/>}/><Route path="/operators" element={<OperatorsPage operator={operator}/>}/><Route path="/audit" element={(isL2&&<AuditPage/>)||<Navigate to="/requests" replace/>}/><Route path="*" element={<Navigate to="/requests" replace/>}/></Routes></main></div>
 </div></BrowserRouter>;
}
export function ShellContainer(){
 const [operator,setOperator]=useState<Operator|null>(null);
 const [loading,setLoading]=useState(true);
 const [error,setError]=useState('');
 const [retry,setRetry]=useState(0);
 useEffect(()=>{
  const abort=new AbortController();
  setLoading(true);setError('');
  function expired(){setOperator(null);}
  window.addEventListener('internal-session-expired',expired);
  void request<Operator>('/internal/protected/auth/me','GET',undefined,abort.signal).then((v)=>{if(!abort.signal.aborted)setOperator(v);}).catch((failure)=>{if(!abort.signal.aborted&&!(failure instanceof ApiError&&failure.status===401))setError(errorText(failure));}).finally(()=>{if(!abort.signal.aborted)setLoading(false);});
  return ()=>{abort.abort();window.removeEventListener('internal-session-expired',expired);};
 },[retry]);
 if(loading)return <main className="center-state" role="status">{__('Checking operator session…')}</main>;
 if(error)return <main className="center-state"><p role="alert">{error}</p><Button onClick={()=>setRetry((v)=>v+1)} dataTestId="session-retry">{__('Retry')}</Button></main>;
 if(!operator)return <LoginPage onLogin={setOperator}/>;
 return <Workspace operator={operator} onLogout={()=>setOperator(null)}/>;
}
