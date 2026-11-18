import { useLayoutEffect, useRef, useState, type CSSProperties, type ReactNode } from 'react';
import {
  ArrowRight, BookOpen, CalendarDays, GraduationCap, ListTodo, Monitor,
  Moon, Sparkles, Sun, UsersRound, WalletCards, BrainCircuit, ClipboardCheck,
  TrendingUp, Workflow, MessageSquareText, UserRoundPlus, Users, LibraryBig,
  Pencil, BellRing, Trophy, School,
} from 'lucide-react';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Shell/theme';

import {
  packAuthIllustrations,
  type ArtId,
  type ArtPlacement,
} from './authArtworkLayout';

const themeChoices: { value: ThemePreference; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: 'System', icon: Monitor },
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
];

const artworkIds: ArtId[] = [
  'crm', 'student', 'tasks', 'finance', 'ai', 'audit', 'growth',
  'automation', 'communication', 'admissions', 'hr', 'library',
  'quote', 'calendar', 'book', 'pencil', 'bell', 'trophy', 'school',
  'squiggle', 'motto',
];

/** Decorative objects are packed into unused measured space, never on the form. */
const FloatingSchoolArtwork = ({
  register,
  placements,
}: {
  register: (id: ArtId, node: HTMLDivElement | null) => void;
  placements: Partial<Record<ArtId, ArtPlacement>>;
}) => {
  const item = (id: ArtId, extraClass: string, children: ReactNode) => {
    const placement = placements[id];
    const style: CSSProperties | undefined = placement
      ? { left: placement.x, top: placement.y, visibility: 'visible' }
      : undefined;
    return (
      <div
        key={id}
        ref={node => register(id, node)}
        className={`auth-float auth-float--${id} ${extraClass}`}
        style={style}
        aria-hidden="true"
      >
        {children}
      </div>
    );
  };

  return (
    <div className="auth-floating-art" aria-hidden="true">
      {item('crm', 'auth-float-card auth-float-card--hero',
        <>
          <div className="auth-float-card__title"><UsersRound size={18}/><strong>Lead CRM</strong></div>
          <div className="auth-float-card__rule"/>
          <div className="auth-float-card__journey"><span>Enquiry</span><ArrowRight size={13}/><span>Admission</span></div>
          <div className="auth-float-card__contact"><i/><b/><span>✓</span></div>
          <div className="auth-float-card__contact"><i/><b/><span>✓</span></div>
        </>,
      )}
      {item('student', 'auth-float-card auth-float-card--hero',
        <>
          <div className="auth-float-card__title"><GraduationCap size={18}/><strong>Student Lifecycle 360°</strong></div>
          <small>Unified student view · UCV</small>
          <div className="auth-float-card__rule"/>
          <div className="auth-float-card__student">
            <span className="auth-float-card__avatar"/>
            <span className="auth-float-card__lines"><i/><i/></span>
            <span className="auth-float-card__ring"/>
          </div>
          <div className="auth-float-card__tags"><span>Academics</span><span>Growth</span></div>
        </>,
      )}
      {item('tasks', 'auth-float-card auth-float-card--small',
        <><ListTodo size={17}/><span><strong>Tasks</strong><small>Daily school work</small></span></>,
      )}
      {item('finance', 'auth-float-card auth-float-card--small',
        <><WalletCards size={17}/><span><strong>Finance &amp; Fees</strong><small>School finances</small></span></>,
      )}
      {item('ai', 'auth-float-card auth-float-card--small auth-float-card--ai',
        <><BrainCircuit size={17}/><span><strong>AI Assistance</strong><small>Smarter school workflows</small></span></>,
      )}
      {item('audit', 'auth-float-card auth-float-card--small',
        <><ClipboardCheck size={17}/><span><strong>Audit & Compliance</strong><small>Always ready</small></span></>,
      )}
      {item('growth', 'auth-float-card auth-float-card--small',
        <><TrendingUp size={17}/><span><strong>Learn & Grow</strong><small>Student progress</small></span></>,
      )}
      {item('automation', 'auth-float-card auth-float-card--small',
        <><Workflow size={17}/><span><strong>Automation</strong><small>Less repetitive work</small></span></>,
      )}
      {item('communication', 'auth-float-card auth-float-card--small',
        <><MessageSquareText size={17}/><span><strong>Communication</strong><small>School connections</small></span></>,
      )}
      {item('admissions', 'auth-float-card auth-float-card--small',
        <><UserRoundPlus size={17}/><span><strong>Admissions</strong><small>Every first step</small></span></>,
      )}
      {item('hr', 'auth-float-card auth-float-card--small',
        <><Users size={17}/><span><strong>People & HR</strong><small>Teams together</small></span></>,
      )}
      {item('library', 'auth-float-card auth-float-card--small',
        <><LibraryBig size={17}/><span><strong>School Library</strong><small>Knowledge organized</small></span></>,
      )}
      {item('quote', 'auth-float-quote',
        <>“Everything your school needs, in one calm workspace.”</>,
      )}
      {item('calendar', 'auth-float-doodle',
        <CalendarDays size={24} strokeWidth={1.35}/>,
      )}
      {item('book', 'auth-float-doodle',
        <BookOpen size={24} strokeWidth={1.35}/>,
      )}
      {item('pencil', 'auth-float-doodle',
        <Pencil size={24} strokeWidth={1.35}/>,
      )}
      {item('bell', 'auth-float-doodle',
        <BellRing size={24} strokeWidth={1.35}/>,
      )}
      {item('trophy', 'auth-float-doodle',
        <Trophy size={24} strokeWidth={1.35}/>,
      )}
      {item('school', 'auth-float-doodle',
        <School size={24} strokeWidth={1.35}/>,
      )}
      {item('squiggle', 'auth-float-doodle auth-float-doodle--line',
        <svg viewBox="0 0 140 45" fill="none"><path d="M2 30c16 25 28-30 44-11 12 16 0 29 19 22 22-10 22-43 42-24 12 13 18 19 31 12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/></svg>,
      )}
      {item('motto', 'auth-float-motto',
        <>LEARN <span>·</span> ORGANIZE <span>·</span> GROW</>,
      )}
    </div>
  );
};

