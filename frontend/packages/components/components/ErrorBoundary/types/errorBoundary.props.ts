import type { ReactNode } from 'react';

export interface ErrorBoundaryProps {
  readonly children: ReactNode;
  readonly resetKey?: string;
  readonly fallback?: (reset: () => void) => ReactNode;
  readonly onError?: (error: Error, componentStack: string) => void;
}

export interface ErrorBoundaryState {
  readonly hasError: boolean;
}

export type ErrorBoundaryReset = () => void;
