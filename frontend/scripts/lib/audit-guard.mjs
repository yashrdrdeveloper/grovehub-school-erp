/**
 * Pure audit regression evaluator. The allowlist is temporary debt, not approval
 * for new issue types or an increase in the known legacy count.
 */
export const parseAuditOutput = (output) => {
  const observed = new Map();
  const lines = output.split('\n');

  for (const line of lines) {
    if (!line.startsWith('GH_AUDIT|')) {
      continue;
    }

    const [, filename, issues] = line.split('|');
    if (!filename || !issues || issues === 'NO_PATTERN_FLAG') {
      continue;
    }

    const codes = issues.split(',').map((issue) => issue.split('@')[0]);
    observed.set(filename, new Set(codes));
  }

  return observed;
};

export const evaluateAudit = (output, baseline) => {
  const observed = parseAuditOutput(output);
  const failures = [];
  const counts = {};
  const allowlisted = baseline.allowlisted ?? {};
  const maxByCode = baseline.maxByCode ?? {};
  const findingsCounter = { current: 0 };

  for (const [filename, codes] of observed) {
    const permitted = new Set(allowlisted[filename] ?? []);

    for (const code of codes) {
      findingsCounter.current += 1;
      counts[code] = (counts[code] ?? 0) + 1;

      if (!permitted.has(code)) {
        failures.push(filename + ': newly introduced guideline violation ' + code);
      }

      if (!Object.hasOwn(maxByCode, code)) {
        failures.push(filename + ': retired or unapproved audit code ' + code);
      }
    }
  }

  for (const [code, ceiling] of Object.entries(maxByCode)) {
    if ((counts[code] ?? 0) > ceiling) {
      failures.push(code + ': ' + counts[code] + ' findings exceed ceiling ' + ceiling);
    }
  }

  if (observed.size > baseline.maxFlaggedFiles) {
    failures.push(
      'Flagged files: ' + observed.size + ' exceed ceiling ' + baseline.maxFlaggedFiles,
    );
  }

  if (findingsCounter.current > baseline.maxFindings) {
    failures.push('Findings: ' + findingsCounter.current + ' exceed ceiling ' + baseline.maxFindings);
  }

  const summaryLine = output.split('\n').find((line) => line.startsWith('GH_AUDIT_SUMMARY|'));
  if (!summaryLine) {
    failures.push('Audit output is incomplete: missing GH_AUDIT_SUMMARY');
  } else {
    try {
      const summary = JSON.parse(summaryLine.slice('GH_AUDIT_SUMMARY|'.length));
      if (summary.flagged !== observed.size) {
        failures.push('Audit summary flagged-file count does not match the observed entries');
      }
      if (Object.entries(counts).some(([code, count]) => summary.counts[code] !== count)) {
        failures.push('Audit summary issue counts do not match the observed entries');
      }
    } catch {
      failures.push('Audit summary contains invalid JSON');
    }
  }

  const stale = [];
  for (const [filename, permitted] of Object.entries(allowlisted)) {
    for (const code of permitted) {
      if (!observed.get(filename)?.has(code)) {
        stale.push(filename + ': ' + code);
      }
    }
  }

  for (const entry of stale) {
    failures.push('Obsolete allowlist entry: ' + entry);
  }

  return { failures, counts, findings: findingsCounter.current, flaggedFiles: observed.size, stale };
};
