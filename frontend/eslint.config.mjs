import js from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  { ignores: ['**/dist/**', '**/node_modules/**'] },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ['**/*.{js,cjs,mjs,ts,tsx}'],
    languageOptions: {
      globals: {
        __: 'readonly',
        window: 'readonly',
        document: 'readonly',
        localStorage: 'readonly',
        process: 'readonly',
        console: 'readonly',
        fetch: 'readonly',
        AbortSignal: 'readonly',
        Buffer: 'readonly',
        URL: 'readonly',
        require: 'readonly',
        module: 'readonly',
        global: 'readonly',
      },
    },
  },
  {
    files: ['**/*.cjs'],
    rules: {
      '@typescript-eslint/no-require-imports': 'off',
    },
  },
  {
    files: ['**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/no-explicit-any': 'error',
      'no-restricted-syntax': [
        'error',
        {
          selector: 'ConditionalExpression',
          message: 'Use named helpers, typed descriptors, or JSX renderers instead of ternary expressions.',
        },
      ],
    },
  },
);
