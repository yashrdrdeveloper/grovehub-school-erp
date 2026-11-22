export type ArtId =
  | 'crm' | 'student' | 'tasks' | 'finance' | 'ai' | 'audit' | 'growth'
  | 'automation' | 'communication' | 'admissions' | 'hr' | 'library'
  | 'calendar' | 'book' | 'pencil' | 'bell' | 'trophy' | 'school'
  | 'quote' | 'squiggle' | 'plane' | 'motto';

export type ArtPlacement = { x: number; y: number };
export type ArtBounds = { width: number; height: number };
export type ArtRect = ArtBounds & { x: number; y: number };
export type ArtBox = { id: ArtId; width: number; height: number };

type Space = ArtRect;
type Rng = () => number;
const margin = 12;
const padding = 18;
const heroModules: ArtId[] = ['crm', 'student'];
const subModules: ArtId[] = ['tasks', 'finance', 'ai', 'audit', 'growth', 'automation', 'communication', 'admissions', 'hr', 'library'];
const lineArt: ArtId[] = ['plane', 'squiggle'];
const smallIcons: ArtId[] = ['calendar', 'book', 'pencil', 'bell', 'trophy', 'school'];
const optionalAccents: ArtId[] = ['motto'];
const makeRandom = (seed: number): Rng => {
  let value = seed >>> 0;
  return () => {
    value = (Math.imul(value, 1664525) + 1013904223) >>> 0;
    return value / 4294967296;
  };
};
const overlaps = (a: ArtRect, b: ArtRect) =>
  a.x < b.x + b.width + padding &&
  a.x + a.width + padding > b.x &&
  a.y < b.y + b.height + padding &&
  a.y + a.height + padding > b.y;

const shuffle = <T,>(items: T[], rng: Rng): T[] => {
  const result = [...items];
  for (let i = result.length - 1; i > 0; i -= 1) {
    const j = Math.floor(rng() * (i + 1));
    const current = result[i];
    const selected = result[j];
    if (current === undefined || selected === undefined) continue;
    result[i] = selected;
    result[j] = current;
  }
  return result;
};

/** The sign-in panel is a protected region. Free-space candidates cover the
 * entire available viewport, with four gutters around the panel. Each entity
 * either receives a non-overlapping rectangle or is omitted. A stable seed
 * makes each refresh unique but avoids random jitter during resizing.
 */
