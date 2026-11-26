/** Calendar boundaries use the browser's local timezone; the API receives UTC instants. */
export const getUserSetupDateBoundary = (value: string | undefined, offset: number): string | undefined => {
  if (!value || !/^\d{4}-\d{2}-\d{2}$/.test(value)) return undefined;
  const year = Number(value.slice(0, 4));
  const month = Number(value.slice(5, 7));
  const day = Number(value.slice(8, 10));
  const date = new Date(year, month - 1, day);
  if (date.getFullYear() !== year || date.getMonth() !== month - 1 || date.getDate() !== day) return undefined;
  date.setDate(date.getDate() + offset);
  return date.toISOString();
};
