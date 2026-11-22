import { Component } from 'react';

import { Box } from '../Box';
import { Button } from '../Button';
import { TYPOGRAPHY_TAG, TYPOGRAPHY_VARIANT, Typography } from '../Typography';

import { ERROR_BOUNDARY_COPY } from './constants/errorBoundary.copy';
import { ERROR_BOUNDARY_TEST_IDS } from './constants/errorBoundary.testIds';
import type { ErrorBoundaryProps, ErrorBoundaryState } from './types/errorBoundary.props';

/**
 * Render-time failure isolation. Network and form errors remain owned by their
 * domain; internal error details are never shown to parents or students.
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = {
    hasError: false,
  };

  static getDerivedStateFromError(): ErrorBoundaryState {
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: { componentStack?: string | null }): void {
    this.props.onError?.(error, info.componentStack ?? '');
  }

  componentDidUpdate(previous: ErrorBoundaryProps): void {
    if (previous.resetKey !== this.props.resetKey && this.state.hasError) {
      this.setState({ hasError: false });
    }
  }

  reset = (): void => {
    this.setState({ hasError: false });
  };

  render() {
    if (!this.state.hasError) {
      return this.props.children;
    }

    if (this.props.fallback) {
      return this.props.fallback(this.reset);
    }

    return (
      <Box
        as="section"
        role="alert"
        aria-live="assertive"
        className="gh-error-boundary"
        dataTestId={this.props.dataTestId ?? ERROR_BOUNDARY_TEST_IDS.ROOT}
      >
        <Typography as={TYPOGRAPHY_TAG.H2} variant={TYPOGRAPHY_VARIANT.HEADING}>
          {ERROR_BOUNDARY_COPY.TITLE}
        </Typography>
        <Typography as={TYPOGRAPHY_TAG.P}>
          {ERROR_BOUNDARY_COPY.DESCRIPTION}
        </Typography>
        <Button type="button" dataTestId={ERROR_BOUNDARY_TEST_IDS.RETRY} onClick={this.reset}>
          {ERROR_BOUNDARY_COPY.RETRY}
        </Button>
      </Box>
    );
  }
}
