import { useState } from 'react';
import type { CSSProperties } from 'react';

import { Box } from '../Box';

import { getAvatarHue } from './helpers/avatar.getHue';
import { getAvatarInitials } from './helpers/avatar.getInitials';
import styles from './Avatar.module.css';
import type { AvatarProps } from './types/avatar.props';

export const Avatar = ({ name, src, fallback, dataTestId }: AvatarProps) => {
  const [failedSource, setFailedSource] = useState<string>();
  const style = { '--avatar-hue': getAvatarHue(name) } as CSSProperties;
  const showImage = Boolean(src && src !== failedSource);

  return (
    <Box className={styles.avatar} style={style} dataTestId={dataTestId} aria-label={name}>
      {showImage && <img src={src} alt="" onError={() => setFailedSource(src)} />}
      {!showImage && (fallback ?? getAvatarInitials(name))}
    </Box>
  );
};
