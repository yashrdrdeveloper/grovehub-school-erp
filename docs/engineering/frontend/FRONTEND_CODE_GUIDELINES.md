# Frontend Code Guidelines

> **PackageManager** : use pnpm over npm always and react latest stable version

> **Status:** Mandatory Phase 0 engineering standard — Revision 3

> **Applies to:** `frontend/`, `packages/base`, `packages/business`, `packages/components`, `packages/domain`, `packages/shell`, every domain, page, template, component, hook, HOC, helper, util, constant, service, reader, mock, type, story and test

> **Primary stack:** React + TypeScript + Vite

> **Testing:** Jest + React Testing Library

> **Core principles:** **SRP**, **OCP**, **DRY**, composition, configuration-driven rendering, strong typing, explicit public boundaries, testability

---

## 1. Objective

The frontend must be designed as a collection of small, composable, reusable, strongly typed units.

The architecture must remain predictable even when the product grows to many domains and hundreds of screens.

Mandatory principles:

1. **SRP — Single Responsibility Principle**

2. **OCP — Open/Closed Principle**

3. **DRY — Don't Repeat Yourself**

4. Composition over inheritance

5. Configuration/registry/strategy over repeated conditional markup

6. Pure functions for transformations

7. Custom hooks for reusable stateful behavior

8. Explicit public APIs through `index.ts`

9. No deep imports into another unit's internals

10. Reuse the nearest authoritative type instead of redefining primitives

11. Presentational components contain no business orchestration

12. Containers/hooks/services own side effects and orchestration

13. Repeated UI patterns become platform-level primitives

14. `switch` is restricted

15. Avoid `if / else if` chains when a strategy/config map can express the behavior

16. Naming must explain intent without opening the file

17. Public contracts should remain stable once consumed

18. Architecture rules should be enforced automatically where practical

---

## 2. Top-Level Package Architecture

The frontend package hierarchy is mandatory. It separates pure technical primitives, shared business logic, shared UI primitives, isolated product domains, and the application-composition shell.

```text
frontend/
├── app/
│   ├── bootstrap/
│   ├── main.tsx
│   └── index.ts
│
├── packages/
│   ├── base/
│   │   ├── constants/
│   │   ├── helpers/
│   │   ├── utils/
│   │   ├── readers/
│   │   ├── types/
│   │   ├── RUNBOOK.md
│   │   └── index.ts
│   │
│   ├── business/
│   │   ├── services/
│   │   ├── helpers/
│   │   ├── utils/
│   │   ├── readers/
│   │   ├── constants/
│   │   ├── types/
│   │   ├── RUNBOOK.md
│   │   └── index.ts
│   │
│   ├── components/
│   │   ├── components/
│   │   ├── templates/
│   │   ├── builders/
│   │   ├── renderers/
│   │   ├── hooks/
│   │   ├── hocs/
│   │   ├── layouts/
│   │   ├── types/
│   │   ├── RUNBOOK.md
│   │   └── index.ts
│   │
│   ├── domain/
│   │   ├── catalog/
│   │   ├── delivery/
│   │   ├── testing/
│   │   ├── workbench/
│   │   └── ...
│   │
│   └── shell/
│       ├── Shell.tsx
│       ├── ShellContainer.tsx
│       ├── providers/
│       ├── state/
│       ├── registry/
│       ├── layout/
│       ├── skeleton/
│       ├── errors/
│       ├── types/
│       ├── RUNBOOK.md
│       └── index.ts
│
├── assets/
├── config/
├── test/
│   ├── fixtures/
│   ├── mocks/
│   └── setup/
├── vite.config.ts
├── tsconfig.json
├── eslint.config.*
├── jest.config.*
└── package.json
```

### Package responsibilities

| Package | Responsibility | May contain business logic? | May contain UI? | May call backend APIs? |
| --- | --- | ---: | ---: | ---: |
| `packages/base` | Pure, generic technical primitives | No | No | No |
| `packages/business` | Shared product/business contracts and logic used by multiple domains | Yes | No | Yes |
| `packages/components` | Shared platform UI, hooks, HOCs, templates, builders and renderers | No domain-specific business logic | Yes | No direct API calls |
| `packages/domain/<name>` | Domain-specific screens, state, orchestration and domain behavior | Yes | Yes | Yes, through domain services |
| `packages/shell` | Application composition, shared provider initialization, domain registration/binding, global layout and startup skeleton | **No** | Shell/chrome only | **No** |

### `frontend/app`

`frontend/app` is the executable browser entry point.

It should remain intentionally thin.

Its responsibility is to:

- perform minimal browser/bootstrap initialization
- resolve root DOM mounting
- load environment/bootstrap configuration required before React starts
- mount the application shell

Conceptually:

```tsx
import { createRoot } from 'react-dom/client';
import { PerihellionShell } from '@shell';

createRoot(document.getElementById('root')!).render(
  <PerihellionShell />,
);
```

`app` must not become the place where every domain route, reducer, provider, navigation item, or feature is manually wired.

That belongs in `packages/shell`.

### `packages/base`

`base` is the dependency floor. Examples:

```text
getDateTime
EMPTY_OBJECT
EMPTY_ARRAY
NOOP
safeParseNumber
createPropertyReader
```

`base` must remain pure and domain-independent. It must not import `business`, `components`, `domain`, `shell`, or `app`.

### `packages/business`

`business` contains reusable business logic shared by multiple domains. Examples:

```text
business/services/application/
business/services/environment/
business/utils/application/
business/helpers/permissions/
business/readers/application/
```

Do not promote logic into `business` merely because two files look similar. Promote it when the same business concept is intentionally shared and should evolve together.

### `packages/components`

`components` is the shared platform UI package. It owns reusable components such as `FormBuilder`, `TableManager`, `ModalManager`, `StatusBadge`, common hooks/HOCs, layouts, renderers and templates.

Shared components must not import a product domain or `packages/shell`.

### `packages/domain`

Every actual product domain lives under `packages/domain/<domain-name>`. Domains own pages, domain state, domain services and domain-specific UI.

A domain may consume:

```text
base
business
components
another domain's root public API when explicitly required
```

A domain must never import `packages/shell`.

### `packages/shell`

`shell` is the **application-composition package**.

It binds the frontend together without owning product business logic.

It is responsible for:

- the top-level `ShellContainer`
- the top-level presentational shell/chrome
- initialization of shared runtime providers
- Redux store/provider construction
- query-client provider construction
- router provider construction
- theme/design-system provider construction
- global error boundary/fallback wiring
- top-level telemetry provider wiring
- assembling domain registrations
- composing domain routes
- composing domain navigation
- composing domain reducers/middleware
- application-level layout
- startup/loading skeleton
- domain mount/unmount registration where supported

It must **not**:

- call domain/backend APIs directly
- contain product use cases
- contain authorization decisions
- contain domain reducers
- contain domain service logic
- contain domain entity transformations
- deep-import any domain internal path

The shell consumes domains only through their root public APIs.

Example:

```ts
import { catalogDomain } from '@domain/catalog';
import { deliveryDomain } from '@domain/delivery';
```

Never:

```ts
import { catalogReducer }
  from '@domain/catalog/core/state/catalog.slice';
```

The shell is the composition root of the frontend, not a new business domain.

---

## 3. Domain Package Structure

Every domain follows the same predictable unit structure.

```text
frontend/packages/domain/<domain-name>/
├── core/
│   ├── routes/
│   ├── state/
│   ├── providers/
│   ├── guards/
│   ├── registry/
│   ├── constants/
│   ├── types/
│   └── index.ts
│
├── pages/
├── templates/
├── components/
├── constants/
├── helpers/
├── utils/
├── hooks/
├── hocs/
├── types/
│
├── services/
│   ├── <entity-or-use-case>/
│   └── index.ts
│
├── mock/
│   ├── builders/
│   ├── fixtures/
│   ├── handlers/
│   ├── services/
│   ├── state/
│   └── index.ts
│
├── RUNBOOK.md
└── index.ts
```

Examples:

```text
frontend/packages/domain/catalog/
frontend/packages/domain/delivery/
frontend/packages/domain/testing/
frontend/packages/domain/workbench/
frontend/packages/domain/migration/
```

A domain's root `index.ts` is its runtime public API. The `mock/index.ts` file is a separate test/development public API.

Do not create services under individual components. API services belong at the domain level or, when intentionally shared across domains, in `packages/business`.

Do not create empty folders only to satisfy the template. Create the folder when the domain owns that responsibility.

---

## 4. Domain `core/`

`core/` owns domain-level wiring.

```text
domain/core/
├── routes/
│   ├── domain.routes.ts
│   ├── types/
│   └── index.ts
│
├── state/
│   ├── domain.slice.ts
│   ├── domain.actions.ts
│   ├── domain.selectors.ts
│   ├── domain.middleware.ts
│   ├── types/
│   └── index.ts
│
├── providers/
├── guards/
├── constants/
├── types/
└── index.ts
```

Responsibilities may include:

\- domain routes

\- domain Redux state registration

\- domain actions/selectors

\- domain context/providers

\- domain feature guards

\- domain bootstrap

\- domain configuration

\- domain permission mappings

\- public domain contracts

A domain should be mountable from its public API.

```ts
import { deliveryDomain } from '@/packages/delivery';
```

---

## 5. Public API Rule — `index.ts`

Every logical unit exposes its supported public API through `index.ts`.

External consumers must import through the nearest public boundary.

### Correct

```ts
import {
  DeploymentStatus,
  type DeploymentStatusProps,
} from '@/packages/delivery/components/DeploymentStatus';
```

### Incorrect

```ts
import { DeploymentStatus }
  from '@/packages/delivery/components/DeploymentStatus/DeploymentStatus';
import { mapStatus }
  from '@/packages/delivery/components/DeploymentStatus/helpers/deploymentStatus.mapper';
```

Internal files are private unless intentionally exported.

---

## 6. `index.ts` Is a Contract, Not a Dumping Ground

Good:

```ts
export { DeploymentStatus } from './DeploymentStatus';
export { DeploymentStatusContainer } from './DeploymentStatusContainer';
export type {
  DeploymentStatusProps,
  DeploymentStatusViewModel,
} from './types';
```

Avoid:

