# Regression Testing Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Applies to:** Browser-driven product regression testing for Perihellion under the canonical `unit-test → sanity → regression` model  
> **Primary framework:** Playwright + TypeScript  
> **Primary purpose:** Protect supported user journeys, domain integrations, cross-module behavior, and previously working functionality after sanity has established baseline platform health  
> **Core principles:** deterministic tests, isolated state, user-visible behavior, stable selectors, real integration paths, reproducibility, debuggability, and fast feedback

---

# 1. Objective

Regression tests verify that previously supported product behavior continues to work after changes.

Perihellion has only three top-level test suites:

```text
UNIT TEST
    ↓
SANITY
    ↓
REGRESSION
```

Regression is not a replacement for unit-test or sanity.

Use:

```text
UNIT TEST
    isolated implementation/business behavior

SANITY
    fast assembled-system and real-infrastructure confidence

REGRESSION
    supported user workflows through stable public behavior
```

There is no separate top-level `smoke`, `critical`, `functional`, or `E2E` suite.

Regression is the product-level browser suite and should remain smaller than unit coverage while protecting important workflows.

---

# 2. Core Philosophy

A regression test should answer:

> Can a real user still complete this supported product workflow after our change?

Tests should verify behavior through stable public interfaces.

Preferred:

```text
Browser
  ↓
Perihellion UI
  ↓
Public backend APIs
  ↓
Real application services
  ↓
Real database/infrastructure where practical
```

Avoid tests that know internal React implementation details, Redux structure, private domain files, or backend database internals.

---

# 3. Mandatory Rules

0. Regression runs within the canonical `unit-test → sanity → regression` model; sanity owns fast startup/health/real-dependency confidence.


1. Playwright is the standard browser regression framework.
2. Regression tests are written in TypeScript.
3. Every test must be independently runnable unless explicitly marked as a serial workflow.
4. Tests must not depend on execution order.
5. Tests must not depend on data created by another test.
6. Tests must not use arbitrary sleeps.
7. `page.waitForTimeout(...)` is forbidden in committed regression tests except a documented debugging-only case.
8. Prefer Playwright web-first assertions.
9. Prefer accessible/user-facing locators.
10. CSS selectors and XPath are a last resort.
11. Test IDs are allowed only when semantic selectors cannot provide a stable contract.
12. Regression tests must test user-visible behavior, not implementation details.
13. Core product regression tests should use real backend behavior.
14. Mocking is allowed only at intentional external boundaries or in explicitly classified mocked suites.
15. Test data must be deterministic and isolated.
16. Tests must clean up owned data or use disposable environments.
17. Credentials and secrets must never be committed.
18. Authentication setup must be reusable.
19. Test retries must not hide flaky tests.
20. A consistently flaky test is a defect.
21. Regression failures must produce enough artifacts to diagnose the problem.
22. Tests must support local execution and CI execution.
23. CI should fail when a non-quarantined regression test fails.
24. Critical regression flows must be runnable separately as a fast gate.
25. Domain tests must be organized by product/domain responsibility.
26. Regression tests must not deep-import frontend domain internals.
27. Regression tests must not query application databases directly to assert user-facing behavior.
28. Setup through supported APIs is preferred over slow UI setup unless the setup journey itself is under test.
29. Assertions must verify meaningful business outcomes, not only that an element exists.
30. The suite must remain maintainable as the product grows.

---

# 4. Repository Location

Regression tests belong at repository level because they validate the assembled product, not only the frontend package.

Recommended:

```text
perihellion/
├── frontend/
├── backend/
├── testing/
│   └── regression/
│       ├── playwright.config.ts
│       ├── tests/
│       ├── fixtures/
│       ├── auth/
│       ├── domains/
│       ├── components/
│       ├── helpers/
│       ├── builders/
│       ├── readers/
│       ├── constants/
│       ├── types/
│       ├── data/
│       ├── setup/
│       ├── teardown/
│       ├── reporters/
│       └── README.md
│
└── docs/
    └── engineering/
        └── REGRESSION_TESTING.md
```

Do not place product-wide regression tests inside an individual React component.

---

# 5. Recommended Regression Structure

```text
testing/regression/
├── playwright.config.ts
│
├── domains/
│   ├── authentication/
│   ├── catalog/
│   ├── delivery/
│   ├── testing/
│   ├── workbench/
│   ├── governance/
│   └── ...
│
├── fixtures/
├── auth/
├── components/
├── helpers/
├── builders/
├── readers/
├── constants/
├── types/
├── data/
├── setup/
├── teardown/
├── reporters/
└── README.md
```

Domain tests live inside each domain's own `tests/` folder.

Do not create top-level:

```text
smoke/
critical/
e2e/
```

Risk/criticality can be metadata, not a separate suite or directory.

Do not create empty folders only to satisfy the blueprint.

---

# 6. Domain Test Organization

Regression tests should reflect product domains.

Example:

```text
testing/regression/domains/catalog/
├── pages/
│   ├── ApplicationListPage/
│   └── ApplicationDetailsPage/
│
├── flows/
│   ├── createApplication.flow.ts
│   ├── updateApplication.flow.ts
│   └── deleteApplication.flow.ts
│
├── assertions/
├── builders/
├── types/
├── tests/
│   ├── application-create.spec.ts
│   ├── application-update.spec.ts
│   └── application-delete.spec.ts
│
└── index.ts
```