export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => {
  const { preference, setPreference } = useWorkspaceTheme();
  const stageRef = useRef<HTMLDivElement>(null);
  const panelRef = useRef<HTMLElement>(null);
  const itemRefs = useRef<Partial<Record<ArtId, HTMLDivElement>>>({});
  const seedRef = useRef(Math.floor(Math.random() * 4294967295));
  const [placements, setPlacements] = useState<Partial<Record<ArtId, ArtPlacement>>>({});

  const register = (id: ArtId, node: HTMLDivElement | null) => {
    if (node) itemRefs.current[id] = node;
    else delete itemRefs.current[id];
  };

  useLayoutEffect(() => {
    const stage = stageRef.current;
    const panel = panelRef.current;
    if (!stage || !panel) return;

    let frame = 0;
    const recalculate = () => {
      window.cancelAnimationFrame(frame);
      frame = window.requestAnimationFrame(() => {
        const stageBounds = stage.getBoundingClientRect();
        const formBounds = panel.getBoundingClientRect();
        const boxes = artworkIds.flatMap(id => {
          const node = itemRefs.current[id];
          return node ? [{ id, width: node.offsetWidth, height: node.offsetHeight }] : [];
        });
        const next = packAuthIllustrations(
          { width: stageBounds.width, height: stageBounds.height },
          {
            x: formBounds.left - stageBounds.left,
            y: formBounds.top - stageBounds.top,
            width: formBounds.width,
            height: formBounds.height,
          },
          boxes,
          seedRef.current,
        );
        setPlacements(current =>
          JSON.stringify(current) === JSON.stringify(next) ? current : next,
        );
      });
    };

    const observer = new ResizeObserver(recalculate);
    observer.observe(stage);
    observer.observe(panel);
    for (const node of Object.values(itemRefs.current)) {
      if (node) observer.observe(node);
    }
    recalculate();
    return () => {
      observer.disconnect();
      window.cancelAnimationFrame(frame);
    };
  }, [title]);

  return (
    <main className="auth-page auth-page--consistent">
      <header className="auth-topbar">
        <div className="auth-brand auth-brand--school" aria-label="RDR School">
          <SchoolMonogram/>
          <strong>RDR School</strong>
        </div>
        <div className="auth-theme-switcher" role="group" aria-label="Appearance">
          {themeChoices.map(({ value, label, icon: Icon }) => (
            <button key={value} type="button" title={label} aria-label={`${label} theme`}
              aria-pressed={preference === value}
              className={preference === value ? 'auth-theme-switcher__choice is-active' : 'auth-theme-switcher__choice'}
              onClick={() => setPreference(value)}>
              <Icon size={15}/><span>{label}</span>
            </button>
          ))}
        </div>
      </header>
      <div ref={stageRef} className="auth-stage">
        <FloatingSchoolArtwork register={register} placements={placements}/>
        <section ref={panelRef} className="auth-panel" aria-label={title}>{children}</section>
      </div>
      <footer className="auth-footer"><span>© {new Date().getFullYear()} GrooveHub</span></footer>
    </main>
  );
};
