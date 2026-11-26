import { existsSync } from 'node:fs';
import path from 'node:path';

/** Keep logical package owners stable when a package moves outside the app. */
export const createSourceLayout = (appRoot, packageNames = ['base']) => {
  const localRoot = path.join(appRoot, 'packages');
  const sharedRoot = path.resolve(appRoot, '../libs/frontend');
  const sharedPackages = packageNames.flatMap((name) => {
    const local = path.join(localRoot, name);
    const shared = path.join(sharedRoot, name);
    const localExists = existsSync(local);
    const sharedExists = existsSync(shared);
    if (localExists && sharedExists) throw new Error(`Duplicate package owner: ${name}`);
    if (sharedExists) return [[name, shared]];
    if (!localExists) throw new Error(`Missing package owner: ${name}`);
    return [];
  });
  const sources = new Map(sharedPackages);
  const resolvePackage = (relative) => {
    const [name, ...rest] = relative.split('/');
    const directory = sources.get(name);
    if (directory) return path.join(directory, ...rest);
    return path.join(localRoot, relative);
  };
  const getPackageRelative = (filename) => {
    const shared = path.relative(sharedRoot, filename);
    if (!shared.startsWith('..') && !path.isAbsolute(shared)) return shared;
    return path.relative(localRoot, filename);
  };
  const getFrontendRelative = (filename) => {
    const shared = path.relative(sharedRoot, filename);
    if (!shared.startsWith('..') && !path.isAbsolute(shared)) {
      return path.join('packages', shared).replaceAll(path.sep, '/');
    }
    return path.relative(appRoot, filename).replaceAll(path.sep, '/');
  };
  const resolveFrontend = (relative) => {
    if (relative.startsWith('packages/')) return resolvePackage(relative.slice('packages/'.length));
    return path.join(appRoot, relative);
  };
  return {
    sharedPackageSources: sharedPackages,
    packageSourceRoots: [localRoot, ...sources.values()],
    resolvePackageSource: resolvePackage,
    packageRelativePath: getPackageRelative,
    frontendRelativePath: getFrontendRelative,
    resolveFrontendSource: resolveFrontend,
  };
};

export const {
  sharedPackageSources,
  packageSourceRoots,
  resolvePackageSource,
  packageRelativePath,
  frontendRelativePath,
  resolveFrontendSource,
} = createSourceLayout(process.cwd());
