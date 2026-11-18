/** Short, original school-focused lines. A single quote is selected per page load. */
export const authIllustrationQuotes = [
  'Everything your school needs, in one calm workspace.',
  'From first enquiry to every milestone.',
  'Every learner has a story worth supporting.',
  'Less paperwork. More time for learning.',
  'One school. Every journey, connected.',
  'Great schools grow through everyday care.',
  'Make every school day a little easier.',
  'From admissions to alumni, stay connected.',
  'Clearer workflows, brighter school days.',
  'Keep the details organized and the people first.',
  'Every small improvement helps a student thrive.',
  'Better coordination makes room for better teaching.',
  'Where school operations meet student growth.',
  'Support every step of the learning journey.',
  'Bring your school community closer together.',
  'The right information, when it matters.',
  'Simplify the routine. Celebrate the remarkable.',
  'Organize today. Inspire tomorrow.',
  'Thoughtful tools for meaningful school moments.',
  'One connected view, countless possibilities.',
] as const;

export const selectAuthIllustrationQuote = (seed: number): string =>
  authIllustrationQuotes[(seed >>> 0) % authIllustrationQuotes.length];