Tests should consume the domain regression public boundary rather than deep-importing arbitrary files.

---

# 7. Public Regression Test Boundary

Each reusable regression domain can expose its public helpers through `index.ts`.

Example:

```ts
export { ApplicationListPage } from './pages/ApplicationListPage';
export { createApplication } from './flows/createApplication.flow';

export type {
  CreateApplicationFlowParams,
} from './types';
```

Consumers should use:

```ts
import {
  ApplicationListPage,
  createApplication,
} from '@regression/domains/catalog';
```

Do not deep-import another regression domain's internals.

---

# 8. Regression Classification

All tests in this document belong to the single top-level suite:

```text
REGRESSION
```

Use metadata for secondary characteristics such as:

```text
domain
risk
destructive
read-only/write
slow
migration
provider capability
```

Do not create separate regression subclasses that become additional top-level suites.

Example:

```ts
test(
  '[REG-CATALOG-001] developer can create an application',
  {
    tag: ['@regression', '@catalog', '@high-risk'],
  },
  async ({ page }) => {
    // ...
  },
);
```

---

# 9. Recommended Tags

Use consistent tags.

```text
@regression

@high-risk
@destructive
@read-only
@write
@slow

@auth
@catalog
@delivery
@testing
@workbench
@migration
@governance
@platform
```

Do not use `@sanity` or `@high-risk` to create hidden parallel suites.

Criticality is represented as risk metadata such as `@high-risk`.

Do not create uncontrolled one-off tags.

Document supported tags in the regression README.

---

# 10. Test IDs

Important regression cases should have stable IDs.

Example:

```text
REG-AUTH-001
REG-CATALOG-001
REG-DELIVERY-014
REG-WORKBENCH-007
```

Example title:

```ts
test(
  '[REG-CATALOG-001] user can create an application',
  async ({ page }) => {
    // ...
  },
);
```

Stable IDs help with:

- defect references
- release evidence
- test management
- historical results
- requirements mapping
- regression reports

---

# 11. What Belongs in Regression

Good candidates:

- login/logout
- organization/project/application creation
- permissions
- navigation
- domain CRUD journeys
- cross-domain workflows
- One-Go missions
- deployment flows
- approval flows
- Workbench context switching
- environment creation
- test execution
- rollback
- provider setup
- configuration changes
- migration workflows
- critical error paths
- browser refresh/resume behavior
- role-based behavior

Do not turn every small UI condition into a browser regression test.

---

# 12. Locators

Use Playwright locators that represent how a user finds the element.

Priority:

```text
1. getByRole
2. getByLabel
3. getByPlaceholder
4. getByText when semantically stable
5. getByTestId
6. CSS selector only when unavoidable
7. XPath only as an exceptional last resort
```

Preferred:

```ts
await page
  .getByRole('button', { name: 'Create application' })
  .click();
```

Avoid:

```ts
await page.locator(
  '#root > div:nth-child(2) > div > button',
).click();
```

Avoid:

```ts
await page.locator('//div[3]/button[2]').click();
```

---

# 13. `data-testid` Policy

`data-testid` is a supported testing contract when semantic locators are insufficient.

Example:

```tsx
<div data-testid="execution-status" />
```

Test:

```ts
await expect(
  page.getByTestId('execution-status'),
).toHaveText('Succeeded');
```

Do not add a test ID to every DOM node.

Use stable semantic names.

Good:

```text
application-create-submit
execution-status
workbench-log-panel
```

Bad:

```text
button1
div2
test123
```

---

# 14. Never Use Arbitrary Sleeps

Forbidden:

```ts
await page.waitForTimeout(5000);
```

Use product state.

Examples:

```ts
await expect(
  page.getByRole('heading', {
    name: 'Application details',
  }),
).toBeVisible();
```

```ts
await expect(
  page.getByTestId('execution-status'),
).toHaveText('Succeeded');
```

```ts
await expect
  .poll(async () => getExecutionStatus(executionId))
  .toBe('SUCCEEDED');
```

Tests wait for meaningful state, not guessed time.

---

# 15. Web-First Assertions

Prefer Playwright assertions that automatically retry.

Good:

```ts
await expect(locator).toBeVisible();
await expect(locator).toHaveText('Ready');
await expect(locator).toBeEnabled();
await expect(locator).toHaveCount(3);
```

Avoid manually polling the DOM where Playwright assertions already solve the problem.

---

# 16. Assertions Must Verify Business Outcomes

Weak:

```ts
await expect(page.locator('table')).toBeVisible();
```

Better:

```ts
await expect(
  page.getByRole('row', {
    name: /crm.*production/i,
  }),
).toBeVisible();
```

A regression test should prove the expected outcome.

---

# 17. Page / Screen Objects

Use page/screen objects to encapsulate stable UI interaction.

Example:

