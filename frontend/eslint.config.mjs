import js from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ['**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/no-explicit-any': 'error',
      'no-restricted-syntax': ['error', {
        selector: 'ConditionalExpression',
        message: 'Use named helpers, typed descriptors, or JSX renderers instead of ternary expressions.',
      }],
    },
  },
);
