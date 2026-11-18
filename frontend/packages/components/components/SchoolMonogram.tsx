import { cn } from '@Base/utilities';

export interface SchoolMonogramProps {
  className?: string;
}

/**
 * RDR School's distinctive connected-letter logo, traced from the approved
 * navy reference artwork. Do not replace the symbol with an RDR text glyph.
 */
export const SchoolMonogram = ({ className }: SchoolMonogramProps) => (
  <span className={cn('school-monogram', className)} aria-hidden="true">
    <svg className="school-monogram__symbol" viewBox="0 0 553 229" xmlns="http://www.w3.org/2000/svg" fill="none" focusable="false">
      <path d="M0 1 1 49 121 49 128 51 134 55 141 64 143 70 143 83 138 93 130 100 120 104 0 104 0 228 50 228 51 151 75 151 138 228 206 227 140 149 141 147 161 140 171 133 180 124 187 114 194 98 196 87 195 59 191 46 179 26 161 11 145 4 132 1Z" fill="currentColor"/>
      <path d="M160 1 171 7 187 21 195 32 203 49 266 49 286 54 297 60 309 71 315 79 321 94 324 111 322 130 319 139 310 154 300 164 285 172 270 175 241 175 239 156 240 80 208 79 208 88 205 102 198 118 191 128 177 142 161 152 222 227 278 227 294 225 312 220 328 212 345 199 357 185 369 164 375 146 374 145 376 141 378 124 378 104 372 74 364 56 350 36 340 26 322 14 309 8 280 1Z" fill="currentColor"/>
      <path d="M315 1 329 7 347 19 361 33 371 48 461 49 475 55 482 63 485 71 484 86 478 95 470 101 459 104 391 104 389 139 383 162 374 181 362 198 351 209 339 218 320 228 401 227 402 151 412 151 478 228 552 227 478 148 498 143 515 133 527 121 537 102 540 89 540 64 537 51 528 33 515 19 505 12 492 6 462 0Z" fill="currentColor"/>
    </svg>
  </span>
);