```ts
export * from './helpers';
export * from './constants';
export * from './internal';
```

unless every export is intentionally public.

Use barrels mainly at package/unit boundaries. Inside a unit, direct sibling imports are allowed when they prevent circular dependencies.

---

## 7. Dependency Direction and Domain Isolation

The dependency graph is intentionally one-directional.

```text
                               app
                                ↓
                              shell
                    ┌───────────┼───────────┐
                    ↓           ↓           ↓
              domain/<name>  components    base
               ↙      ↓
          business   components
              ↓         ↓
             base      base
```

A more explicit form:

```text
base
  ↑
  ├──────── business
  ├──────── components
  │
  └──────── domain
             ↑
             │
           shell
             ↑
             │
            app
```

The diagram is about dependency direction, not runtime ownership. `shell` may also directly consume `components` and `base`.

Mandatory boundaries:

```text
base        -> third-party libraries only

business    -> base

components  -> base

domain      -> base
            + business
            + components
            + another domain's ROOT PUBLIC API only when explicitly required

shell       -> base
            + components
            + domain ROOT PUBLIC APIs

app         -> shell
            + minimal bootstrap/config dependencies only
```

Forbidden dependency directions:

```text
base        ✗ business
base        ✗ components
base        ✗ domain
base        ✗ shell

business    ✗ domain
business    ✗ shell

components  ✗ domain
components  ✗ shell

domain      ✗ shell
domain      ✗ app

shell       ✗ domain internal paths

app         ✗ domain internal paths
app         ✗ direct domain composition that belongs in shell
```

`business` and `components` must not import a product domain.

A domain must never deep-import another domain.

### Allowed cross-domain import

```ts
import {
  ApplicationSelector,
  type ApplicationSummary,
} from '@domain/catalog';
```

### Forbidden

```ts
import { selectCatalogInternalState }
  from '@domain/catalog/core/state/catalog.internal.selectors';

import { ApplicationRow }
  from '@domain/catalog/components/ApplicationTable/components/ApplicationRow';
```

If Catalog wants another domain to use `ApplicationSelector`, Catalog explicitly exports it from its root `index.ts`.

The exporting domain owns the contract. The consuming domain does not reach into implementation folders.

### Domain import rule

Cross-domain imports must look conceptually like:

```ts
import { something } from '@domain/<domain-name>';
```

and never like:

```ts
import { something } from '@domain/<domain-name>/<internal-path>';
```

The shell follows the same rule.

CI/ESLint must enforce:

- no cross-domain deep imports
- no domain-to-shell imports
- no components-to-shell imports
- no business-to-shell imports
- no app-level deep imports into domains
- shell may import only domain root public APIs

---

## 8. Application Shell Package

`packages/shell` is the frontend **composition root**.

It exists so global application wiring does not leak into:

```text
frontend/app
individual domains
shared components
business services
base utilities
```

The shell binds domain capabilities together while remaining business-logic-free.

### Recommended structure

```text
frontend/packages/shell/
├── Shell.tsx
├── ShellContainer.tsx
│
├── providers/
│   ├── ReduxProvider/
│   ├── QueryProvider/
│   ├── RouterProvider/
│   ├── ThemeProvider/
│   ├── TelemetryProvider/
│   ├── ErrorBoundaryProvider/
│   └── index.ts
│
├── state/
│   ├── createShellStore.ts
│   ├── reducerRegistry.ts
│   ├── middlewareRegistry.ts
│   ├── types/
│   └── index.ts
│
├── registry/
│   ├── domains.registry.ts
│   ├── routes.registry.ts
│   ├── navigation.registry.ts
│   ├── commands.registry.ts
│   ├── types/
│   └── index.ts
│
├── layout/
│   ├── ShellLayout/
│   └── index.ts
│
├── skeleton/
│   ├── ShellSkeleton.tsx
│   ├── ShellSkeleton.stories.tsx
│   ├── types/
│   └── index.ts
│
├── errors/
├── constants/
├── helpers/
├── types/
├── RUNBOOK.md
└── index.ts
```

Do not create empty folders only for symmetry.

### `Shell.tsx`

`Shell.tsx` is the presentational application frame.

It may render:

```text
application chrome
global header
navigation slots
main content outlet
global loading/skeleton state
global error/fallback UI
```

It must not:

```text
fetch domain data
dispatch domain business actions
interpret backend entities
make permission decisions
contain domain-specific conditions
```

### `ShellContainer.tsx`

`ShellContainer.tsx` owns **composition/orchestration only**.

Typical responsibilities:

```text
initialize global providers
construct provider nesting
construct global Redux runtime
consume domain registrations
create route/navigation registries
select top-level loading/error fallback
pass prepared shell configuration to Shell.tsx
```

It is allowed to orchestrate infrastructure.

It is not allowed to implement a product use case.

### Shared provider initialization

Global providers should be initialized once in the shell.

Examples:

```text
Redux Provider
TanStack Query / QueryClient Provider
Router Provider
Theme / Design System Provider
Global Error Boundary
Telemetry Context
Internationalization Provider when introduced
other truly application-wide technical providers
```

Do not initialize the same provider independently inside every domain.

A domain may still own a **domain-specific provider** inside its own `core/providers/` when the provider is not truly global.

### Redux ownership

The shell owns **Redux runtime composition**, not domain state semantics.

The separation is:

```text
Domain
  owns:
    slice
    reducers
    actions
    selectors
    domain middleware
    state types
    registration metadata

Shell
  owns:
    store creation
    root reducer composition
    middleware composition
    Redux <Provider>
```

A domain should expose only the state registration it intentionally wants the shell to mount.

Conceptually:

```ts
export const catalogDomain = {
  id: DOMAIN_ID.CATALOG,
  routes: catalogRoutes,
  navigation: catalogNavigation,
  state: {
    key: CATALOG_STATE_KEY,
    reducer: catalogReducer,
  },
} satisfies FrontendDomainRegistration;
```

The shell consumes:

```ts
import { catalogDomain } from '@domain/catalog';
import { deliveryDomain } from '@domain/delivery';

export const DOMAIN_REGISTRY = Object.freeze([
  catalogDomain,
  deliveryDomain,
]) satisfies readonly FrontendDomainRegistration[];
```

The shell must never know:

```text
how catalog reducer calculates state
how delivery permissions work
how testing fetches data
how workbench maps backend responses
```

### Registration contract placement

A domain must not import `packages/shell` merely to implement the shell's registration type.

The shared `FrontendDomainRegistration` contract must live at a dependency-safe lower level that both domains and shell can consume.

Preferred placement:

```text
packages/components/types/registration/
```

because it is a frontend composition contract shared between product domains and the application shell.

If the registration contract is later split into pure technical state contracts and visual contracts, pure technical contracts may move to `packages/base/types/registration/`.

The important rule is:

> **Domain registration types cannot create a `domain -> shell` dependency.**

### Domain registration

Every mountable domain should export one explicit registration object from its root public API.

Example:

```ts
export {
  catalogDomain,
} from './core/registry';
```

Consumers:

```ts
import { catalogDomain } from '@domain/catalog';
```

The shell binds registrations.

No other package should need to manually know every domain's internal route/reducer/provider files.

### Routes

Domains own their routes.

The shell aggregates them.

```text
Domain route definitions
        ↓
Domain public registration
        ↓
Shell route registry
        ↓
Router Provider
```

`frontend/app` does not manually concatenate domain routes.

### Navigation

Domains expose navigation contributions as configuration.

The shell aggregates and passes them to shared navigation components.

The shell must not contain code such as:

```ts
if (domain === 'catalog') {
  ...
} else if (domain === 'delivery') {
  ...
}
```

Use registration/configuration.

### Commands and extension points

Where domains contribute:

```text
command palette commands
navigation
route definitions
dashboard registrations
global extension slots
```

they should expose typed configuration through their public registration.

The shell aggregates it.

### Shell skeleton

The shell owns the **application startup skeleton**, not domain-specific loading UI.

Examples:

```text
initial app bootstrap
provider initialization
session/bootstrap context resolution
route registry startup
lazy domain bundle loading
```

`ShellSkeleton` should remain generic.

Domain-specific skeletons stay inside the domain.

Because the shell skeleton is public visual infrastructure, it should have Storybook coverage.

### Shell and business logic

Forbidden inside `packages/shell`:

```text
createApplication()
approveRelease()
calculateReadiness()
mapDeploymentResponse()
fetchCatalogApplications()
authorization business decisions
migration calculations
provider-specific product behavior
```

Those belong to domains/business services.

### Shell and API calls

The shell must not become a "global API service."

If application startup needs business data, the owning domain should expose the required capability/hook/registration through its public contract.

Only truly technical bootstrap/configuration required before domains exist may be handled at app/shell bootstrap level.

### Shell public API

Normal consumers should use:

```ts
import {
  PerihellionShell,
} from '@shell';
```

`frontend/app` is the primary runtime consumer.

Domains must not consume the shell.

### Shell testing

Test shell behavior around:

```text
provider composition
domain registration
route aggregation
navigation aggregation
reducer/middleware registration
startup skeleton
global error fallback
```

Do not duplicate domain business tests inside shell tests.

### Shell runbook

`packages/shell/RUNBOOK.md` must document:

```text
shell responsibility
provider nesting/order
how to register a domain
how routes are aggregated
how Redux reducers/middleware are registered
how startup skeleton works
how global errors are handled
what is explicitly forbidden from shell
how to test shell composition
```

---

## 9. Pages Live at Domain Level

```text
domain/pages/
├── DeliveryHomePage/
├── ReleaseDetailsPage/
├── PipelineRunsPage/
└── index.ts
```

Recommended page unit:

```text
ReleaseDetailsPage/
├── ReleaseDetailsPage.tsx
├── ReleaseDetailsPageContainer.tsx
├── components/
├── templates/
├── hooks/
├── helpers/
├── constants/
├── types/
└── index.ts
```

Pages compose screens. They must not become giant business-logic files.

---

## 10. Presentational / Container Pattern

Stateful feature components and pages should use a presentational/container split.

### Presentational component

