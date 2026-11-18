import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { ShellContainer } from '@Shell';

import './styles.css';

const element = document.getElementById('root');
if (!element) {
  throw new Error('Missing #root application mount point');
}

createRoot(element).render(
  <StrictMode>
    <ShellContainer />
  </StrictMode>,
);