```ts
export class ApplicationListPage {
  constructor(
    private readonly page: Page,
  ) {}

  async open(): Promise<void> {
    await this.page.goto('/applications');
  }

  async createApplication(
    name: string,
  ): Promise<void> {
    await this.page
      .getByRole('button', {
        name: 'Create application',
      })
      .click();

    await this.page
      .getByLabel('Application name')
      .fill(name);

    await this.page
      .getByRole('button', {
        name: 'Create',
      })
      .click();
  }
}
```

Page objects must not become God Objects.

Split by coherent screen/component responsibility.

---

# 18. Component Objects

Reusable complex UI primitives can have regression component objects.

Example:

```text
components/
├── Table/
├── Form/
├── Modal/
├── Drawer/
├── EnvironmentSelector/
└── ...
```

Example:

```ts
export class ModalComponent {
  constructor(
    private readonly root: Locator,
  ) {}

  async confirm(): Promise<void> {
    await this.root
      .getByRole('button', { name: 'Confirm' })
      .click();
  }
}
```

This reduces duplicated interaction logic.

---

# 19. Flows

A flow represents a reusable business interaction composed from page/component objects.

Example:

```ts
export const createApplication:
  CreateApplicationFlow = async ({
    page,
    application,
  }) => {
    const applications =
      new ApplicationListPage(page);

    await applications.open();
    await applications.createApplication(
      application.name,
    );
  };
```

Do not put assertions for unrelated use cases into one giant flow.

---

# 20. Test Data Builders

Use builders for reusable test data.

Example:

```ts
export const buildApplication:
  BuildApplication = (overrides = {}) => ({
    name: `regression-${crypto.randomUUID()}`,
    description: 'Regression application',
    ...overrides,
  });
```

Use unique test-owned values.

Do not rely on a permanently existing record such as:

```text
test-app-1
```

unless it is explicitly immutable reference data.

---

# 21. Data Isolation

Each test must own its mutable data.

Preferred:

```text
Test
 ↓
unique organization/project/application/resource
 ↓
test operates
 ↓
cleanup or disposable environment teardown
```

Parallel tests must not mutate the same entity.

Use:

- generated UUIDs
- worker-specific prefixes
- run IDs
- tenant/project isolation
- disposable environments

---

# 22. Setup Through APIs

Use supported API setup when UI setup is not the behavior being tested.

Example:

```text
API:
create organization
create project
create application

Browser:
test deployment workflow
```

This makes regression tests faster and more focused.

If the test is specifically testing "Create Application through UI", perform that part through UI.

---

# 23. Never Seed Through the Database Directly

Regression tests should not insert directly into application tables.

Bad:

```text
INSERT INTO applications ...
```

Preferred:

```text
Public API
Admin API
Approved test setup API
Supported product workflow
```

Direct database setup tightly couples regression tests to persistence internals.

---

# 24. Cleanup

Tests must leave the environment usable.

Preferred cleanup:

```ts
test.afterEach(async ({ api }) => {
  await cleanupOwnedResources(api);
});
```

Cleanup should be idempotent.

Cleanup errors should be reported without masking the original test failure.

Disposable preview environments may be destroyed as a whole instead.

---

# 25. Authentication

Do not login through the UI before every test.

Use reusable authentication state.

Recommended:

```text
auth/
├── setup/
│   ├── admin.setup.ts
│   ├── developer.setup.ts
│   └── viewer.setup.ts
│
└── states/
    ├── admin.json
    ├── developer.json
    └── viewer.json
```

Generated auth-state files containing credentials/tokens must not be committed.

---

# 26. Playwright Projects for Roles

Use projects where roles/browser configurations differ.

Example:

```ts
projects: [
  {
    name: 'admin',
    use: {
      storageState: '.auth/admin.json',
    },
  },
  {
    name: 'developer',
    use: {
      storageState: '.auth/developer.json',
    },
  },
];
```

Do not duplicate the same test logic simply to change credentials when a parameterized/role-aware approach is clearer.

---

# 27. Secrets

Secrets must come from:

- CI secret storage
- environment variables
- approved secret providers
- ephemeral credentials

Never commit:

```text
passwords
tokens
API keys
private keys
cookies
storage-state secrets
```

---

# 28. Mocking Philosophy

Regression tests should prefer real product integration.

Default:

```text
Frontend
   ↓
Real backend
   ↓
Real PostgreSQL
   ↓
Reference providers / test infrastructure
```

Mocks are appropriate for:

- unavailable third-party SaaS
- expensive external providers
- destructive external operations
- deterministic simulation of rare external failures
- explicitly isolated browser behavior tests

Do not mock the core backend in a test that claims to be full-product regression.

---

# 29. Frontend Domain Mock Packages

The frontend's:

```text
packages/domain/<domain>/mock
```

is primarily for:

- component tests
- Storybook
- frontend development
- explicitly mocked browser suites

Full regression tests should not depend on frontend-private mock implementation details.

If a regression suite intentionally uses a mock provider, that must be obvious from the suite name/tag.

---

# 30. Network Interception

Playwright `page.route(...)` may be used when the test explicitly owns the mocked boundary.

Example:

```ts
await page.route(
  '**/external-provider/**',
  async (route) => {
    await route.fulfill({
      status: 200,
      json: mockedProviderResponse,
    });
  },
);
```

Do not silently intercept core product APIs in normal regression suites.

---

# 31. Response Readers