```tsx
export const ReleaseSummary: React.FC<ReleaseSummaryProps> = (props) =>{
  {
    title,
    status,
    owner,
    isLoading,
    onRefresh,
 } = props;
 return (
    <SummaryCard
      title={title}
      status={status}
      owner={owner}
      isLoading={isLoading}
      onRefresh={onRefresh}
    />
  );
}
```

Presentational components:

\- receive prepared props

\- render UI

\- emit callbacks

\- do not call APIs

\- do not access Redux directly

\- do not know provider/vendor implementations

\- do not build complex request payloads

\- do not perform large business transformations

\- do not own orchestration side effects

### Container

```tsx
export const ReleaseSummaryContainer = () => {
  const {
    viewModel,
    isLoading,
    refresh,
  } = useReleaseSummary();
  return (
    <ReleaseSummary
      {...viewModel}
      isLoading={isLoading}
      onRefresh={refresh}
    />
  );
};
```

The container owns orchestration.

Every page and stateful/reusable feature should use this model. Tiny pure visual atoms do not need an empty container solely for ceremony.

---

## 11. Component Unit Structure

A full component unit may contain presentation, orchestration, local helpers, hooks, HOCs and nested UI units, but **not reusable backend API services**.

```text
ComponentName/
├── ComponentName.tsx
├── ComponentNameContainer.tsx
│
├── components/
├── templates/
│
├── constants/
│   ├── componentName.useCase.constant.ts
│   ├── types/
│   │   └── componentName.useCase.constant.types.ts
│   └── index.ts
│
├── helpers/
│   ├── componentName.useCase.mapper.ts
│   ├── componentName.useCase.selector.ts
│   ├── types/
│   ├── __tests__/
│   └── index.ts
│
├── utils/
│   ├── componentName.useCase.util.ts
│   ├── types/
│   ├── __tests__/
│   └── index.ts
│
├── hooks/
│   └── useSomething/
│       ├── useSomething.ts
│       ├── constants/
│       ├── helpers/
│       ├── types/
│       ├── __tests__/
│       └── index.ts
│
├── hocs/
│   └── withSomething/
│       ├── withSomething.tsx
│       ├── constants/
│       ├── helpers/
│       ├── types/
│       ├── __tests__/
│       └── index.ts
│
├── types/
│   ├── ComponentName.props.ts
│   ├── ComponentName.state.ts
│   ├── ComponentName.viewModel.ts
│   ├── ComponentName.useCase.ts
│   └── index.ts
└── index.ts
```

A component may consume a domain/business service through an injected hook or container, but the component folder does not own the service contract.

This avoids coupling reusable UI to a single API endpoint and prevents duplicate service definitions when the same backend entity is used by several components.

---

## 12. Nested Components

Components may own nested components.

```text
PipelineTable/
├── PipelineTable.tsx
├── components/
│   ├── PipelineRow/
│   │   ├── PipelineRow\.tsx
│   │   ├── types/
│   │   └── index.ts
│   └── index.ts
└── index.ts
```

Nested units remain private unless exported through the parent unit's public API.

---

## 13. Templates

Templates are reusable composition/layout units.

```text
templates/
└── ApplicationDetailsTemplate/
    ├── ApplicationDetailsTemplate.tsx
    ├── components/
    ├── constants/
    ├── helpers/
    ├── hooks/
    ├── types/
    └── index.ts
```

Templates primarily handle:

\- layout

\- composition

\- slots

\- panel placement

\- responsive structure

\- repeated screen skeletons

Business orchestration belongs in containers/hooks.

---

## 14. Hooks

Hooks own reusable stateful behavior.

Do not pollute components with fetch/state/transformation logic.

### Avoid

```tsx
export const SomethingContainer = () => {
  const [data, setData] = useState([]);
  const [isLoading, setLoading] = useState(false);
  useEffect(() => {
    setLoading(true);
    fetchSomething()
      .then((response) => setData(response.items))
      .finally(() => setLoading(false));
  }, []);
  const grouped = data.reduce(...);
  const filtered = grouped.filter(...);
  return <Something data={filtered} />;
};
```

### Prefer

```ts
const {
  data,
  isLoading,
  refresh,
} = useFetchSomething();
const viewModel = useSomethingViewModel(data);
```

Then:

```tsx
<Something
  {...viewModel}
  isLoading={isLoading}
  onRefresh={refresh}
/>
```

---

## 15. Hook Structure

```text
hooks/
└── useFetchSomething/
    ├── useFetchSomething.ts
    ├── constants/
    ├── helpers/
    ├── types/
    ├── __tests__/
    └── index.ts
```

Hooks should have one coherent responsibility.

Compose focused hooks instead of creating one giant hook.

```ts
const query = useApplicationQuery(applicationId);
const permissions = useApplicationPermissions(applicationId);
const viewModel = useApplicationViewModel(query.data);
const actions = useApplicationActions(applicationId);
```

---

## 16. HOCs

Use HOCs only for wrapper-style composition.

Typical uses:

\- permission wrapper

\- feature gating

\- error boundary

\- telemetry wrapper

\- context/provider wrapping

```text
hocs/
└── withPermission/
    ├── withPermission.tsx
    ├── constants/
    ├── helpers/
    ├── types/
    ├── __tests__/
    └── index.ts
```

Prefer hooks when a wrapper is not needed.

---

## 17. Services, Backend Entities, and Response Readers

API services live only at one of these levels:

```text
packages/domain/<domain>/services/
packages/business/services/
```

They must not live inside component folders.

Use a domain service when the API/entity is owned by one domain. Promote it to `packages/business/services/<entity>` only when the same business contract is intentionally shared across multiple domains.

Recommended service unit:

```text
services/<entity-or-use-case>/
├── fetchSomething.service.ts
├── updateSomething.service.ts
├── index.ts
│
├── types/
│   ├── request/
│   │   ├── fetchSomething.request.ts
│   │   └── index.ts
│   ├── response/
│   │   ├── something.response.ts
│   │   └── index.ts
│   ├── fetchSomething.service.types.ts
│   └── index.ts
│
├── response/
│   └── readers/
│       ├── somethingResponse.reader.ts
│       ├── types/
│       ├── __tests__/
│       └── index.ts
│
└── __tests__/
```

Implementation:

```ts
export const fetchSomething: FetchSomething = async (request) => {
  return api.post<SomethingResponse>(
    '/something',
    request,
  );
};
```

Callable type:

```ts
export type FetchSomething = (
  request: FetchSomethingRequest,
) => Promise<SomethingResponse>;
```

## Response readers

Readers are the approved property-access layer for backend entities/responses when a property is repeatedly consumed.

They provide:

- one named place to read backend entity fields
- strongly typed property access
- refactor-safe reuse
- a place to centralize null/default semantics when explicitly required
- consistency across mappers, selectors and hooks

Readers are pure. They do not call APIs, mutate data, make authorization decisions, or contain UI behavior.

A generic property-reader helper may live in `packages/base` and may internally use Lodash:

```ts
import _property from 'lodash/property';

export const createPropertyReader = <
  TEntity,
  TKey extends keyof TEntity,
>(key: TKey): ((entity: TEntity) => TEntity[TKey]) => {
  return _property(key as string) as (
    entity: TEntity,
  ) => TEntity[TKey];
};
```

Backend response reader:

```ts
import { createPropertyReader } from '@base/readers';
import type { SomethingResponse } from '../../types/response';

const readSomethingId = createPropertyReader<
  SomethingResponse,
  'somethingId'
>('somethingId');

const readDisplayName = createPropertyReader<
  SomethingResponse,
  'displayName'
>('displayName');

export const somethingResponseReader = Object.freeze({
  somethingId: readSomethingId,
  displayName: readDisplayName,
});
```

Consumer:

```ts
const somethingId =
  somethingResponseReader.somethingId(response);
```

Do not recreate:

```ts
const getSomethingId = (response: SomethingResponse) =>
  response.somethingId;
```

in ten different files when the backend entity already has a shared reader.

Readers should normally be associated with the service/backend entity they represent.

---

## 18. Service Responsibility and Placement

Services handle transport concerns and expose stable typed business-facing methods.

Typical responsibilities:

- endpoint
- HTTP method
- request serialization
- response parsing
- cancellation
- normalized transport errors
- transport telemetry
- request/response contracts
- response readers for backend entities

Services should not contain large UI-specific transformations.

Preferred flow:

```text
Backend API
    ↓
Domain/Business Service
    ↓
Response Type + Reader
    ↓
Mapper / Selector / Hook
    ↓
ViewModel
    ↓
Presentational Component
```

Placement decision:

```text
Used by one domain
    -> packages/domain/<domain>/services

Shared by multiple domains as the same business contract
    -> packages/business/services/<entity>

Component-specific API call
    -> still domain/business service; never component/services
```

Do not move an API into `business` only to avoid imports. Shared placement must represent true shared ownership.

---

## 19. Request / Response Types Are First-Class Contracts

```ts
export interface SomethingResponse {
  somethingId: string;
  displayName: string;
  createdAt: string;
}
```

Do not duplicate a field type when it already exists.

Bad:

```ts
export interface SomethingProps {
  somethingId: string;
}
```

Better:

```ts
export interface SomethingProps {
  somethingId: SomethingResponse['somethingId'];
}
```

Also prefer:

```ts
export type SomethingSummary = Pick<
  SomethingResponse,
  'somethingId' | 'displayName'
>;
```

or:

```ts
export type SomethingWithoutMetadata = Omit<
  SomethingResponse,
  'metadata'
>;
```

---

## 20. Nearest Authoritative Type Rule

When the same concept exists at multiple levels, reuse the nearest authoritative top-level type.

Priority:

```text
Canonical domain entity / API contract
        ↓
Domain public type
        ↓
Component props/view model
        ↓
Local helper type
```

Example:

```ts
type DeploymentId = DeploymentResponse['deploymentId'];
```

not:

```ts
type DeploymentId = string;
```

---

## 21. Type Definitions Stay Adjacent in `types/`

Implementation files should focus on implementation.

Helper:

```ts
export const getSomething: GetSomething = (params) => {
  // implementation
};
```

Type:

```ts
export type GetSomething = (
  params: GetSomethingParams,
) => GetSomethingResult;
```

Params/result:

```ts
export interface GetSomethingParams {
  somethingId: SomethingResponse['somethingId'];
}
export interface GetSomethingResult {
  label: string;
  isActive: boolean;
}
```

