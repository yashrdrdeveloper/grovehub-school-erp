import { CreditCard, MessageCircle, ShieldCheck } from '@Components/icons';

export const DASHBOARD_WORKSPACES = [
  {
    title: __('Admissions & dues'),
    hint: __('Applications, payments and follow-ups'),
    icon: CreditCard,
    tone: 'blue',
  },
  {
    title: __('Complaints & requests'),
    hint: __('Conversations that need a response'),
    icon: MessageCircle,
    tone: 'pink',
  },
  {
    title: __('Audit & compliance'),
    hint: __('Checks, reviews and approvals'),
    icon: ShieldCheck,
    tone: 'green',
  },
] as const;
