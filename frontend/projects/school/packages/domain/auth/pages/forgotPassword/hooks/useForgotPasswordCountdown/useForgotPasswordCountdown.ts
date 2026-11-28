import { useCallback, useEffect, useRef, useState } from 'react';

import { FORGOT_PASSWORD_COOLDOWN_MS } from '../../constants/forgotPassword.cooldown';
import type { UseForgotPasswordCountdown } from './types/useForgotPasswordCountdown';

/** Wall-clock countdown remains accurate when the browser tab is inactive. */
export const useForgotPasswordCountdown: UseForgotPasswordCountdown = () => {
  const [secondsLeft, setSecondsLeft] = useState(0);
  const deadlineRef = useRef(0);

  useEffect(() => {
    if (secondsLeft <= 0) return;
    const interval = window.setInterval(() => {
      setSecondsLeft(Math.max(0, Math.ceil((deadlineRef.current - Date.now()) / 1000)));
    }, 250);
    return () => window.clearInterval(interval);
  }, [secondsLeft > 0]);

  const startCountdown = useCallback(() => {
    deadlineRef.current = Date.now() + FORGOT_PASSWORD_COOLDOWN_MS;
    setSecondsLeft(Math.ceil(FORGOT_PASSWORD_COOLDOWN_MS / 1000));
  }, []);
  const isCoolingDown = useCallback(() => Date.now() < deadlineRef.current, []);

  return { secondsLeft, startCountdown, isCoolingDown };
};