Backend/API entities used in setup or assertions may use typed readers.

Example:

```ts
export const applicationResponseReader =
  Object.freeze({
    applicationId:
      createPropertyReader<
        ApplicationResponse,
        'applicationId'
      >('applicationId'),

    name:
      createPropertyReader<
        ApplicationResponse,
        'name'
      >('name'),
  });
```

This keeps API entity access consistent with frontend service-reader conventions.

---

# 32. Lodash Import Rule

When Lodash is used, import only the required module.

Correct:

```ts
import _property from 'lodash/property';
import _groupBy from 'lodash/groupBy';
import _uniqBy from 'lodash/uniqBy';
import _noop from 'lodash/noop';
```

Avoid:

```ts
import _ from 'lodash';
```

Avoid:

```ts
import * as lodash from 'lodash';
```

---

# 33. Fixtures

Use Playwright fixtures for reusable test dependencies.

Good fixture responsibilities:

- authenticated page
- API client
- test organization
- project
- application
- browser context
- role
- cleanup registry

Example:

```ts
export const test = base.extend<
  RegressionFixtures
>({
  application: async ({ api }, use) => {
    const application =
      await createTestApplication(api);

    await use(application);

    await deleteTestApplication(
      api,
      application.id,
    );
  },
});
```

Fixtures should be composable and focused.

---

# 34. Worker Fixtures

Use worker-scoped fixtures only when expensive resources can safely be shared.

Examples:

- immutable reference dataset
- dedicated worker tenant
- reusable authenticated account

Never share mutable entities across parallel tests accidentally.

---

# 35. Global Setup

Global setup should contain only environment-level concerns.

Examples:

- authenticate reusable roles
- verify target environment
- create immutable test bootstrap data
- validate required feature capabilities

Do not build the entire test state globally.

---

# 36. Environment Safety

Every regression execution must know its target environment.

Recommended environment values:

```text
local
preview
test
staging
production
```

Destructive tests must refuse to run against production.

Example:

```ts
if (
  targetEnvironment === 'production' &&
  isDestructiveSuite
) {
  throw new Error(
    'Destructive regression tests cannot run in production',
  );
}
```

Do not rely only on developer memory.

---

# 37. Production Regression Policy

Production should normally run only explicitly approved, restricted read-only regression checks.

Examples:

- home page loads
- login works
- approved high-risk read-only navigation
- health-visible workflows
- safe synthetic tenant operations where explicitly designed

Do not run normal destructive regression against real customer data.

---

# 38. Browser Matrix

Default pull-request gate should use one primary browser unless cross-browser behavior is relevant.

Recommended:

```text
PR:
  Chromium

Nightly / release:
  Chromium
  Firefox
  WebKit
```

This keeps PR feedback fast while still protecting cross-browser compatibility.

---

# 39. Responsive Coverage

High-risk responsive layouts should include explicit projects/viewports where the product supports them.

Do not multiply every regression case across every viewport.

Choose representative high-risk journeys.

---

# 40. Parallelism

Tests should be parallel by default when isolation permits.

A test requiring serial execution must explain why.

Avoid:

```ts
test.describe.configure({
  mode: 'serial',
});
```

as a workaround for shared state or poor isolation.

---

# 41. Retries

Recommended:

```text
local: 0
CI:    1 or 2
```

Retries exist to capture useful diagnostics for occasional infrastructure/browser noise.

They must not make known flaky tests acceptable.

Track first-run failures separately.

---

# 42. Flaky Test Policy

A flaky test is a defect.

When a test becomes flaky:

```text
1. identify cause
2. fix synchronization/data/isolation
3. add diagnostics
4. rerun repeatedly
5. only quarantine temporarily when necessary
```

Do not simply increase retry count.

---

# 43. Quarantine Policy

A quarantined regression test must have:

- explicit annotation/tag
- linked issue
- owner
- reason
- date
- expected fix/removal date

Example:

```ts
test.fixme(
  true,
  'PERI-1234: intermittent provider callback race',
);
```

Quarantine must be temporary and visible.

---

# 44. Traces

Use Playwright traces for failure investigation.

Recommended CI configuration:

```ts
trace: 'on-first-retry',
```

Trace should capture:

- DOM snapshots
- actions
- network
- console
- screenshots

Store traces as CI artifacts.

---

# 45. Screenshots

Recommended:

```ts
screenshot: 'only-on-failure',
```

Do not generate screenshots for every passing step unless a compliance/evidence workflow explicitly requires them.

---

# 46. Video

Recommended:

```ts
video: 'retain-on-failure',
```

Video is diagnostic evidence, not a substitute for meaningful assertions.

---

# 47. Example Playwright Configuration

```ts
import {
  defineConfig,
  devices,
} from '@playwright/test';

export default defineConfig({
  testDir: './tests',

  fullyParallel: true,

  forbidOnly: Boolean(process.env.CI),

  retries: process.env.CI ? 1 : 0,

  workers: process.env.CI
    ? undefined
    : undefined,

  reporter: [
    ['list'],
    ['html', {
      outputFolder: 'playwright-report',
      open: 'never',
    }],
  ],

  use: {
    baseURL:
      process.env.PLAYWRIGHT_BASE_URL ??
      'http://localhost:8080',

    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
      },
    },
  ],
});
```

