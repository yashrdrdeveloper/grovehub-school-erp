/** Calendar boundaries use the browser's local timezone; the API receives UTC instants. */
export const getUserSetupDateBoundary = (value: string | undefined, offset: number): string | undefined => {
  if (!value || !/^\d{4}-\d{2}-\d{2}$/.test(value)) return undefined;
  const [year, month, day] = value.split('-').map(Number);
  const date = new Date(year, month - 1, day);
  if (date.getFullYear() !== year || date.getMonth() !== month - 1 || date.getDate() !== day) return undefined;
  date.setDate(date.getDate() + offset);
  return date.toISOString();
};
