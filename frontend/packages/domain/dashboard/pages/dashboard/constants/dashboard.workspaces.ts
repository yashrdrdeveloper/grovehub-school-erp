import { CreditCard, MessageCircle, ShieldCheck } from '@Components/icons';

export const DASHBOARD_WORKSPACES = [
  {
    title: 'Admissions & dues',
    hint: 'Applications, payments and follow-ups',
    icon: CreditCard,
    tone: 'blue',
  },
  {
    title: 'Complaints & requests',
    hint: 'Conversations that need a response',
    icon: MessageCircle,
    tone: 'pink',
  },
  {
    title: 'Audit & compliance',
    hint: 'Checks, reviews and approvals',
    icon: ShieldCheck,
    tone: 'green',
  },
] as const;