Configuration should be centralized.

Individual tests must not redefine global behavior without a real need.

---

# 48. Timeouts

Use meaningful timeout levels.

Recommended conceptual separation:

```text
expect timeout
action timeout
navigation timeout
test timeout
long-running workflow timeout
```

Do not solve slowness by setting every timeout to several minutes.

Long-running Perihellion missions should use explicit polling/assertion helpers.

---

# 49. Long-Running Product Operations

For operations such as:

- deployments
- test runs
- migrations
- mission executions
- resource provisioning

do not keep arbitrary UI sleeps.

Use status polling:

```ts
await expect
  .poll(
    async () =>
      executionReader.status(
        await getExecution(executionId),
      ),
    {
      timeout: EXECUTION_TIMEOUT_MS,
    },
  )
  .toBe(EXECUTION_STATUS.SUCCEEDED);
```

Then assert the UI reflects the final state where relevant.

---

# 50. Polling Helpers

Centralize domain polling behavior.

Example:

```text
helpers/polling/
├── waitForExecution.ts
├── waitForDeployment.ts
├── waitForResource.ts
├── types/
└── index.ts
```

Do not repeat retry loops in test files.

---

# 51. API Helpers

Test setup API calls belong in dedicated helpers/clients.

```text
helpers/api/
├── application.api.ts
├── environment.api.ts
├── execution.api.ts
├── types/
└── index.ts
```

Use the product's public API.

Do not call backend private Go packages from Playwright tests.

---

# 52. Test Constants

Use stable constants.

```ts
export const TEST_ROLE = {
  ADMIN: 'ADMIN',
  DEVELOPER: 'DEVELOPER',
  VIEWER: 'VIEWER',
} as const;
```

Derive types:

```ts
export type TestRole =
  (typeof TEST_ROLE)[keyof typeof TEST_ROLE];
```

Avoid magic strings in test logic.

---

# 53. Naming

Test filenames:

```text
<feature>.<behavior>.spec.ts
```

Examples:

```text
application.create.spec.ts
application.permissions.spec.ts
deployment.rollback.spec.ts
workbench.logs.spec.ts
mission.resume.spec.ts
```

Helpers:

```text
application.builder.ts
deployment.status.reader.ts
mission.polling.helper.ts
```

Names should explain intent.

---

# 54. Test Structure

Use Arrange / Act / Assert conceptually.

Example:

```ts
test(
  '[REG-CATALOG-001] user can create an application',
  async ({ page, api }) => {
    // Arrange
    const application =
      buildApplication();

    // Act
    const applications =
      new ApplicationListPage(page);

    await applications.open();

    await applications.createApplication(
      application.name,
    );

    // Assert
    await expect(
      page.getByRole('heading', {
        name: application.name,
      }),
    ).toBeVisible();
  },
);
```

Comments are optional when the structure is already obvious.

---

# 55. One Test — One Primary Behavior

A test should have one primary reason to fail.

Avoid:

```text
create application
edit application
add environment
delete environment
deploy workload
delete application
```

all in one test unless the full lifecycle itself is the intended regression journey.

Split independent behaviors.

---

# 56. Full Journey Tests

A limited number of broader journey tests are valuable.

Examples:

```text
Developer onboarding
Create application → environment → run validation

Delivery
Build → test → deploy → verify → promote

Incident reproduction
Open application → select revision → reproduce → inspect logs

Migration
Assess → map → dry-run → validate
```

These should complement focused tests, not replace them.

---

# 57. Negative Cases

Regression should include high-value negative cases.

Examples:

- unauthorized user cannot deploy
- invalid form cannot submit
- required approval blocks execution
- failed test prevents promotion
- unsupported capability is explained
- expired session redirects to login

Do not duplicate every validation branch already covered by unit tests.

---

# 58. Permissions Regression

For security-sensitive operations, test both:

```text
authorized role → allowed
unauthorized role → blocked
```

The frontend check is not enough.

Where practical, verify the operation does not succeed through the backend even if UI behavior changes.

---

# 59. Accessibility in Regression

High-risk flows should include focused accessibility checks where practical.

Regression tests should also use semantic locators because this naturally encourages accessible UI.

Dedicated accessibility scanning may be implemented separately.

---

# 60. Console Errors

Unexpected browser errors should be detectable.

Use a reusable fixture/helper to capture:

- uncaught exceptions
- severe console errors
- failed high-risk network requests

Allow documented known exceptions only through explicit configuration.

---

# 61. Network Failures

When intentionally testing network failure behavior, make the scenario explicit.

Example:

```text
@regression @network-failure
```

Verify:

- user-visible error
- retry behavior
- no stale success state
- no duplicate mutation

---

# 62. Idempotency

Tests that retry product mutations should verify the product's idempotency semantics where relevant.

Examples:

- create operation with idempotency key
- retry deployment request
- duplicate webhook callback
- repeated approval action

---

# 63. Visual Regression

Visual screenshot comparison is a separate concern from behavioral regression.

Use visual regression selectively for:

- design-system primitives
- high-risk layouts
- charts
- complex Workbench panels
- responsive layouts

Do not make every behavioral test depend on screenshots.

