import { createRoot } from 'react-dom/client';
import { __ as translate } from '@Base/i18n';
import '@groovehub/styles/index.css';
import '@groovehub/components/styles/primitives.css';
import './styles.css';

Object.assign(window, { __: translate });
const root = document.getElementById('root');
if (!root) throw new Error('Missing application mount');
void import('../packages/shell/containers/ShellContainer').then(({ ShellContainer }) => {
  createRoot(root).render(<ShellContainer />);
}).catch(() => { root.textContent = 'Unable to load Internal Tool. Please refresh.'; });
