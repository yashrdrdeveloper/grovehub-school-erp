export type ArtId =
  | 'crm' | 'student' | 'tasks' | 'finance' | 'quote'
  | 'calendar' | 'book' | 'squiggle' | 'motto';

export type ArtPlacement = { x: number; y: number };
export type ArtBounds = { width: number; height: number };
export type ArtRect = ArtBounds & { x: number; y: number };
export type ArtBox = { id: ArtId; width: number; height: number };

type Space = { x: number; y: number; width: number; height: number };
type Rng = () => number;

const gap = 14;
const margin = 7;
const priority: ArtId[] = [
  'crm', 'student', 'tasks', 'finance', 'quote', 'calendar', 'book', 'squiggle', 'motto',
];

const makeRandom = (seed: number): Rng => {
  let value = seed >>> 0;
  return () => {
    value = (Math.imul(value, 1664525) + 1013904223) >>> 0;
    return value / 4294967296;
  };
};

const collides = (a: ArtRect, b: ArtRect) =>
  a.x < b.x + b.width + gap &&
  a.x + a.width + gap > b.x &&
  a.y < b.y + b.height + gap &&
  a.y + a.height + gap > b.y;

/** Break the free canvas into actual rectangles around the protected login card.
 * Every entity is measured independently; a failed fit simply is not rendered.
 */
const availableSpaces = (stage: ArtBounds, panel: ArtRect): Space[] => {
  const left = Math.max(margin, panel.x - gap);
  const right = Math.min(stage.width - margin, panel.x + panel.width + gap);
  const top = Math.max(margin, panel.y - gap);
  const bottom = Math.min(stage.height - margin, panel.y + panel.height + gap);

  return [
    { x: margin, y: margin, width: left - margin, height: stage.height - margin * 2 },
    { x: right, y: margin, width: stage.width - right - margin, height: stage.height - margin * 2 },
    { x: left, y: margin, width: right - left, height: top - margin },
    { x: left, y: bottom, width: right - left, height: stage.height - bottom - margin },
  ].filter(space => space.width > 0 && space.height > 0);
};

export const packAuthIllustrations = (
  stage: ArtBounds,
  panel: ArtRect,
  boxes: ArtBox[],
  seed: number,
): Partial<Record<ArtId, ArtPlacement>> => {
  const result: Partial<Record<ArtId, ArtPlacement>> = {};
  const occupied: ArtRect[] = [panel];
  const random = makeRandom(seed);
  const spaces = availableSpaces(stage, panel);
  const byId = new Map(boxes.map(box => [box.id, box]));

  // The primary CRM / UCV illustrations start on opposite sides where there
  // is room. Their actual positions and the positions of all other objects
  // vary by seed and by measured available space.
  const leftFirst = random() > .5;
  const heroPreference = {
    crm: leftFirst ? 0 : 1,
    student: leftFirst ? 1 : 0,
  };

  for (const id of priority) {
    const box = byId.get(id);
    if (!box) continue;
    // Rotated card artwork has a slightly larger visual footprint.
    const paddedWidth = box.width + (id === 'crm' || id === 'student' ? 14 : 2);
    const paddedHeight = box.height + (id === 'crm' || id === 'student' ? 12 : 2);
    const candidates = spaces
      .map((space, index) => ({ space, index }))
      .filter(({ space }) => space.width >= paddedWidth && space.height >= paddedHeight);
    if (!candidates.length) continue;

    if (id === 'crm' || id === 'student') {
      candidates.sort((a, b) =>
        Number(b.index === heroPreference[id]) - Number(a.index === heroPreference[id]));
    } else {
      // A different stable order each page load, without jitter on resize.
      candidates.sort(() => random() - .5);
    }

    let chosen: ArtRect | undefined;
    for (const { space } of candidates) {
      const availableX = space.width - paddedWidth;
      const availableY = space.height - paddedHeight;
      for (let attempt = 0; attempt < 80; attempt += 1) {
        const x = space.x + (attempt === 0 ? availableX / 2 : random() * availableX);
        const y = space.y + (attempt === 0 ? availableY / 2 : random() * availableY);
        const rect = { x, y, width: paddedWidth, height: paddedHeight };
        if (!occupied.some(other => collides(rect, other))) {
          chosen = rect;
          break;
        }
      }
      if (chosen) break;
    }
    if (!chosen) continue;
    occupied.push(chosen);
    result[id] = {
      x: Math.round(chosen.x + (chosen.width - box.width) / 2),
      y: Math.round(chosen.y + (chosen.height - box.height) / 2),
    };
  }

  return result;
};