Visual baselines must be intentionally reviewed.

---

# 64. Storybook vs Regression

Shared components require Storybook stories.

Storybook validates isolated component states.

Playwright regression validates assembled product behavior.

```text
Storybook:
Component in isolation

Playwright Regression:
Real user journey in assembled product
```

Do not duplicate the entire Storybook catalog as browser regression tests.

---

# 65. CI Execution Order

Recommended product test sequence:

```text
1. format/lint/static
2. unit-test
3. start required infrastructure
4. run migrations
5. start backend
6. wait for readiness
7. start frontend
8. wait for readiness
9. run sanity
10. run regression as required by the pipeline
11. publish reports/artifacts
12. cleanup
```

If sanity fails, regression may be skipped because baseline platform health is not established.

---

# 66. Pull Request Gate

Recommended PR gate:

```text
unit-test
sanity
targeted/high-risk regression when the change affects user behavior
```

Regression selection may use domain/risk metadata when dependency mapping is reliable.

Do not add a separate smoke or critical gate.

---

# 67. Main Branch / Nightly

Recommended:

```text
unit-test
sanity
all regression
Chromium
Firefox
WebKit where supported
```

Run repeated/flakiness detection for high-risk regression cases periodically.

---

# 68. Release Gate

Before release:

```text
unit-test
sanity
all required regression
required cross-browser projects
upgrade regression
installation regression
permissions regression
```

Release failures must be triaged before shipping.

---

# 69. Sharding

Large suites should support Playwright sharding.

Example:

```bash
npx playwright test --shard=1/4
```

Tests must remain independent so sharding does not change correctness.

---

# 70. CI Artifacts

On failure retain:

```text
HTML report
trace
screenshot
video
browser console
relevant product logs
test metadata
run ID
commit SHA
environment
```

A developer should be able to diagnose most failures without rerunning locally first.

---

# 71. Correlation IDs

Where the product exposes request/correlation IDs, regression helpers should capture them for failed operations.

This allows:

```text
Playwright failure
      ↓
correlation ID
      ↓
backend logs/traces
```

This is especially important for Perihellion's Workbench and Mission Runtime.

---

# 72. Test Evidence

For high-risk/release regression, record:

- test ID
- product version
- commit SHA
- environment
- browser/project
- start/end time
- result
- failure reason
- artifacts
- retry count

This can later feed Perihellion's own evidence model.

---

# 73. No Hidden State

Tests must not depend on:

- previous browser session
- browser cache from a different test
- manual environment setup
- an engineer's machine state
- previously executed test order
- undocumented feature flags

Everything required must be created, configured, or declared.

---

# 74. Environment Capability Checks

If a regression test requires a capability such as:

```text
Kubernetes runtime
database branching
CI provider
migration connector
```

the test should fail or skip explicitly based on declared environment capability.

Never silently behave differently because a provider happens to be missing.

---

# 75. Skip Policy

A skipped test must have a reason.

Bad:

```ts
test.skip(true);
```

Good:

```ts
test.skip(
  !capabilities.has(CAPABILITY.DB_BRANCHING),
  'Environment does not provide DB branching capability',
);
```

Do not use skip to hide defects.

---

# 76. Test Ownership

Every regression domain should have an owner.

Example metadata:

```text
domain: delivery
owner: delivery-platform
risk: high
```

Failed tests should route to a responsible team/domain.

---

# 77. Runbook

`testing/regression/README.md` or `RUNBOOK.md` must explain:

- how to install Playwright
- how to run locally
- how to choose environment
- required environment variables
- how authentication works
- how to run tags
- how to run one domain
- how to debug
- how to open trace viewer
- how to update visual snapshots if used
- how cleanup works
- how quarantine works
- CI commands
- common failure modes

---

# 78. Local Commands

Canonical repository commands:

```bash
make test-regression
```

Useful focused execution may be exposed without creating new suites, for example:

```bash
make test-regression-domain DOMAIN=catalog
```

Equivalent package commands may be provided.

Examples:

```bash
npx playwright test
```

```bash
npx playwright test testing/regression/domains/catalog
```

```bash
npx playwright test --grep @high-risk
```

```bash
npx playwright test --ui
```

Do not create separate regression Make targets for smoke or critical classifications.

---

# 79. Debugging

Preferred debugging tools:

```text
Playwright UI mode
trace viewer
headed mode
inspector
screenshots
video
browser console
network logs
backend correlation IDs
```

Do not commit debug sleeps or `test.only`.

CI must use:

```ts
forbidOnly: true
```

---

# 80. Regression Coverage Model

Coverage is not measured only by number of tests.

Track coverage across:

```text
Product domain
High-risk user journey
Role
Environment
Capability/provider
Browser
Failure mode
Release risk
```

A thousand low-value tests are not better than a smaller, well-selected regression suite.

---

# 81. Risk-Based Regression

Priority should follow risk.

High priority:

- authentication
- authorization
- irreversible actions
- deployments
- migration
- approvals
- execution engine
- environment lifecycle
- data integrity
- provider switching
- upgrade compatibility

Lower-risk cosmetic behavior belongs primarily in component/visual testing.

---

# 82. Regression for Bug Fixes

