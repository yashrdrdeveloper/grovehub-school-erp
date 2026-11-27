# Frontend workspace

This directory owns the pnpm workspace and lockfile. No repository-root Node package is needed.

- School: projects/school, package groovehub-frontend, port 4173.
- Internal Tool: projects/internal-tool, package @groovehub/internal-tool, port 4174.
- Shared source packages: libs/base, libs/components, libs/styles.

From frontend/ install only a product and its dependencies with `pnpm --filter groovehub-frontend... install --frozen-lockfile` or `pnpm --filter @groovehub/internal-tool... install --frozen-lockfile`. Run the corresponding filter's build, test or dev script. Install all workspace dependencies with `pnpm install --frozen-lockfile`. Existing aliases refer only to shared source or app-owned code. Package publishing and version-pinned independent checkouts are not implemented yet; do not replace workspace dependencies with registry versions until the packages are tested and published.
