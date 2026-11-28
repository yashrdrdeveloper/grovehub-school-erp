import type { ReactElement } from 'react';
import type { RecoveryProps } from '../../../types';

export type RecoveryArtworkProps = Pick<RecoveryProps, 'view' | 'purpose'>;
export type RecoveryArtworkComponent = (props: RecoveryArtworkProps) => ReactElement;
