import { useEffect, useRef } from 'react';
import { Box } from '../../../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../../../Button';
import type { MobileInfiniteEndProps } from './types/mobileInfiniteEnd.props';
import styles from './MobileInfiniteEnd.module.css';

export const MobileInfiniteEnd = ({
  hasNext,
  loading,
  error,
  onLoadMore,
}: MobileInfiniteEndProps) => {
  const sentinel = useRef<HTMLElement>(null);
  useEffect(() => {
    if (
      !hasNext ||
      loading ||
      error ||
      !onLoadMore ||
      !sentinel.current ||
      typeof IntersectionObserver === 'undefined'
    )
      return;
    const observer = new IntersectionObserver((entries) => {
      if (entries.some((entry) => entry.isIntersecting)) onLoadMore();
    });
    observer.observe(sentinel.current);
    return () => observer.disconnect();
  }, [hasNext, loading, error, onLoadMore]);
  if (!hasNext && !loading) return null;

  return (
    <Box className={styles.end} ref={sentinel}>
      {loading && <Box role="status">{__('Loading records…')}</Box>}
      {!loading && hasNext && (
        <Button variant={BUTTON_VARIANT.GHOST} size={BUTTON_SIZE.SMALL} onClick={onLoadMore}>
          {error && __('Retry loading')}
          {!error && __('Load more')}
        </Button>
      )}
    </Box>
  );
};