Reusable helpers, services, selectors, mappers, factories, and public callbacks should have explicit callable types.

---

## 22. Constants

Scalar constants use `SCREAMING_SNAKE_CASE`.

```ts
export const HELLO_WORLD = 'HELLO_WORLD';
```

Use `as const` for strongly typed constant maps.

```ts
export const EXECUTION_STATUS = {
  CREATED: 'CREATED',
  RUNNING: 'RUNNING',
  FAILED: 'FAILED',
} as const;
```

Derived type:

```ts
export type ExecutionStatus =
  (typeof EXECUTION_STATUS)[keyof typeof EXECUTION_STATUS];
```

If runtime freeze is required:

```ts
export const EXECUTION_STATUS = Object.freeze({
  CREATED: 'CREATED',
  RUNNING: 'RUNNING',
  FAILED: 'FAILED',
} as const);
```

Prefer `as const` object maps over TypeScript `enum` unless interoperability requires an enum.

---

## 23. Constant Folder Structure

```text
constants/
├── executionStatus.constant.ts
├── types/
│   └── executionStatus.constant.types.ts
└── index.ts
```

Derive union types from the constant source instead of manually duplicating string unions.

---

## 24. Helpers

Helpers should be small, deterministic where possible, and single-purpose.

```text
helpers/
├── deployment.request.mapper.ts
├── deployment.response.mapper.ts
├── deployment.status.selector.ts
├── types/
│   ├── deployment.request.mapper.types.ts
│   ├── deployment.response.mapper.types.ts
│   └── deployment.status.selector.types.ts
├── __tests__/
│   ├── deployment.request.mapper.test.ts
│   ├── deployment.response.mapper.test.ts
│   └── deployment.status.selector.test.ts
└── index.ts
```

Avoid meaningless files:

```text
helper.ts
common.ts
misc.ts
functions.ts
```

---

## 25. Helpers vs Utils and Lodash Usage

Use `helpers` for domain/use-case-aware logic:

- request mapper
- response mapper
- view-model builder
- selector
- domain parser
- domain formatter

Use `utils` for generic technical operations:

- array utility
- string utility
- date/time utility
- generic object utility
- generic calculation
- generic retry calculation

Pure generic utilities belong in `packages/base`. Shared business-aware utilities belong in `packages/business/utils/<business-area>`. Domain-specific utilities remain in the domain.

Do not place domain behavior into `packages/base`.

## Lodash rule

Lodash is approved for object/collection operations where it makes behavior clearer and more reliable.

Always import the individual function.

Correct:

```ts
import _noop from 'lodash/noop';
import _groupBy from 'lodash/groupBy';
import _property from 'lodash/property';
import _uniqBy from 'lodash/uniqBy';
```

Forbidden:

```ts
import _ from 'lodash';
import * as lodash from 'lodash';
```

The underscore-prefixed local alias makes Lodash usage immediately visible during review.

Do not use Lodash merely to replace a trivial, clearer native expression. Use it for consistent reusable calculations/object operations where it improves readability or correctness.

---

## 26. Use-Case-Based File Naming

Names must communicate intent.

Good:

```text
payloadCreation.request.mapper.ts
applicationDetails.response.mapper.ts
resourceTable.columns.builder.ts
deploymentStatus.badge.mapper.ts
migrationPlan.step.selector.ts
```

Bad:

```text
helper1.ts
data.ts
mapper.ts
util.ts
common.ts
```

---

## 27. Unit Tests

Use **Jest** for unit testing.

Use **React Testing Library** for component behavior tests.

Pure functions require adjacent tests.

```text
helpers/
├── buildDeploymentPayload.ts
└── __tests__/
    └── buildDeploymentPayload.test.ts
```

Test:

\- happy path

\- empty values

\- optional values

\- invalid values

\- boundaries

\- mapping correctness

\- immutability

\- error behavior

\- unsupported strategy behavior

---

## 28. `switch` Is Restricted

Do not use `switch` for behavior selection.

Bad:

```ts
switch (rendererType) {
  case RENDERER_TYPES.TABLE:
    return renderTable(params);
  case RENDERER_TYPES.FORM:
    return renderForm(params);
  default:
    return null;
}
```

Use a strategy map:

```ts
const RENDERER_FACTORY = {
  [RENDERER_TYPES.TABLE]: renderTable,
  [RENDERER_TYPES.FORM]: renderForm,
} satisfies Record<RendererType, RendererStrategy>;
```

```ts
export const executeRenderer: ExecuteRenderer = ({
  type,
  params,
}) => {
  const renderer = RENDERER_FACTORY[type];
  return renderer(params);
};
```

`Record` + `satisfies` should be used where practical to make the registry exhaustive.

---

## 29. Avoid `if / else if` Chains

Bad:

```ts
if (role === 'ADMIN') {
  return adminConfig;
} else if (role === 'DEVELOPER') {
  return developerConfig;
} else if (role === 'VIEWER') {
  return viewerConfig;
}
```

Prefer:

```ts
const ROLE_CONFIG = {
  [ROLE.ADMIN]: adminConfig,
  [ROLE.DEVELOPER]: developerConfig,
  [ROLE.VIEWER]: viewerConfig,
} satisfies Record<Role, RoleConfig>;
```

```ts
const config = ROLE_CONFIG[role];
```

A simple guard `if` is acceptable:

```ts
if (!applicationId) {
  return null;
}
```

Prefer guard clauses over nested logic.

---

## 30. Strategy / Factory Pattern

When behavior changes by type, use a strategy registry.

```ts
const ACTION_STRATEGIES = {
  [ACTION_TYPE.BUILD]: executeBuild,
  [ACTION_TYPE.TEST]: executeTest,
  [ACTION_TYPE.DEPLOY]: executeDeploy,
} satisfies Record<ActionType, ActionStrategy>;
```

```ts
export const executeAction: ExecuteAction = ({
  type,
  params,
}) => {
  const strategy = ACTION_STRATEGIES[type];
  return strategy(params);
};
```

For untrusted runtime input, a guard is fine:

```ts
const strategy = ACTION_STRATEGIES[type];
if (!strategy) {
  throw new UnsupportedActionTypeError(type);
}
return strategy(params);
```

---

## 31. Configuration-Driven Rendering

Repeated UI structures should be driven by configuration.

```ts
const SECTIONS = [
  {
    id: 'summary',
    renderer: ApplicationSummary,
    renderOptions: {
      applicationId,
    },
  },
  {
    id: 'resources',
    renderer: ResourceTable,
    renderOptions: {
      resources,
    },
  },
] satisfies RendererConfig[];
```

Generic renderer:

```tsx
export const ConfigRenderer: ConfigRendererProps = ({
  items,
}) => (
  <>
    {items.map(({ renderer: Renderer, renderOptions, id }) => (
      <Renderer key={id} {...renderOptions} />
    ))}
  </>
);
```

---

## 32. Typed Renderer Contract

```ts
export interface RendererConfig<
  TProps extends object = Record<string, never>,
> {
  id: string;
  renderer: ComponentType<TProps>;
  renderOptions: TProps;
}
```

For heterogeneous renderer collections, use a discriminated contract rather than falling back to `any`.

---

## 33. Shared Platform Components and Storybook

Repeated UI patterns belong in `packages/components`.

Priority shared units include:

```text
FormBuilder
TableManager
DataGrid
FilterBuilder
SearchBar
PageHeader
PageShell
DetailsPanel
DrawerManager
ModalManager
Wizard
StepManager
Tabs
Renderer
FieldRenderer
ActionRenderer
StatusBadge
EmptyState
LoadingState
ErrorState
PermissionGuard
FeatureGuard
ConfirmAction
EntitySelector
EnvironmentSelector
ContextSelector
CodeEditor
JsonViewer
YamlViewer
DiffViewer
Terminal
LogViewer
MetricPanel
TraceViewer
Timeline
ExecutionGraph
DependencyGraph
KeyValueViewer
ObjectInspector
```

Do not duplicate these patterns domain by domain.

Every public reusable component in `packages/components` must include an adjacent Storybook story unless the unit is non-visual.

Recommended component package:

```text
packages/components/components/TableManager/
├── TableManager.tsx
├── TableManagerContainer.tsx        # only when shared orchestration is truly needed
├── TableManager.stories.tsx
├── __tests__/
├── components/
├── hooks/
├── helpers/
├── constants/
├── types/
├── README.md                        # optional for simple units
└── index.ts
```

Stories should cover meaningful states such as:

- default
- loading
- empty
- error
- disabled/read-only
- large data where relevant
- important variants
- interaction examples where useful

Storybook stories are part of the component contract and should be reviewed with the implementation.

Shared components must not call domain services directly. They receive data/configuration/callbacks through typed props.

---

## 34. FormBuilder

Do not repeatedly hand-code the same form infrastructure.

```ts
export const APPLICATION_FORM_CONFIG = [
  {
    name: 'applicationName',
    type: FIELD_TYPE.TEXT,
    label: 'Application name',
  },
  {
    name: 'environmentId',
    type: FIELD_TYPE.SELECT,
    label: 'Environment',
  },
] as const;
```

```tsx
<FormBuilder
  fields={APPLICATION_FORM_CONFIG}
  onSubmit={handleSubmit}
/>
```

Domain-specific renderers should be registered without modifying the core FormBuilder.

---

## 35. TableManager

Do not reimplement sorting/filtering/selection/pagination/actions/loading/empty/error for every table.

```tsx
<TableManager
  columns={columns}
  rows={rows}
  rowActions={rowActions}
  filters={filters}
  pagination={pagination}
/>
```

---

## 36. View Models

Do not pass raw API responses deeply through the UI when the view needs a different shape.

Preferred flow:

```text
API Response
     ↓
response mapper
     ↓
ViewModel
     ↓
Presentational Component
```

Example:

```ts
const viewModel =
  mapDeploymentResponseToDeploymentViewModel(response);
```

---

## 37. Request Mapping

Do not construct complex API payloads inline in components.

```text
Form / UI state
     ↓
request mapper
     ↓
Service Request
     ↓
Service
```

```ts
const request =
  mapCreateApplicationFormToCreateApplicationRequest(formValues);
await createApplication(request);
```