export const packAuthIllustrations = (
  stage: ArtBounds, panel: ArtRect, boxes: ArtBox[], seed: number,
): Partial<Record<ArtId, ArtPlacement>> => {
  const placed: Partial<Record<ArtId, ArtPlacement>> = {};
  if (stage.width <= 0 || stage.height <= 0) return placed;
  const rng = makeRandom(seed);
  const safePanel: ArtRect = {
    x: panel.x - padding, y: panel.y - padding,
    width: panel.width + padding * 2, height: panel.height + padding * 2,
  };
  const occupied: ArtRect[] = [safePanel];

  const sideGutters: Space[] = [
    { x: margin, y: margin, width: Math.max(0, safePanel.x - margin), height: stage.height - margin * 2 },
    { x: safePanel.x + safePanel.width, y: margin,
      width: Math.max(0, stage.width - safePanel.x - safePanel.width - margin), height: stage.height - margin * 2 },
  ];
  const horizontalGutters: Space[] = [
    { x: margin, y: margin, width: stage.width - margin * 2, height: Math.max(0, safePanel.y - margin) },
    { x: margin, y: safePanel.y + safePanel.height,
      width: stage.width - margin * 2,
      height: Math.max(0, stage.height - safePanel.y - safePanel.height - margin) },
  ];
  const allSpaces = [...sideGutters, ...horizontalGutters];
  const sizes = new Map(boxes.map(box => [box.id, box]));
  const flip = rng() < .5;
  const firstHero = shuffle(heroModules, rng);
  const shuffledSubModules = shuffle(subModules, rng);
  const shuffledLines = shuffle(lineArt, rng);
  const shuffledIcons = shuffle(smallIcons, rng);

  // Reserve one opportunity from each category before any optional extras.
  // Trying alternates within each category avoids a large item consuming the
  // space that should be used by the quote or a smaller school illustration.
  // The sign-in form always takes precedence: impossible fits are omitted.
  const requiredGroups: ArtId[][] = [
    firstHero,
    ['quote'],
    shuffledSubModules,
    shuffledLines,
    shuffledIcons,
  ];
  const extraModules = stage.width >= 1250 ? 7 : stage.width >= 850 ? 5 : 2;
  const extraIcons = stage.width >= 1250 ? 5 : stage.width >= 850 ? 3 : 1;
  // Unlike the mandatory groups, each optional item gets its own attempt.
  // This fills all remaining usable space, rather than stopping after a
  // single extra submodule or a single extra school icon.
  const optionalGroups: ArtId[][] = [
    ...firstHero.map(id => [id]),
    ...shuffledSubModules.slice(0, extraModules).map(id => [id]),
    ...shuffledIcons.slice(0, extraIcons).map(id => [id]),
    ...shuffledLines.map(id => [id]),
    ...optionalAccents.map(id => [id]),
  ];
  // Each required category attempts alternatives until one item fits.
  const lineup = [...requiredGroups, ...optionalGroups];
  for (const group of lineup) {
    let categoryPlaced = false;
    for (const id of group) {
      if (categoryPlaced || placed[id]) continue;
      const box = sizes.get(id);
    if (!box || box.width < 1 || box.height < 1) continue;
    // Rotations must fit, including their extended corners.
    const rotate = id === 'crm' || id === 'student' || id === 'tasks' || id === 'finance'
      || id === 'ai' || id === 'audit' || id === 'growth' || id === 'automation';
    const overhang = rotate ? 16 : 6;
    const width = box.width + overhang;
    const height = box.height + overhang;
    let spaces = shuffle(allSpaces, rng);
    if (id === 'crm' || id === 'student') {
      const preferredSide = ((id === 'crm') !== flip ? sideGutters[0] : sideGutters[1]) ?? sideGutters[0];
      if (preferredSide) spaces = [preferredSide, ...spaces.filter(space => space !== preferredSide)];
    } else if (['tasks', 'admissions', 'communication'].includes(id) && placed.crm) {
      const center = placed.crm.x + (sizes.get('crm')?.width ?? 0) / 2;
      spaces.sort((a, b) => Math.abs(a.x + a.width / 2 - center) - Math.abs(b.x + b.width / 2 - center));
    } else if (['finance', 'ai', 'growth', 'audit'].includes(id) && placed.student) {
      const center = placed.student.x + (sizes.get('student')?.width ?? 0) / 2;
      spaces.sort((a, b) => Math.abs(a.x + a.width / 2 - center) - Math.abs(b.x + b.width / 2 - center));
    }

    let accepted: ArtRect | undefined;
    for (const space of spaces) {
      if (space.width < width || space.height < height) continue;
      const maxX = space.width - width;
      const maxY = space.height - height;
      // Try distinct centers, corners and seeded locations. The first large
      // cards establish the composition; smaller cards fill remaining gaps.
      const candidates = [
        { x: maxX / 2, y: maxY / 2 },
        { x: maxX * .16, y: maxY * .22 },
        { x: maxX * .84, y: maxY * .78 },
        ...Array.from({ length: 100 }, () => ({ x: rng() * maxX, y: rng() * maxY })),
      ];
      for (const candidate of candidates) {
        const rect = {
          x: space.x + candidate.x,
          y: space.y + candidate.y,
          width,
          height,
        };
        if (!occupied.some(other => overlaps(rect, other))) {
          accepted = rect;
          break;
        }
      }
      if (accepted) break;
    }
    if (!accepted) continue;
    occupied.push(accepted);
    placed[id] = {
      x: Math.round(accepted.x + overhang / 2),
      y: Math.round(accepted.y + overhang / 2),
    };
    categoryPlaced = true;
    }
  }
  return placed;
};
