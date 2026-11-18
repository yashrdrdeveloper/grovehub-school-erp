import { useState } from 'react';
import { Grid2X2 } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Dialog } from '@Components';

import { workspaceNavigation } from '../registry';

export const AppLauncher = () => {
  const [open, setOpen] = useState(false);

  return (
    <>
      <button
        type="button"
        className="icon-button app-switcher-button"
        aria-label="Open application launcher"
        onClick={() => setOpen(true)}
      >
        <Grid2X2 size={19} />
      </button>
      <Dialog open={open} onOpenChange={setOpen} title="Your workspace" description="Jump to a GrooveHub module.">
        <div className="app-launcher-grid">
          {workspaceNavigation.map(({ label, path, icon: Icon }) => (
            <Link key={path} to={path} className="app-launcher-item" onClick={() => setOpen(false)}>
              <span><Icon size={22} strokeWidth={1.7} /></span>
              <strong>{label}</strong>
            </Link>
          ))}
        </div>
      </Dialog>
    </>
  );
};
