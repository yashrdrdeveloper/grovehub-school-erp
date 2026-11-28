/* Real Chromium regression checks against the production Vite preview.
 * No authentication backend is required: test only public screens and token-missing states.
 * Installed in an isolated runner directory to avoid touching pnpm-lock.yaml.
 */
const assert = require('node:assert/strict');
const path = require('node:path');

const playwrightPackage = process.env.GH_PLAYWRIGHT_PACKAGE;
if (!playwrightPackage) {
  throw new Error('GH_PLAYWRIGHT_PACKAGE must point to the isolated Playwright installation');
}
const { chromium } = require(path.resolve(playwrightPackage));

const origin = process.env.GH_PREVIEW_URL || 'http://127.0.0.1:4173';
const mobileWidths = [
  { width: 320, height: 720 },
  { width: 375, height: 812 },
  { width: 390, height: 844 },
  { width: 720, height: 900 },
];
const publicPaths = ['/', '/about', '/admissions', '/gallery', '/contact'];
const authPaths = ['/login', '/forgot-password', '/reset-password', '/set-password'];

async function visit(browser, pathname, viewport, options = {}) {
  const context = await browser.newContext({
    viewport,
    colorScheme: options.colorScheme || 'light',
    reducedMotion: 'reduce',
  });
  if (options.savedTheme) {
    await context.addInitScript((value) => {
      localStorage.setItem('groovehub.theme', value);
    }, options.savedTheme);
  }
  const page = await context.newPage();
  const exceptions = [];
  page.on('pageerror', (error) => exceptions.push(error.message));
  const response = await page.goto(origin + pathname, { waitUntil: 'domcontentloaded' });
  assert.ok(response && response.status() < 400, pathname + ': page was not served');
  await page.locator('#root').waitFor({ state: 'visible' });
  // A visible root may still contain the async boot view. Wait for the actual
  // initialization contract instead of assuming dynamic imports finish in 120ms.
  await page.waitForFunction(() =>
    typeof window.__ === 'function' &&
    !document.querySelector('[data-testid="i18n-boot-status"]') &&
    Boolean(document.querySelector('#root > *')),
  );
  if (authPaths.includes(pathname)) {
    await page.locator('[data-testid="auth-layout-root"]').waitFor({ state: 'visible' });
  }
  assert.deepEqual(exceptions, [], pathname + ': unhandled browser exceptions');
  const i18n = await page.evaluate(() => ({
    installed: typeof window.__ === 'function',
    sample: typeof window.__ === 'function'
      ? window.__('Hello {{name}}', { name: 'RDR' })
      : null,
    pending: Boolean(document.querySelector('[data-testid="i18n-boot-status"]')),
  }));
  assert.deepEqual(
    i18n,
    { installed: true, sample: 'Hello RDR', pending: false },
    pathname + ': shell must install the global message function before routes mount',
  );
  const geometry = await page.evaluate(() => ({
    scroll: document.documentElement.scrollWidth,
    viewport: window.innerWidth,
    documentHeight: document.documentElement.scrollHeight,
    viewportHeight: window.innerHeight,
  }));
  assert.ok(
    geometry.scroll <= geometry.viewport + 1,
    pathname + ' at ' + viewport.width + 'px: horizontal overflow ' +
      JSON.stringify(geometry),
  );
  const duplicateTestIds = await page.evaluate(() => {
    const seen = new Set();
    const duplicates = [];
    for (const element of document.querySelectorAll('[data-testid]')) {
      const value = element.getAttribute('data-testid');
      if (!value) continue;
      if (seen.has(value)) duplicates.push(value);
      seen.add(value);
    }
    return duplicates;
  });
  assert.deepEqual(duplicateTestIds, [], pathname + ': duplicate data-testid in live DOM');
  return { context, page, geometry };
}

async function checkLayout(browser) {
  for (const viewport of mobileWidths) {
    for (const pathname of authPaths) {
      const { context, page } = await visit(browser, pathname, viewport);
      await page.locator('[data-testid="auth-layout-root"]').waitFor({ state: 'visible' });
      const heading = page.locator('[data-testid="auth-layout-title"]');
      assert.ok(await heading.isVisible(), pathname + ': missing accessible heading');
      assert.equal(await heading.getAttribute('id'), 'rdr-auth-title', pathname + ': heading ID changed');
      if (pathname === '/set-password') {
        assert.equal((await heading.innerText()).trim(), 'Set password');
      }
      if (pathname === '/reset-password') {
        assert.equal((await heading.innerText()).trim(), 'Reset password');
      }
      if (pathname === '/login' && viewport.height >= 800) {
        const submit = page.locator('[data-testid="auth-login-form-button-submit"]');
        const bounds = await submit.boundingBox();
        assert.ok(bounds, 'login submit button is not visible');
        assert.ok(bounds.y + bounds.height <= viewport.height + 1, 'submit button is below viewport');
      }
      await context.close();
      process.stdout.write('PASS auth layout ' + pathname + ' ' + viewport.width + 'px\n');
    }
    for (const pathname of publicPaths) {
      const { context, page } = await visit(browser, pathname, viewport);
      assert.ok((await page.locator('h1').count()) > 0, pathname + ': missing page heading');
      await context.close();
      process.stdout.write('PASS public layout ' + pathname + ' ' + viewport.width + 'px\n');
    }
  }
}

