import { createRoot } from 'react-dom/client';
import { installShellI18n } from '../packages/shell/containers/shellContainer/helpers/shellContainer.installI18n';
import '@groovehub/styles/index.css';
import '@groovehub/components/styles/primitives.css';
import './styles.css';

installShellI18n();
const root = document.getElementById('root');
if (!root) throw new Error('Missing application mount');
void import('../packages/shell/containers/ShellContainer').then(({ ShellContainer }) => {
  createRoot(root).render(<ShellContainer />);
}).catch(() => { root.textContent = 'Unable to load Internal Tool. Please refresh.'; });
