import { access, readdir, readFile } from 'node:fs/promises';
import path from 'node:path';

const source = async pathname => readFile(path.join(process.cwd(), pathname), 'utf8');
const exists = async pathname => {
  try { await access(path.join(process.cwd(), pathname)); return true; }
  catch { return false; }
};
const violations = [];

const pageRoot = 'packages/domain/auth/pages';
for (const entry of await readdir(path.join(process.cwd(), pageRoot), { withFileTypes: true })) {
  if (entry.isFile() && /\.tsx?$/.test(entry.name)) {
    violations.push(`${pageRoot}/${entry.name}: page root cannot own shared components or route/session helpers`);
  }
}

const retired = [
  'packages/domain/auth/pages/AuthAmbientGraphic.tsx',
  'packages/domain/auth/pages/AuthFrame.tsx',
  'packages/domain/auth/pages/AuthMobileArtwork.tsx',
  'packages/domain/auth/pages/loginv2/LoginV2.tsx',
  'packages/domain/auth/pages/loginv2/SoffitGradient.tsx',
  'packages/shell/registry/landingPaths.registry.ts',
  'packages/shell/registry/navigation.registry.ts',
  'packages/domain/publicSite/pages/landing/config/content.ts',
  'packages/domain/publicSite/pages/landing/config/pageRegistry.tsx',
];
for (const pathname of retired) {
  if (await exists(pathname)) violations.push(`${pathname}: retired implementation must remain deleted`);
}

const [art, brand, renderer, routes, protectedLayout, workspaceTemplate, workspaceContainer, catalog, publicIndex, guide] = await Promise.all([
  source('packages/domain/auth/components/authAmbientGraphic/AuthAmbientGraphic.tsx'),
  source('packages/components/components/SchoolMonogram.tsx'),
  source('packages/components/components/SVGRenderer/SVGRenderer.tsx'),
  source('packages/shell/components/appRoutes/AppRoutes.tsx'),
  source('packages/shell/components/protectedLayout/ProtectedLayout.tsx'),
  source('packages/shell/templates/workspaceShell/WorkspaceShell.tsx'),
  source('packages/shell/containers/workspaceShellContainer/WorkspaceShellContainer.tsx'),
  source('packages/shell/constants/shell.catalog.ts'),
  source('packages/domain/publicSite/index.ts'),
  source('../docs/engineering/frontend/FRONTEND_CODE_GUIDELINES.md'),
]);

if (!art.includes('@Assets/auth') || !art.includes('SVGRenderer')) {
  violations.push('auth ambient art must use @Assets/auth through shared SVGRenderer');
}
if (!brand.includes('@Assets/brand') || !brand.includes('SVGRenderer')) {
  violations.push('RDR monogram must use its approved source vector through SVGRenderer');
}
if (!renderer.includes('data-testid') || !renderer.includes('aria-hidden') || !renderer.includes('aria-label')) {
  violations.push('SVGRenderer must expose native selectors and decorative/informative accessibility');
}
if (!routes.includes('PUBLIC_SITE_PATHS') || !publicIndex.includes('PUBLIC_SITE_PATHS')) {
  violations.push('public routes must be owned by publicSite, not duplicated in shell');
}
if (!routes.includes('<RequireSession/>') || !routes.includes('<ProtectedLayout/>')
    || !workspaceTemplate.includes('<Outlet/>')
    || !workspaceContainer.includes('<Shell ')
    || !protectedLayout.includes('<WorkspaceShellContainer ')) {
  violations.push('protected routes must keep one pathless protected layout with a persistent Shell and nested Outlet');
}
if (/PublicLayout|PublicFrame/.test(routes)) {
  violations.push('public school routes must not use a workspace public layout');
}
const workspace = await source('packages/shell/renderer/shell.routes.tsx');
const appRouteConstants = await source('packages/business/app/routes/constants/app.routes.ts');
const appPaths = new Map(
  [...appRouteConstants.matchAll(/([A-Z][A-Z_]+):\s*['"]([^'"]+)['"]/g)]
    .map(match => [match[1], match[2]]),
);
const childRoutes = workspace.split('export const workspaceRoutes = ')[1] || '';
const implemented = new Set(
  [...childRoutes.matchAll(/path:\s*(APP_ROUTE\.[A-Z_]+|['"][^'"]+['"])/g)]
    .map(match => {
      const identifier = match[1];
      if (identifier.startsWith('APP_ROUTE.')) return appPaths.get(identifier.slice('APP_ROUTE.'.length));
      return identifier.slice(1, -1);
    })
    .filter(Boolean),
);
const catalogPaths = [...catalog.matchAll(/path:\s*APP_ROUTE\.([A-Z_]+)/g)]
  .map(match => appPaths.get(match[1])).filter(Boolean);
for (const pathname of catalogPaths) {
  if (!implemented.has(pathname)) violations.push(`catalog path ${pathname} has no implemented protected route`);
}
if (await exists('../docs/engineering/frontend/FRONTEND_REFACTOR_MIGRATION.md')) {
  violations.push('migration rules must remain merged into the single frontend guideline');
}
if (!guide.includes('## 21. Assets') || !guide.includes('## 22. Unified migration inventory')) {
  violations.push('frontend code guideline must own assets/layout rules and migration inventory');
}
if (violations.length) {
  for (const violation of violations) process.stderr.write(`${violation}\n`);
  process.exitCode = 1;
} else {
  process.stdout.write('Artwork ownership, public route registry, persistent shell, catalog and unified guideline verified.\n');
}