---

## 38. Domain State Management

Each domain owns its state implementation.

```text
domain/core/state/
├── domain.slice.ts
├── domain.actions.ts
├── domain.selectors.ts
├── domain.middleware.ts
├── types/
└── index.ts
```

Rules:

\- consumers use exported actions/selectors

\- implementation details remain private

\- large entity collections should be normalized where useful

\- derived state belongs in selectors/hooks

\- do not duplicate the same derived state in multiple places

\- the domain owns its slice/reducer/actions/selectors and domain middleware

\- the domain exports only the registration metadata the shell needs

\- `packages/shell` constructs the shared Redux store and mounts the global Redux provider

\- domains must not import the shell store or shell internals

\- shell code must not directly mutate or interpret domain state

---

## 39. Server State vs Client State

Use the query layer for:

\- API/server state

\- caching

\- stale state

\- retry/refetch

\- request deduplication

Use Redux/domain state for:

\- cross-screen workflows

\- complex client orchestration

\- editors

\- coordinated UI state

\- long-lived domain client state

Do not automatically copy every API response into Redux.

---

## 40. Selectors

Non-trivial derived Redux state belongs in selectors.

Bad:

```tsx
const activeItems = items
  .filter(...)
  .map(...)
  .sort(...);
```

repeated in several components.

Prefer:

```ts
export const selectActiveItems: SelectActiveItems = (state) => {
  ...
};
```

Non-trivial pure selectors require unit tests.

---

## 41. No Vendor Logic in Presentation Code

Bad:

```ts
if (provider === 'DATADOG') {
  ...
}
```

Prefer capability-based behavior:

```ts
const renderer = OBSERVABILITY_RENDERERS[capability];
```

or registration:

```ts
registerObservabilityPanel({
  capability: OBSERVABILITY_CAPABILITY.TRACES,
  renderer: TracePanel,
});
```

UI should depend on normalized contracts/capabilities, not vendor names.

---

## 42. Registry Pattern for Extensible UI

Extensible areas should use registration rather than mega-components.

Example:

```ts
registerWorkbenchPanel({
  id: WORKBENCH_PANEL.LOGS,
  title: 'Logs',
  renderer: LogsPanel,
  capability: CAPABILITY.LOGS,
});
```

Use this pattern for:

\- Workbench panels

\- field renderers

\- action renderers

\- table cell renderers

\- domain routes

\- dashboard widgets

\- command palette actions

\- test-result renderers

\- provider configuration forms

---

## 43. Domain Route Registration

Domains expose route configuration.

```ts
export const deliveryRoutes: DomainRouteConfig[] = [
  ...
];
```

`packages/shell` combines domain route registries through domain public registrations.

`frontend/app` only mounts the shell; it does not manually compose domain routes.

Do not build a single giant route file containing the whole product.

---

## 44. Feature Registration

Feature packages should expose mountable contracts where useful.

```ts
export interface FrontendDomainRegistration {
  id: DomainId;
  routes: DomainRouteConfig[];
  navigation: NavigationConfig[];
  commands?: CommandConfig[];
}
```

The domain exports this registration through its root public API. `packages/shell` consumes the registrations and binds the application together.

This gives us a path toward extension packages/marketplace modules later without turning the shell into business logic.

---

## 45. Permission Configuration

Do not scatter permission strings throughout JSX.

```ts
export const ACTION_PERMISSION_MAP = {
  [ACTION.DEPLOY]: PERMISSION.DEPLOYMENT_EXECUTE,
  [ACTION.ROLLBACK]: PERMISSION.DEPLOYMENT_ROLLBACK,
  [ACTION.VIEW_LOGS]: PERMISSION.LOG_READ,
} as const;
```

Use normalized hooks/guards.

```ts
const canDeploy = usePermission(PERMISSION.DEPLOYMENT_EXECUTE);
```

Frontend permission checks improve UX; backend authorization remains authoritative.

---

## 46. Feature Flags

Feature flag IDs are constants.

```ts
export const FEATURE_FLAG = {
  WORKBENCH_PROFILING: 'WORKBENCH_PROFILING',
} as const;
```

Use a normalized feature hook:

```ts
const isEnabled =
  useFeature(FEATURE_FLAG.WORKBENCH_PROFILING);
```

Do not spread vendor feature-flag SDK calls through the UI.

---

## 47. Error Normalization

UI code consumes normalized errors.

```ts
export interface UiError {
  code: string;
  title: string;
  message: string;
  retryable: boolean;
}
```

Components must not parse raw backend/provider error strings.

---

## 48. Async UI States

Every async unit should deliberately handle:

```text
idle
loading
refreshing
success
empty
partial
error
```

Use reusable primitives:

```tsx
<LoadingState />
<EmptyState />
<ErrorState />
```

---

## 49. Pure Function Preference

Transformation logic should be isolated.

Instead of mapping inside a component:

```tsx
const rows = response.items.map(...);
```

prefer:

```ts
const rows =
  mapApplicationResponseItemsToTableRows(response.items);
```

This improves naming, reuse, testability, and SRP.

---

## 50. Immutability

Do not mutate inputs.

Bad:

```ts
items.sort(compareItems);
```

Prefer:

```ts
const sortedItems = [...items].sort(compareItems);
```

Redux Toolkit/Immer may handle controlled state mutation inside its intended abstraction.

---

## 51. Magic Values

Avoid raw literals in business logic.

Bad:

```ts
if (status === 'RUNNING') {
  ...
}
setTimeout(refresh, 30000);
```

Prefer:

```ts
if (status === EXECUTION_STATUS.RUNNING) {
  ...
}
setTimeout(
  refresh,
  EXECUTION_REFRESH_INTERVAL_MS,
);
```

---

## 52. Naming Conventions

### React components

```text
PascalCase
```

### Hooks

```text
useCamelCase
```

### Functions

```text
camelCase
```

### Constants

```text
SCREAMING_SNAKE_CASE
```

### Types/interfaces

```text
PascalCase
```

### Good filenames

```text
ApplicationCard.tsx
ApplicationCardContainer.tsx
application.response.mapper.ts
application.request.mapper.ts
workbenchPanel.renderer.constant.ts
useApplication.ts
withPermission.tsx
```

---

## 53. Boolean Naming

Use natural boolean names:

```ts
isLoading
isEnabled
hasPermission
canDeploy
shouldRefresh
wasCancelled
```

Avoid ambiguous names such as `flag`, `value`, or `allowedValue`.

---

## 54. Event Handler Naming

Props:

```ts
onSubmit
onCancel
onRefresh
onRowClick
```

Internal handlers:

```ts
handleSubmit
handleCancel
handleRefresh
handleRowClick
```

---

## 55. Component Cleanliness Rule

A presentational component should mostly read like this:

```tsx
export const ApplicationView: ApplicationViewProps = ({
  header,
  rows,
  actions,
  isLoading,
}) => (
  <ApplicationTemplate>
    <PageHeader {...header} actions={actions} />
    <TableManager
      rows={rows}
      loading={isLoading}
    />
  </ApplicationTemplate>
);
```

If the component contains fetching, mapping, permission resolution, request creation, feature resolution, several effects, or provider-specific branches, move those responsibilities out.

---

## 56. No God Components / No God Hooks

A component should not simultaneously own:

\- API calls

\- state machine logic

\- Redux orchestration

\- routing

\- permissions

\- mapping

\- mutations

\- layout

\- dialogs

\- error parsing

Do not solve a God Component by moving everything into one `useEverything()` hook.

Create focused units and compose them.

---

## 57. Component / Function Size Guidance

Review thresholds, not blind hard limits:

```text
Presentational component: \~150 lines
Container:                \~150 lines
Hook:                     \~150 lines
Helper:                   \~100 lines
Service:                  \~100 lines
```

If a file grows substantially, review SRP.

---

## 58. Reusability Hierarchy

Code should live at the lowest correct reusable level.

```text
Platform-wide reusable primitive
            ↑
Cross-domain shared unit
            ↑
Domain reusable unit
            ↑
Page reusable unit
            ↑
Component-local unit
```

Do not put everything into `shared` prematurely.

---

## 59. Promotion Rule

When a local unit becomes reusable:

1. stabilize the contract

2. remove domain-specific assumptions

3. expose generic inputs

4. add tests

5. move/promote to the nearest reusable package

6. preserve compatibility through a re-export when existing consumers already depend on the previous path

This lets architecture improve without breaking all existing consumers.

---

## 60. Avoid Premature Generic Abstractions

OCP does not mean building a generic framework for every two-line function.

Promote an abstraction when:

\- the pattern repeats

\- the extension axis is clear

\- inputs/outputs are sufficiently stable

\- the abstraction reduces conceptual duplication

Do not create an abstraction that is harder to understand than the code it replaces.

---

## 61. Base, Business, Components, and Shell Package Rules

The old generic `platform/shared` split is replaced by three explicit package tiers.

### `packages/base`

Purpose: pure technical building blocks.

```text
packages/base/
├── constants/
│   ├── empty.constant.ts
│   └── index.ts
├── helpers/
├── utils/
│   ├── dateTime/
│   ├── object/
│   ├── array/
│   └── index.ts
├── readers/
│   ├── createPropertyReader.ts
│   └── index.ts
├── types/
├── RUNBOOK.md
└── index.ts
```

Examples:

```ts
export const EMPTY_OBJECT = Object.freeze({});
export const EMPTY_ARRAY = Object.freeze([]);
```

`base` must be:

- domain-independent
- business-independent
- UI-independent
- deterministic/pure wherever practical
- safe for every other frontend package to consume

### `packages/business`

Purpose: shared business logic and reusable backend entity contracts.

```text
packages/business/
├── services/
│   ├── application/
│   ├── environment/
│   └── ...
├── readers/
├── helpers/
├── utils/
│   ├── application/
│   ├── authorization/
│   └── ...
├── constants/
├── types/
├── RUNBOOK.md
└── index.ts
```

`business` may depend on `base` but must not depend on `domain` or shared UI components.

### `packages/components`

Purpose: common reusable platform UI.

```text
packages/components/
├── components/
├── templates/
├── builders/
├── renderers/
├── hooks/
├── hocs/
├── forms/
├── tables/
├── layouts/
├── navigation/
├── feedback/
├── permissions/
├── types/
├── RUNBOOK.md
└── index.ts
```

