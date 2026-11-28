/**
 * Download public image files currently referenced on RDR's legacy website.
 * Usage: pnpm --dir frontend assets:school
 *
 * Photographs may include children. Confirm copyright and consent with the
 * school before reuse on the live site. Do not silently alter the individuals,
 * school premises, awards or other documentary evidence.
 */
import { mkdir, writeFile } from 'node:fs/promises';
import { join, resolve } from 'node:path';

const images = [
  ['slide3.jpg', 'https://rdrschool.com/images/slide3.jpg'],
  ['gallery-13.jpg', 'https://rdrschool.com/images/gallery/13.jpg'],
  ['gallery-6.jpg', 'https://rdrschool.com/images/gallery/6.jpg'],
  ['gallery-10.jpg', 'https://rdrschool.com/images/gallery/10.jpg'],
  ['slide5.jpg', 'https://rdrschool.com/images/slide5.jpg'],
  ['gallery-11.jpg', 'https://rdrschool.com/images/gallery/11.jpg'],
  ['gallery-4.jpg', 'https://rdrschool.com/images/gallery/4.jpg'],
  ['gallery-7.jpg', 'https://rdrschool.com/images/gallery/7.jpg'],
];

const folder = resolve(process.cwd(), 'public/school-images');
await mkdir(folder, { recursive: true });
const completedCounter = { current: 0 };
for (const [name, url] of images) {
  try {
    const response = await fetch(url, {
      signal: AbortSignal.timeout(20000),
      headers: { 'User-Agent': 'RDRSchoolSiteMigration/1.0' },
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const bytes = Buffer.from(await response.arrayBuffer());
    if (bytes.length < 1024 || bytes[0] !== 0xff || bytes[1] !== 0xd8) {
      throw new Error('Not a valid JPEG image');
    }
    if (bytes.length > 15_000_000) throw new Error('Image exceeds 15 MB limit');
    await writeFile(join(folder, name), bytes);
    console.log(`Saved ${name} (${Math.round(bytes.length / 1024)} KB)`);
    completedCounter.current++;
  } catch (error) {
    console.warn(
      `Could not download ${name}: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
}
console.log(
  `Downloaded ${completedCounter.current}/${images.length} original school images to ${folder}.`,
);
if (completedCounter.current !== images.length) process.exitCode = 1;
