import { ChevronDown, Menu } from 'lucide-react';
import type { TopbarComponent } from './types/topbar';

export const Topbar: TopbarComponent = ({ accountLabel, schoolLabel, model }) => (
  <header className="topbar">
    <div className="topbar__left">
      <button type="button" className="icon-button topbar__menu" onClick={model.openSidebar}
        aria-label="Open menu"><Menu size={20}/></button>
      <span className="topbar__breadcrumb">{schoolLabel}</span>
      <span className="topbar__separator">/</span><span className="topbar__muted">Workspace</span>
    </div>
    <div className="topbar__right">
      <button type="button" className="icon-button topbar__quick" onClick={model.openAccountPreferences}
        aria-label="Open account preferences">
        <span className="profile-avatar profile-avatar--small">{accountLabel.charAt(0).toUpperCase()}</span>
        <ChevronDown size={14}/>
      </button>
    </div>
  </header>
);