`components` may depend on `base`. It must not depend on a product domain or `shell`.

Business-aware UI that genuinely belongs to one domain should remain in that domain and may be exposed through the domain public API if other domains need it.

### `packages/shell`

Purpose: application composition and shared runtime initialization.

```text
packages/shell/
├── Shell.tsx
├── ShellContainer.tsx
├── providers/
├── state/
├── registry/
├── layout/
├── skeleton/
├── errors/
├── constants/
├── helpers/
├── types/
├── RUNBOOK.md
└── index.ts
```

`shell` may depend on:

```text
base
components
domain root public APIs
```

It must not deep-import domain internals.

It must not own reusable business services or product use cases.

`shell` is the only package responsible for binding all domain registrations into the assembled frontend runtime.

---

## 62. Public Shared Packages, Stories, and Runbooks

`base`, `business`, `components`, `shell`, and every domain are explicit package boundaries rather than dumping grounds.

Every public shared/composition package must have a `RUNBOOK.md` at package root.

The runbook should explain:

- package purpose
- ownership
- allowed dependency direction
- public import paths
- how to add a new unit
- how to promote a domain unit into the shared package
- common failure modes
- test commands
- Storybook commands where applicable
- deprecation/migration policy
- examples of correct and forbidden usage

For complex/highly reusable units, a unit-level `RUNBOOK.md` or `README.md` should be added beside the unit.

Example:

```text
packages/components/
├── RUNBOOK.md
└── components/
    └── FormBuilder/
        ├── FormBuilder.tsx
        ├── FormBuilder.stories.tsx
        ├── README.md
        └── index.ts
```

Stories are mandatory for public visual components. `packages/shell` also requires Storybook coverage for public visual infrastructure such as `ShellSkeleton`. Runbooks describe engineering/operational usage; stories demonstrate visual/component contracts.

`index.ts` remains the only normal runtime public surface.

---

## 63. Circular Dependencies Are Forbidden

The build/CI should fail on circular dependencies.

Resolve cycles by:

\- moving true shared contracts downward

\- injecting callbacks/contracts

\- using public package APIs

\- using registration/configuration

\- splitting mixed responsibilities

Do not hide cycles using arbitrary lazy imports.

---

## 64. TypeScript Strictness

Recommended baseline:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noImplicitOverride": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "useUnknownInCatchVariables": true
  }
}
```

Avoid `any`.

Use `unknown` and narrow safely.

---

## 65. Avoid Unsafe Casts

Bad:

```ts
const response = data as ApplicationResponse;
```

Prefer schema validation at unsafe boundaries where appropriate.

Use `satisfies` to validate shapes without destroying inference.

```ts
const CONFIG = {
  ...
} satisfies ApplicationConfig;
```

---

## 66. API Boundary Validation

Network data is not trustworthy just because a TypeScript type says so.

For high-risk/external contracts, validate at the boundary.

```ts
const response =
  ApplicationResponseSchema.parse(rawResponse);
```

Do not scatter validation throughout rendering components.

---

## 67. Generated API Types

If OpenAPI generation is introduced:

\- generated models live in a clearly generated package

\- generated files are never manually edited

\- domain public types may wrap/alias generated contracts

\- consumers still use domain public boundaries, not generated internals everywhere

---

## 68. API Client

All services use one shared API client abstraction.

```ts
api.get(...)
api.post(...)
api.put(...)
api.patch(...)
api.delete(...)
```

The API client owns:

\- base URL

\- auth attachment

\- correlation IDs

\- serialization

\- normalized transport errors

\- cancellation

\- retry policy where appropriate

\- telemetry

Services do not recreate those behaviors.

---

## 69. Cancellation

Long-running or rapidly changing requests should support cancellation.

Examples:

\- typeahead search

\- route change

\- Workbench context switch

\- tab change

\- repeated polling

Obsolete responses must not overwrite newer state.

---

## 70. Security Rules

\- never hardcode secrets

\- never place secrets in query strings

\- sanitize untrusted HTML

\- do not render arbitrary raw HTML

\- frontend permission checks never replace backend authorization

\- prevent unsafe redirect URLs

\- redact sensitive values from logs

\- avoid unnecessary persistence of sensitive data

\- do not expose internal stack traces

\- avoid `eval` and dynamic code execution

\- follow CSP requirements

---

## 71. Accessibility

Platform primitives should support accessibility by default:

\- semantic HTML

\- keyboard navigation

\- focus management

\- visible focus

\- accessible labels

\- form-error association

\- screen-reader announcements for important live states

\- modal focus trapping

\- correct table semantics

\- ARIA only when needed

Accessibility should live in reusable primitives so every domain inherits it.

---

## 72. Performance

\- lazy-load domain routes

\- lazy-load heavy Workbench panels

\- virtualize large tables/lists

\- avoid unnecessary re-renders

\- memoize only when useful

\- debounce high-frequency user inputs

\- cancel obsolete requests

\- avoid duplicate derived state

\- avoid huge shared bundles

\- keep provider-specific code split where practical

Do not apply `useMemo`/`useCallback` blindly.

---

## 73. Observability

Frontend telemetry goes through one abstraction.

```ts
trackUiEvent(
  UI_EVENT.DEPLOYMENT_STARTED,
  {
    applicationId,
    environmentId,
  },
);
```

Do not scatter vendor analytics SDKs throughout the frontend.

---

## 74. Logging

Use a common frontend logger abstraction.

Avoid committed production `console.log` calls.

The logger should support:

\- level

\- correlation ID

\- domain

\- safe metadata

\- redaction

\- environment behavior

---

## 75. Capability-Based UI

Render features based on normalized capabilities, not vendor names.

Good:

```ts
const canProfile =
  capabilities.includes(CAPABILITY.PROFILING_CPU);
```

Bad:

```ts
const canProfile = provider === 'DATADOG';
```

---

## 76. OCP Review Rule

When adding a new type, renderer, provider, panel, or action ask:

> Can I add this by registering/configuring a new implementation, or do I have to edit a giant existing conditional?

If a central conditional must be edited every time, the extension point probably needs redesign.

---

## 77. SRP Review Rule

Ask:

> Can this file's responsibility be described in one short sentence?

Good:

> Maps deployment API response to a deployment view model.

Bad:

> Fetches deployment, checks permissions, maps status, creates actions, tracks telemetry, renders UI, and opens dialogs.

---

## 78. DRY Review Rule

DRY means avoiding duplication of the **same concept that should evolve together**.

It does not mean extracting every two similar lines into a generic helper.

---

## 79. Styling and Design Tokens

Use one primary styling strategy consistently.

Platform UI primitives should own their styles.

Avoid domain CSS that modifies platform internals.

Use shared design tokens for:

\- spacing

\- typography

\- radius

\- elevation

\- breakpoints

\- semantic colors

\- z-index layers

\- motion durations

---

## 80. Stable IDs, Not Labels

Do not use user-visible labels as logic keys.

Bad:

```ts
if (tab.label === 'Logs') {
  ...
}
```

Prefer:

```ts
if (tab.id === WORKBENCH_TAB.LOGS) {
  ...
}
```

Labels can change. IDs are contracts.

---

## 81. Date / Time

Use one approved date/time abstraction.

\- treat backend timestamps explicitly

\- format through reusable helpers/components

\- do not manually concatenate date strings

\- timezone-sensitive behavior must be explicit

---

## 82. Optional Values

Do not use empty strings as a universal replacement for missing data.

Recommended convention:

\- `undefined` for absent optional frontend values

\- `null` only where the external/domain contract explicitly uses semantic `null`

Normalize external values in mappers.

---

## 83. Full Domain Example

```text
frontend/packages/domain/delivery/
├── core/
│   ├── routes/
│   ├── state/
│   ├── providers/
│   └── index.ts
│
├── pages/
│   ├── ReleaseDetailsPage/
│   └── PipelineRunsPage/
│
├── templates/
│   └── ReleaseDetailsTemplate/
│
├── components/
│   ├── ReleaseSummary/
│   ├── PipelineRunTable/
│   └── DeploymentStatus/
│
├── hooks/
│   ├── useRelease/
│   └── usePipelineRuns/
│
├── services/
│   ├── getRelease.service.ts
│   ├── startRelease.service.ts
│   └── types/
│
├── helpers/
├── constants/
├── types/
└── index.ts
```

---

## 84. Full Component Example

```text
DeploymentStatus/
├── DeploymentStatus.tsx
├── DeploymentStatusContainer.tsx
│
├── components/
│   └── DeploymentStatusBadge/
│       ├── DeploymentStatusBadge.tsx
│       ├── types/
│       └── index.ts
│
├── constants/
│   ├── deploymentStatus.constant.ts
│   ├── types/
│   │   └── deploymentStatus.constant.types.ts
│   └── index.ts
│
├── helpers/
│   ├── deploymentStatus.response.mapper.ts
│   ├── types/
│   │   └── deploymentStatus.response.mapper.types.ts
│   ├── __tests__/
│   │   └── deploymentStatus.response.mapper.test.ts
│   └── index.ts
│
├── hooks/
│   └── useDeploymentStatus/
│       ├── useDeploymentStatus.ts
│       ├── helpers/
│       ├── types/
│       └── index.ts
│
├── services/
│   ├── getDeploymentStatus.service.ts
│   ├── types/
│   │   ├── request/
│   │   ├── response/
│   │   └── getDeploymentStatus.service.types.ts
│   └── index.ts
│
├── types/
│   ├── DeploymentStatus.props.ts
│   ├── DeploymentStatus.viewModel.ts
│   └── index.ts
│
└── index.ts
```

---

## 85. Factory Example

```ts
export const ACTION_FACTORY = {
  [ACTION_TYPE.CREATE]: createAction,
  [ACTION_TYPE.UPDATE]: updateAction,
  [ACTION_TYPE.DELETE]: deleteAction,
} satisfies Record<ActionType, ActionExecutor>;
```

```ts
export const executeAction: ExecuteAction = ({
  actionType,
  params,
}) => {
  const execute = ACTION_FACTORY[actionType];
  return execute(params);
};
```

---

## 86. Config-Driven Actions

```ts
const ACTION_CONFIG = [
  {
    id: ACTION.REFRESH,
    label: 'Refresh',
    icon: RefreshIcon,
    onExecute: refresh,
  },
  {
    id: ACTION.DEPLOY,
    label: 'Deploy',
    icon: DeployIcon,
    permission: PERMISSION.DEPLOYMENT_EXECUTE,
    onExecute: deploy,
  },
] satisfies ActionConfig[];
```

Render through a reusable `ActionRenderer`.

---

## 87. Composition Instead of Boolean-Prop Explosion

Avoid:

```tsx
<Table
  isMigration
  isReadOnly
  showCost
  useCompactMode
  enableProviderActions
