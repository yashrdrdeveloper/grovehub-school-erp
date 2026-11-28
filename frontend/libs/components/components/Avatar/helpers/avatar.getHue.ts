/** Stable name hash, independent of theme, ordering and render cycle. */
export const getAvatarHue = (name: string): number => {
  const hash = Array.from(name.trim().normalize('NFKC').toLowerCase()).reduce(
    (value, character) => Math.imul(value ^ (character.codePointAt(0) ?? 0), 16777619),
    2166136261,
  );
  return (hash >>> 0) % 360;
};
