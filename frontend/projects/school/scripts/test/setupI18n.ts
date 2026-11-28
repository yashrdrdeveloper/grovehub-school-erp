import { __ as resolveMessage } from '@Base/i18n';

/** Jest unit tests install the real pure formatter; do not mock message logic. */
window.__ = resolveMessage;