/>
```

Prefer composable configuration:

```tsx
<TableManager
  columns={columns}
  actions={actions}
  capabilities={capabilities}
  layout={layout}
/>
```

A growing list of behavioral booleans usually means the abstraction is closed to extension.

---

## 88. State Colocation

Keep state as close as possible to the feature that owns it.

Promote state upward only when multiple consumers truly need the same source.

Do not place every value in Redux.

---

## 89. Side Effects

Side effects must be easy to locate.

Expected locations:

\- services

\- hooks

\- middleware

\- containers

\- explicitly named effect modules

Presentational components should not hide side effects.

---

## 90. Reusable Fetch Hook Contract

A data-fetch hook should expose a semantic frontend contract, not query-library internals.

```ts
export interface UseFetchApplicationResult {
  application: ApplicationViewModel | undefined;
  isLoading: boolean;
  error: UiError | undefined;
  refresh: () => Promise<void>;
}
```

---

## 91. Testing by Unit Type

### Pure helper

Direct Jest unit test.

### Mapper

Test every important transformation.

### Selector

Test with controlled state fixtures.

### Hook

Test behavior with React Testing Library/hook utilities.

### Presentational component

Test visible behavior and callbacks.

### Container

Test orchestration at meaningful boundaries; do not duplicate lower-level tests.

### Page

Test important integrated user flows, not private implementation details.

---

## 92. Domain Mock Mechanism

Mocks are first-class package boundaries for testing, Storybook and local development.

Each domain may expose a dedicated mock entry point:

```text
packages/domain/catalog/
├── mock/
│   ├── builders/
│   │   └── application.builder.ts
│   ├── fixtures/
│   │   └── application.fixture.ts
│   ├── handlers/
│   │   └── application.msw.ts
│   ├── services/
│   │   └── application.service.mock.ts
│   ├── state/
│   └── index.ts
├── index.ts
└── RUNBOOK.md
```

Runtime imports use:

```ts
import {
  ApplicationSelector,
} from '@domain/catalog';
```

Tests/stories/dev-only code may use:

```ts
import {
  buildApplication,
  catalogHandlers,
} from '@domain/catalog/mock';
```

A consuming domain must not mock another domain by importing its internal files.

Preferred boundaries:

- **MSW** for HTTP/network mocks
- service mocks for isolated hook/container tests
- typed builders for backend response entities
- fixtures for stable examples
- Jest mocks only at actual public boundaries when required

Production runtime code must not import `@domain/*/mock`.

CI should include an architecture check preventing mock entry points from being imported by production source.

Mock data must use the same public request/response types as the real domain/service contracts so contract drift becomes visible.

Mock true boundaries, not internal helper implementation details.

---

## 93. Test Data Builders

Prefer builders/factories over repeated large test objects.

```ts
const application = buildApplicationResponse({
  name: 'CRM',
});
```

---

## 94. Snapshot Tests

Do not use giant snapshots as the main correctness mechanism.

Prefer explicit behavior assertions.

Small focused snapshots are allowed when they add clear value.

---

## 95. ESLint / Architecture Enforcement

The repository should automate these rules where practical:

\- strict TypeScript

\- no explicit `any`

\- restricted `switch`

\- package-boundary imports

\- no circular dependencies

\- React Hooks rules

\- unused imports/variables

\- import ordering

\- consistent type imports

\- complexity limits

\- restricted deep imports

\- restricted `console`

\- restricted direct API client use outside services

\- restricted cross-domain internal imports

\- forbidden `domain -> shell` imports

\- forbidden `components -> shell` and `business -> shell` imports

\- shell imports from domains restricted to domain root public APIs

\- app domain composition restricted in favor of `packages/shell`

Rules should not rely only on reviewer memory.

---

## 96. Formatting

Use one automatic formatter.

CI verifies formatting.

Formatting should not consume review time.

---

## 97. Import Ordering

Recommended groups:

```text
1. React / external dependencies
2. base
3. business
4. components
5. domain public packages
6. shell (only from app/shell-owned code; domains must never import it)
7. local unit imports
8. type imports
9. styles/assets where applicable
```

Automate ordering.

---

## 98. Path Aliases and Public Import Paths

Use stable aliases that communicate architectural ownership.

Recommended:

```json
{
  "compilerOptions": {
    "paths": {
      "@app/*": ["app/*"],
      "@base/*": ["packages/base/*"],
      "@business/*": ["packages/business/*"],
      "@components/*": ["packages/components/*"],
      "@domain/*": ["packages/domain/*"],
      "@shell/*": ["packages/shell/*"],
      "@shell": ["packages/shell/index.ts"]
    }
  }
}
```

Examples:

```ts
import { EMPTY_OBJECT } from '@base/constants';
import { applicationService } from '@business/services/application';
import { TableManager } from '@components/components/TableManager';
import { ApplicationSelector } from '@domain/catalog';
import { PerihellionShell } from '@shell';
```

Cross-domain code must use only the domain root public API:

```ts
import { Something } from '@domain/some-domain';
```

Never:

```ts
import { Something }
  from '@domain/some-domain/components/internal/Something';
```

Test/Storybook-only mocks use the explicit secondary entry point:

```ts
import { buildSomething } from '@domain/some-domain/mock';
```

`@shell` is primarily consumed by `frontend/app` and shell-owned code.

Domains must never do:

```ts
import { something } from '@shell';
```

Avoid long `../../../../` imports.

Aliases do not bypass dependency rules. ESLint/architecture tests must still enforce allowed package direction.

---

## 99. Validation

Form validation belongs in reusable schemas/helpers.

Do not scatter validation logic through JSX.

Example:

```ts
export const CREATE_APPLICATION_SCHEMA = ...;
```

Validation errors should map to stable field IDs.

---

## 100. Forms Do Not Persist Directly

Preferred flow:

```text
Form View
   ↓
Container/Hook
   ↓
Request Mapper
   ↓
Service
```

This keeps forms reusable and testable.

---

## 101. Dialog / Drawer Orchestration

Avoid excessive boolean state:

```ts
isDeleteOpen
isEditOpen
isInfoOpen
isRollbackOpen
```

For complex surfaces, prefer a dialog/drawer manager contract.

---

## 102. Command / Action Model

Complex product surfaces such as Workbench should use reusable command contracts.

```ts
export interface UiCommand {
  id: CommandId;
  label: string;
  permission?: PermissionId;
  execute: CommandExecutor;
}
```

The same action can then appear in:

\- buttons

\- menus

\- command palette

\- context menus

\- keyboard shortcuts

---

## 103. Public Contract Stability

Anything exported through `index.ts` is a public contract.

Before breaking it:

1. identify consumers

2. prefer additive changes

3. deprecate first where practical

4. document migration

5. preserve compatibility adapters when inexpensive

---

## 104. Recommended Code Review Questions

Every frontend PR should answer:

\- Is SRP preserved?

\- Is OCP preserved?

\- Is true conceptual duplication avoided?

\- Is the unit consumed through `index.ts`?

\- Are deep imports prevented?

\- Are authoritative types reused?

\- Are API calls isolated in services?

\- Is mapping outside presentation components?

\- Are hooks focused?

\- Is `switch` avoided?

\- Are avoidable `if / else` chains replaced with maps/strategies/config?

\- Is UI configuration-driven where useful?

\- Does a reusable platform primitive already exist?

\- Should this pattern become a platform primitive?

\- Are pure functions tested?

\- Are async states handled?

\- Is accessibility preserved?

\- Are errors normalized?

\- Is vendor-specific logic kept out of presentation code?

\- Is there a circular dependency?

\- Is the public API intentional and stable?

---

## 105. Definition of Done — Frontend Unit

```text
[ ] Responsibility is clear and singular
[ ] Public API is exported through index.ts
[ ] External consumers do not deep-import internals
[ ] TypeScript types are explicit
[ ] Nearest authoritative types are reused
[ ] Primitive type duplication is avoided
[ ] No switch statement exists
[ ] Avoidable if/else chains are replaced by map/strategy/config
[ ] Presentation/container responsibilities are separated
[ ] API calls live in services
[ ] Complex stateful logic lives in hooks/container/state layer
[ ] Pure transformations live in helpers/selectors
[ ] Constants are immutable and strongly typed
[ ] Request/response contracts are typed
[ ] Errors are normalized
[ ] Loading/empty/error states are covered
[ ] Pure functions have Jest tests
[ ] React behavior has appropriate RTL tests
[ ] Accessibility is checked
[ ] No circular dependency exists
[ ] No forbidden deep imports exist
[ ] Lint passes
[ ] TypeScript passes
[ ] Tests pass
[ ] Formatting passes
[ ] Reusable public contract is documented
```

---

## 106. Target Architecture Shape

Preferred application composition flow:

```text
frontend/app
     ↓
packages/shell
     ↓
Domain Registrations
     ↓
Routes + Navigation + Redux Registration
     ↓
Domain Container / Hook
     ↓
Service + State + Permissions
     ↓
Mapper / Selector
     ↓
ViewModel
     ↓
Presentational Component
     ↓
Shared Components / Base Primitives
```

The shell exists above domains. It composes them; domains do not depend on it.

For extensibility:

```text
Configuration
      ↓
Registry
      ↓
Strategy / Renderer
      ↓
Reusable Unit
```

Avoid:

```text
Huge Component
      ↓
API Calls
      ↓
Redux
      ↓
Provider Checks
      ↓
Nested if/else
      ↓
Mapping
      ↓
Duplicated HTML
```

---

## 107. Phase 0 Frontend Deliverables

Before Phase 0 is complete, frontend engineering must demonstrate the final package model:

```text
frontend/
├── app/
└── packages/
    ├── base/
    ├── business/
    ├── components/
    ├── domain/
    │   └── example-domain/
    └── shell/
```

Required reference implementation:

1. `packages/base` with at least one generic utility and immutable constant such as `EMPTY_OBJECT`
2. `packages/business` with one shared typed business/service unit
3. `packages/components` with one reusable component and Storybook stories
4. one example domain under `packages/domain/<name>`
5. `packages/shell` with `Shell.tsx` + `ShellContainer.tsx`
6. shell-owned global provider composition
7. shell-owned Redux store/provider initialization
8. domain Redux registration consumed by the shell without deep imports
9. domain route registration consumed by the shell
10. domain navigation registration consumed by the shell
11. generic application startup `ShellSkeleton` with Storybook coverage
12. `frontend/app` mounting the shell rather than directly composing domains
13. one container/presentational component pair inside the example domain
14. one domain-level service with typed request/response
15. one typed backend response reader
16. one typed helper + adjacent Jest test
17. one focused custom hook
18. one config-driven renderer
19. one strategy/factory map
20. public root `index.ts` boundaries
21. explicit `@domain/<name>/mock` mechanism with typed builder/fixture and MSW handler
22. lint rule preventing cross-domain deep imports
23. lint rule preventing `domain -> shell`
24. lint rule restricting shell imports to domain root public APIs
25. lint/architecture rule preventing component-level API services
26. architecture rule preventing production imports from mock entry points
27. circular dependency detection
28. strict TypeScript
29. approved Lodash per-function import rule
30. Jest + React Testing Library setup
31. Storybook setup
32. `RUNBOOK.md` for `base`, `business`, `components`, and `shell`
33. example domain `RUNBOOK.md`
34. Woodpecker CI frontend checks

The example domain plus `packages/shell` become the reference implementation for humans and coding agents.

The Phase 0 architecture test should prove:

```text
app mounts shell
shell initializes shared providers
shell loads domain registrations
domain reducers/routes/navigation become available
domain never imports shell
shell never deep-imports domain internals
```

---

## 108. Engineering Documentation and Runbooks

These guidelines live with the rest of the repository engineering rules.

```text
docs/
└── engineering/
    ├── README.md
    ├── FRONTEND_CODE_GUIDELINES.md
    ├── BACKEND_CODE_GUIDELINES.md
    ├── API_CONTRACT_GUIDELINES.md
    ├── PROVIDER_SDK_GUIDELINES.md
    ├── DATABASE_GUIDELINES.md
    ├── TESTING_GUIDELINES.md
    ├── DEPLOYMENT_GUIDELINES.md
    ├── DOCKER_GUIDELINES.md
    ├── KUBERNETES_GUIDELINES.md
    ├── SHELL_SCRIPT_GUIDELINES.md
    ├── SECURITY_GUIDELINES.md
    ├── OBSERVABILITY_GUIDELINES.md
    ├── GIT_AND_PR_GUIDELINES.md
    ├── RELEASE_GUIDELINES.md
    └── DOCUMENTATION_GUIDELINES.md
```

Frontend package runbooks live with the code they govern:

```text
frontend/packages/base/RUNBOOK.md
frontend/packages/business/RUNBOOK.md
frontend/packages/components/RUNBOOK.md
frontend/packages/shell/RUNBOOK.md
frontend/packages/domain/<domain>/RUNBOOK.md
```

Root `CONTRIBUTING.md` should link to the frontend guideline and package runbooks.

Runbooks are not a substitute for code comments or Storybook. They document package ownership, extension rules, public usage, development commands and failure/recovery guidance.

---

## 109. GitHub Repository Guardrails

Recommended repository files:

```text
.github/
├── CODEOWNERS
├── pull_request_template.md
├── ISSUE_TEMPLATE/
└── workflows/            # only where GitHub automation is actually used
CONTRIBUTING.md
SECURITY.md
CODE_OF_CONDUCT.md
README.md
```

Woodpecker remains the default/reference CI. GitHub metadata is for collaboration and repository governance, not a requirement to use GitHub Actions.

---

## 110. Frontend PR Checklist

Recommended PR template block:

```text
Frontend Architecture
- [ ] SRP preserved
- [ ] OCP extension point considered
- [ ] No conceptual duplication
- [ ] Correct package chosen: base / business / components / domain / shell
- [ ] No switch for behavior selection
- [ ] No avoidable if/else chain
- [ ] No forbidden deep import
- [ ] Cross-domain imports use only domain root public API
- [ ] Domains do not import `packages/shell`
- [ ] Shell imports domains only through root public APIs
- [ ] Global provider/store/route composition lives in shell, not app/domains
- [ ] Shell contains no product business logic
- [ ] No circular dependency
- [ ] index.ts public boundary respected
- [ ] Existing authoritative types reused
- [ ] API calls exist only in domain/business services
- [ ] No component-level service was introduced
- [ ] Backend entity reader reused/added where appropriate
- [ ] Lodash imports are function-level (lodash/<fn>)
- [ ] Presentational components contain no business orchestration
- [ ] Stateful behavior extracted into hooks/container
- [ ] Pure functions have tests
- [ ] Config-driven rendering used where appropriate
- [ ] Shared component reuse considered
- [ ] New public shared visual component has Storybook stories
- [ ] Shared package/shell/domain runbook updated when contract changed
- [ ] Mock entry point updated when a public domain contract needs test support
- [ ] Production code does not import domain mocks
- [ ] Accessibility checked
```

---

## 111. Final Mandatory Rules

1. Everything meaningful is a unit.
2. Each unit has one clear responsibility.
3. Frontend code is classified into `base`, `business`, `components`, `domain`, or `shell`.
4. `base` contains only pure generic technical code and is the dependency floor.
5. `business` contains reusable shared business logic/services, never domain UI.
6. `components` contains reusable platform UI and never imports a product domain or shell.
7. Domains live under `packages/domain/<name>`.
8. `packages/shell` is the application-composition root for the frontend.
9. `frontend/app` remains a thin executable/bootstrap entry point and mounts the shell.
10. The shell initializes truly global providers once.
11. The shell owns global Redux store/provider composition.
12. Domains own their reducers/actions/selectors/domain middleware and expose only registration metadata required by shell.
13. The shell aggregates domain routes/navigation/state registrations through domain root public APIs.
14. The shell must never deep-import domain internals.
15. Domains, `business`, `components`, and `base` must never import `packages/shell`.
16. The shell contains no product business logic.
17. The shell does not directly call domain/backend APIs.
18. The shell owns generic application startup skeleton/fallback behavior; domain-specific skeletons remain in domains.
19. Every public unit exposes an intentional `index.ts`.
20. Consumers use the public boundary; forbidden deep imports are blocked.
21. Cross-domain imports use only `@domain/<name>`.
22. Domains expose only what they intentionally want other domains/shell to consume.
23. Domain mocks are exposed separately through `@domain/<name>/mock`.
24. Production code must never import domain mocks.
25. Pages live at domain level.
26. API services live at domain or business level, never inside component folders or shell.
27. Backend entity/response readers live with their service/entity contract.
28. Readers are pure typed accessors and contain no business decisions.
29. Stateful pages/features use container + presentational separation.
30. Presentational components contain minimal business logic.
31. Hooks encapsulate focused reusable stateful behavior.
32. Helpers own pure domain/use-case transformations.
33. Utils are placed at the nearest correct level: base, business, or domain.
34. Types stay adjacent to the unit they describe.
35. Reuse authoritative types with indexed access, `Pick`, `Omit`, etc.
36. Reusable functions have explicit callable types.
37. Constants use `SCREAMING_SNAKE_CASE`.
38. Constant maps use `as const`; runtime immutability may use `Object.freeze`.
39. Derive constant union types instead of manually duplicating them.
40. Lodash imports are function-specific, e.g. `import _noop from 'lodash/noop'`.
41. Full-package Lodash imports are forbidden.
42. `switch` is restricted for behavior selection.
43. Avoid `if / else if` chains when strategy/config maps fit.
44. Guard clauses are allowed and encouraged.
45. Configuration-driven UI is preferred over duplicated markup.
46. Use typed renderer/strategy/factory registries.
47. Repeated UI patterns become `packages/components` primitives.
48. Public shared visual components require Storybook stories.
49. `base`, `business`, `components`, and `shell` require package runbooks.
50. Pure functions require adjacent Jest tests.
51. Do not create God Components.
52. Do not replace God Components with God Hooks.
53. Capability-based behavior is preferred over vendor-name checks.
54. Circular dependencies are forbidden.
55. Strict TypeScript is mandatory.
56. `any` is forbidden except documented exceptional boundaries.
57. Public contracts are stable once consumed.
58. Promote reusable units deliberately and preserve compatibility when necessary.
59. Enforce architecture rules automatically wherever practical.
60. SRP, OCP, and DRY are mandatory review criteria.
61. Names must make intent obvious.
62. Extension should happen through registration/configuration rather than giant conditional files.
63. A new contributor should understand the architecture by inspecting the shared package runbooks, shell runbook, and one reference domain.

---

## 112. Guiding Principle

The frontend should continuously move toward:

```text
Pure Base Primitives
        +
Shared Business Contracts
        +
Reusable Platform Components
        +
Isolated Product Domains
        +
Application Composition Shell
        +
Explicit Public Boundaries
        +
Typed Services + Readers
        +
Configuration / Strategy / Registry
        +
Stories + Tests + Runbooks
```

The dependency model should remain easy to explain:

```text
app
 ↓
shell
 ├──────────────→ components ──→ base
 │
 └──────────────→ domain
                    ├────────→ business ──→ base
                    ├────────→ components ──→ base
                    └────────→ base
```

The key isolation rule is:

```text
domain  ✗ shell
business ✗ shell
components ✗ shell

shell → domain ROOT PUBLIC API only
app   → shell
```

> **Target result:** clean files, predictable package ownership, deliberate public APIs, isolated domains, one explicit application composition root, reusable backend entity contracts, reusable UI, typed readers, realistic mocks, discoverable stories/runbooks, minimal branching, and a frontend that remains maintainable as Perihellion grows.

---