Every meaningful escaped bug should normally add regression protection at the lowest effective layer.

If only browser-level product behavior can reproduce the failure, add a Playwright regression test.

The new test should fail before the fix and pass after it where practical.

Use the defect ID in metadata/title when useful.

---

# 83. Provider Regression

For provider-backed capabilities, separate:

```text
Capability regression
Provider conformance
Provider-specific regression
```

A product workflow should primarily assert capability semantics.

Provider-specific suites verify adapter behavior.

Do not copy the full product regression suite once per provider unless required.

---

# 84. Reference Provider Matrix

Perihellion's reference stack should have a defined regression matrix.

Example:

```text
CI          Woodpecker
Cache       Valkey
Events      NATS JetStream
Storage     MinIO/S3
Secrets     OpenBao
Runtime     local/Docker initially
Database    PostgreSQL
```

As additional providers become supported, conformance suites protect their contracts.

---

# 85. Upgrade Regression

Self-hosted product releases require regression for:

```text
fresh install
upgrade from supported previous version
database migration
configuration compatibility
rollback where supported
existing user/project/application data
```

Upgrade regression must not be postponed until GA.

---

# 86. Clone-and-Run Regression

Because Perihellion is open source, every release should prove the documented quick-start path.

Conceptually:

```text
clean environment
    ↓
clone repository / release bundle
    ↓
bootstrap
    ↓
start platform
    ↓
complete setup
    ↓
login
    ↓
create first project/application
    ↓
run a small high-risk workflow
```

This is a release-gating regression path.

---

# 87. Installation Profiles

Regression should eventually cover:

```text
Lite
Team
Enterprise reference
Air-gapped reference
```

Phase 0/early development may begin with Lite.

Do not claim support for a profile that is not tested.

---

# 88. API/UI Contract Regression

When the frontend depends on a backend contract, protect it through:

- backend public-contract validation
- generated/typed frontend contracts where applicable
- regression journey
- compatibility checks

Do not use browser regression as the only API compatibility mechanism.

---

# 89. Feature Flags

Regression tests must explicitly declare relevant feature configuration.

Do not assume a developer's local feature flags.

If the feature is off:

- skip with explicit capability/feature reason, or
- use a dedicated project/environment where enabled

---

# 90. Clock / Time-Sensitive Tests

Avoid tests that depend on wall-clock timing.

For expiration/countdown/time-window behavior:

- use controlled test configuration
- use backend test facilities where supported
- assert with reasonable bounds
- avoid exact millisecond assumptions

---

# 91. File Upload / Download

Use deterministic fixtures under:

```text
testing/regression/data/files/
```

Keep fixtures small unless testing size behavior.

Generated downloaded files should go to test-run temporary storage, not the repository.

---

# 92. Cleanup Registry

Complex tests should register created resources centrally.

Concept:

```ts
cleanup.register({
  type: 'application',
  id: applicationId,
});
```

Teardown then performs cleanup in safe reverse order.

This is preferable to scattered cleanup logic in every test.

---

# 93. Test Metadata

Recommended metadata model:

```ts
interface RegressionTestMetadata {
  id: string;
  domain: string;
  owner: string;
  risk: 'low' | 'medium' | 'high';
  destructive: boolean;
}
```

Metadata may later feed reporting/evidence.

---

# 94. Reporter

Use Playwright's standard reporters initially.

Recommended:

```text
list
html
junit in CI when needed
```

A custom reporter should be introduced only when Perihellion needs structured evidence ingestion.

---

# 95. Test Results Must Be Reproducible

A failed CI test should be reproducible using:

```text
commit SHA
environment/profile
browser
seed/run ID
test ID
trace
```

If failures cannot be reproduced or diagnosed, improve the test harness.

---

# 96. Architecture Boundaries

Regression code can depend on:

```text
Playwright
public test helpers
public API clients
test builders
regression page objects
test configuration
```

Regression code should not depend on:

```text
frontend domain internals
Redux internal state
private React hooks
backend service internals
backend repositories
database tables
provider implementation internals
```

Test the product contract.

---

# 97. No DOM Implementation Assertions

Avoid assertions such as:

```text
React component class names
framework-generated IDs
internal wrapper hierarchy
Redux state shape
CSS module hashes
```

Assert user-observable behavior.

---

# 98. Avoid Over-Generalized Test Frameworks

Do not create a huge custom test DSL before repeated needs are proven.

Use Playwright directly for simple behavior.

Promote helpers only when they represent stable reusable concepts.

SRP, OCP, and DRY still apply, but avoid abstraction for abstraction's sake.

---

# 99. Regression PR Checklist

```text
Regression Testing
- [ ] High-risk behavior has appropriate regression coverage
- [ ] No arbitrary waits
- [ ] Stable semantic locators used
- [ ] Test data is isolated
- [ ] Test does not depend on execution order
- [ ] Setup uses supported API where appropriate
- [ ] Core backend is not unnecessarily mocked
- [ ] Secrets are not committed
- [ ] Assertions verify business outcomes
- [ ] Cleanup is reliable
- [ ] Test can run independently
- [ ] Test supports CI
- [ ] Failure artifacts are useful
- [ ] No test.only
- [ ] Skips have documented reasons
- [ ] Flakiness has not been hidden with retries
- [ ] Cross-domain helpers use public boundaries
```