async function checkThemes(browser) {
  for (const testCase of [
    { width: 390, height: 844, savedTheme: 'light', colorScheme: 'dark', expected: 'dark' },
    { width: 390, height: 844, savedTheme: 'dark', colorScheme: 'light', expected: 'light' },
    { width: 1280, height: 800, savedTheme: 'light', colorScheme: 'dark', expected: 'light' },
    { width: 1280, height: 800, savedTheme: 'dark', colorScheme: 'light', expected: 'dark' },
  ]) {
    const { context, page } = await visit(browser, '/login', {
      width: testCase.width,
      height: testCase.height,
    }, testCase);
    const actual = await page.evaluate(() => document.documentElement.dataset.theme);
    assert.equal(actual, testCase.expected, 'system theme preference mismatch');
    await context.close();
  }
  process.stdout.write('PASS mobile system theme and desktop saved preference\n');
}

async function checkKeyboardAndLabels(browser) {
  const { context, page } = await visit(browser, '/login', { width: 1280, height: 800 });
  const email = page.locator('[data-testid="auth-login-form-input-email"]');
  const password = page.locator('[data-testid="auth-login-form-input-password"]');
  const toggle = page.locator('[data-testid="auth-login-form-password-visibility"]');
  assert.ok(await email.isVisible() && await password.isVisible());
  assert.ok((await email.getAttribute('aria-label')) || (await page.locator('label[for="' + await email.getAttribute('id') + '"]').count()), 'email has no accessible label');
  assert.equal(await password.getAttribute('type'), 'password');
  await toggle.focus();
  await page.keyboard.press('Enter');
  assert.equal(await password.getAttribute('type'), 'text');
  await page.keyboard.press('Enter');
  assert.equal(await password.getAttribute('type'), 'password');
  await context.close();
  process.stdout.write('PASS keyboard password visibility and form labels\n');
}


async function checkMobileFormsAndValidation(browser) {
  for (const viewport of [
    { width: 320, height: 720 },
    { width: 375, height: 812 },
    { width: 390, height: 844 },
  ]) {
    for (const testCase of [
      { pathname: '/login', submitId: 'auth-login-form-button-submit' },
      { pathname: '/forgot-password', submitId: 'forgot-password-submit' },
    ]) {
      const { context, page } = await visit(browser, testCase.pathname, viewport);
      const submit = page.locator('[data-testid="' + testCase.submitId + '"]');
      assert.ok(await submit.isVisible(), testCase.pathname + ': submit is hidden');
      const bounds = await submit.boundingBox();
      assert.ok(bounds, testCase.pathname + ': submit has no layout box');
      assert.ok(
        bounds.y >= 0 && bounds.y + bounds.height <= viewport.height + 1,
        testCase.pathname + ' ' + viewport.width + 'px: submit outside visible viewport',
      );
      await context.close();
    }
  }

  const { context, page } = await visit(browser, '/forgot-password', {
    width: 375,
    height: 812,
  });
  await page.locator('[data-testid="forgot-password-submit"]').click();
  const emailField = page.locator('input[type="email"]');
  assert.ok(await emailField.isVisible(), 'Forgot password: email input must be visible');
  const invalid = await emailField.getAttribute('aria-invalid');
  const alert = page.locator('[role="alert"]');
  assert.ok(invalid === 'true' || (await alert.count()) > 0, 'Forgot password: missing validation feedback');
  await context.close();
  process.stdout.write('PASS mobile auth submission reachability and form validation\n');
}

async function checkShortViewportAndZoom(browser) {
  for (const viewport of [
    { width: 320, height: 568 },
    { width: 640, height: 360 },
  ]) {
    for (const pathname of ['/login', '/forgot-password']) {
      const { context } = await visit(browser, pathname, viewport);
      await context.close();
    }
  }
  process.stdout.write('PASS short-viewport horizontal overflow regression\n');
}

(async () => {
  const browser = await chromium.launch({ headless: true });
  try {
    await checkLayout(browser);
    await checkMobileFormsAndValidation(browser);
    await checkShortViewportAndZoom(browser);
    await checkThemes(browser);
    await checkKeyboardAndLabels(browser);
  } finally {
    await browser.close();
  }
})().catch((error) => {
  process.stderr.write(error.stack + '\n');
  process.exitCode = 1;
});
