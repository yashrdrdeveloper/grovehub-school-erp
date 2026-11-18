export type ArtId =
  | 'crm' | 'student' | 'tasks' | 'finance' | 'ai' | 'audit' | 'growth'
  | 'automation' | 'communication' | 'admissions' | 'hr' | 'library'
  | 'calendar' | 'book' | 'pencil' | 'bell' | 'trophy' | 'school'
  | 'quote' | 'squiggle' | 'motto';

export type ArtPlacement = { x: number; y: number };
export type ArtBounds = { width: number; height: number };
export type ArtRect = ArtBounds & { x: number; y: number };
export type ArtBox = { id: ArtId; width: number; height: number };

type Space = ArtRect;
type Rng = () => number;
const margin = 12;
const padding = 18;
const priority: ArtId[] = [
  'crm', 'student', 'tasks', 'finance', 'ai', 'audit', 'growth', 'automation',
  'communication', 'admissions', 'hr', 'library', 'quote',
  'calendar', 'book', 'pencil', 'bell', 'trophy', 'school', 'squiggle', 'motto',
];
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
    [result[i], result[j]] = [result[j], result[i]];
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

  for (const id of priority) {
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
      const preferredSide = (id === 'crm') !== flip ? sideGutters[0] : sideGutters[1];
      spaces = [preferredSide, ...spaces.filter(space => space !== preferredSide)];
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
  }
  return placed;
};
