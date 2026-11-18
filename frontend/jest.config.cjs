module.exports = {
  testEnvironment: 'jsdom',
  testMatch: ['<rootDir>/packages/**/*.test.ts', '<rootDir>/packages/**/*.test.tsx'],
  transform: { '^.+\\.tsx?$': ['ts-jest', { tsconfig: 'tsconfig.json' }] },
  moduleNameMapper: {
    '^@Base/(.*)$': '<rootDir>/packages/base/$1',
    '^@Business/(.*)$': '<rootDir>/packages/business/$1',
    '^@Components/(.*)$': '<rootDir>/packages/components/$1',
    '^@Domain/(.*)$': '<rootDir>/packages/domain/$1',
    '^@Auth/(.*)$': '<rootDir>/packages/auth/$1',
    '^@Shell/(.*)$': '<rootDir>/packages/shell/$1',
  },
};