---

# 100. Definition of Done — Regression Test

A regression test is complete when:

```text
[ ] It protects a meaningful supported behavior
[ ] It has a clear test ID when appropriate
[ ] It belongs to regression and uses supported domain/risk metadata correctly
[ ] It runs independently
[ ] It uses deterministic isolated data
[ ] It does not rely on arbitrary sleeps
[ ] It uses stable locators
[ ] It uses web-first assertions
[ ] It verifies meaningful outcome
[ ] Setup is efficient
[ ] Cleanup is safe/idempotent
[ ] Secrets are externalized
[ ] Failure artifacts are produced
[ ] It passes repeatedly
[ ] It passes in CI
[ ] It does not deep-import product internals
[ ] It does not query DB internals
[ ] It has an owner/domain
[ ] Documentation is updated when new harness behavior is introduced
```

---

# 101. Phase 0 Regression Deliverables

Phase 0 should establish the regression harness, not hundreds of browser cases.

Required Phase 0 deliverables:

```text
[ ] Playwright installed and pinned
[ ] playwright.config.ts
[ ] Chromium local/CI project
[ ] common fixture architecture
[ ] reusable authentication-state mechanism
[ ] API setup helper
[ ] test-data builder
[ ] cleanup mechanism
[ ] one page/screen object
[ ] one reusable component object
[ ] one response reader
[ ] regression tag
[ ] domain/risk metadata convention
[ ] one authenticated regression journey
[ ] one negative permission regression
[ ] trace on first retry where retry is enabled
[ ] screenshot on failure
[ ] video retained on failure according to policy
[ ] HTML report
[ ] CI execution after sanity
[ ] RUNBOOK / README
[ ] no arbitrary-wait lint/review rule
```

Suggested reference journey:

```text
Sanity passes
      ↓
Start Perihellion
      ↓
Login
      ↓
Create organization/project/application
      ↓
Open application
      ↓
Verify persisted application data through public behavior
      ↓
Logout
```

This proves the regression harness end to end without introducing smoke or critical suites.

---

# 102. Suggested Root Engineering Documentation

```text
docs/engineering/
├── FRONTEND_CODE_GUIDELINES.md
├── BACKEND_CODE_GUIDELINES.md
├── REGRESSION_TESTING.md
├── API_CONTRACT_GUIDELINES.md
├── PROVIDER_SDK_GUIDELINES.md
├── DATABASE_GUIDELINES.md
├── TESTING_GUIDELINES.md
├── DEPLOYMENT_GUIDELINES.md
├── SHELL_SCRIPT_GUIDELINES.md
├── SECURITY_GUIDELINES.md
└── GIT_AND_PR_GUIDELINES.md
```

`REGRESSION_TESTING.md` owns browser-level product regression rules.

`TESTING_GUIDELINES.md` defines the complete three-suite model and owns how database, contract, provider, performance, security, and resilience responsibilities map into unit-test, sanity, or regression.

---

# 103. Final Mandatory Rules

1. Playwright + TypeScript is the browser regression standard.
2. Regression is one of exactly three top-level suites: unit-test, sanity, regression.
3. Sanity owns fast startup/health/real-dependency confidence before broader regression.
4. No separate smoke, critical, functional, or E2E suite exists.
5. Test behavior, not implementation.
6. No arbitrary waits.
7. Use web-first assertions.
8. Use semantic locators first.
9. Use `data-testid` only as a deliberate stable contract.
10. Every test owns its mutable state.
11. Regression tests are independent by default.
12. Core regression uses real backend behavior through public boundaries.
13. Mocks are explicit and limited to intentional external boundaries.
14. Setup via API is preferred when setup UI is not being tested.
15. Never seed/assert application state through private database internals.
16. Authentication is reusable and secrets remain external.
17. Page objects encapsulate UI interaction but remain small.
18. Reusable UI interaction belongs in component objects.
19. Reusable business journeys belong in flows.
20. Builders create deterministic/unique test data.
21. Retry never excuses flakiness.
22. Quarantine is explicit, owned, linked, temporary, and regression-only.
23. Failed tests retain useful traces/artifacts.
24. Regression domains respect public boundaries.
25. No deep imports into frontend/backend internals.
26. Domain/risk metadata is standardized without creating new suites.
27. Release regression includes clone/install/upgrade flows as product maturity grows.
28. Provider behavior is protected through capability validation plus targeted regression, not duplicated full suites.
29. Regression evidence must be reproducible.
30. Production runs only explicitly approved safe regression checks.
31. SRP/OCP/DRY apply to test code without creating an unnecessary framework.
32. A regression suite that cannot be trusted is worse than a smaller deterministic suite.

---

# 104. Guiding Principle

The target is:

```text
Real User Behavior
       +
Stable Public Contracts
       +
Deterministic Test Data
       +
Playwright Auto-Waiting
       +
Isolated Tests
       +
Strong Diagnostics
       =
Trusted Regression Suite
```

> **Perihellion regression testing should tell us with high confidence whether a real supported engineering workflow still works — without brittle selectors, arbitrary sleeps, hidden shared state, or implementation coupling.**
