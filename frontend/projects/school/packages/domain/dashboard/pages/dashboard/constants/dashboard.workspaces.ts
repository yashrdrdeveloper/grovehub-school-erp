import { CreditCard, MessageCircle, ShieldCheck } from '@Components/icons';

import { DASHBOARD_TONE } from './dashboard.tone';

export const DASHBOARD_WORKSPACES = [
  {
    title: __('Admissions & dues'),
    hint: __('Applications, payments and follow-ups'),
    icon: CreditCard,
    tone: DASHBOARD_TONE.BLUE,
  },
  {
    title: __('Complaints & requests'),
    hint: __('Conversations that need a response'),
    icon: MessageCircle,
    tone: DASHBOARD_TONE.PINK,
  },
  {
    title: __('Audit & compliance'),
    hint: __('Checks, reviews and approvals'),
    icon: ShieldCheck,
    tone: DASHBOARD_TONE.GREEN,
  },
] as const;
