module.exports = {
  testEnvironment: 'jsdom',
  setupFiles: ['<rootDir>/scripts/test/polyfills.cjs'],
  testMatch: [
    '<rootDir>/packages/**/*.test.ts',
    '<rootDir>/packages/**/*.test.tsx',
  ],
  transform: {
    '^.+\\.tsx?$': [
      'ts-jest',
      { tsconfig: { module: 'CommonJS', jsx: 'react-jsx' } },
    ],
  },
  moduleNameMapper: {
    '^.+\\.(css|scss)$': '<rootDir>/scripts/test/styleMock.cjs',
    '^@Base$': '<rootDir>/packages/base',
    '^@Base/(.*)$': '<rootDir>/packages/base/$1',
    '^@Business$': '<rootDir>/packages/business',
    '^@Business/(.*)$': '<rootDir>/packages/business/$1',
    '^@Components$': '<rootDir>/packages/components',
    '^@Components/(.*)$': '<rootDir>/packages/components/$1',
    '^@Domain$': '<rootDir>/packages/domain',
    '^@Domain/(.*)$': '<rootDir>/packages/domain/$1',
    '^@Styles
    '^@Shell/(.*)$': '<rootDir>/packages/shell/$1',
  },
};
: '<rootDir>/packages/styles',
    '^@Styles/(.*)
    '^@Shell/(.*)$': '<rootDir>/packages/shell/$1',
  },
};
: '<rootDir>/packages/styles/$1',
    '^@Shell
    '^@Shell/(.*)$': '<rootDir>/packages/shell/$1',
  },
};
: '<rootDir>/packages/shell',
    '^@Shell/(.*)$': '<rootDir>/packages/shell/$1',
  },
};
