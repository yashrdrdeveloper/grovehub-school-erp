/** Original RDR photo references, not stock images. Run `pnpm assets:school`
 * to copy the school's originals into public/school-images. If an image
 * has not been copied yet, the UI falls back to the verified legacy URL.
 */
export const photographs = [
  {
    src: '/school-images/slide3.jpg',
    originalSrc: 'https://rdrschool.com/images/slide3.jpg',
    title: __('Learning through achievement'),
    alt: __('RDR School students with awards and teachers at a school event'),
    category: __('Achievements'),
    position: 'center 42%',
  },
  {
    src: '/school-images/gallery-13.jpg',
    originalSrc: 'https://rdrschool.com/images/gallery/13.jpg',
    title: __('The joy of performing'),
    alt: __('RDR School students performing a group cultural dance at a school celebration'),
    category: __('Culture'),
    position: 'center 38%',
  },
  {
    src: '/school-images/gallery-6.jpg',
    originalSrc: 'https://rdrschool.com/images/gallery/6.jpg',
    title: __('Everyday moments of discovery'),
    alt: __('Young RDR School children playing outside at school'),
    category: __('Campus Life'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-10.jpg',
    originalSrc: 'https://rdrschool.com/images/gallery/10.jpg',
    title: __('Confidence in every move'),
    alt: __('Students participating in a taekwondo demonstration at an RDR School event'),
    category: __('Sports & Activities'),
    position: 'center 35%',
  },
  {
    src: '/school-images/slide5.jpg',
    originalSrc: 'https://rdrschool.com/images/slide5.jpg',
    title: __('Growing together'),
    alt: __('RDR School students smiling together in their school uniforms'),
    category: __('Student Life'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-11.jpg',
    originalSrc: 'https://rdrschool.com/images/gallery/11.jpg',
    title: __('Celebrating our traditions'),
    alt: __('RDR School students celebrating cultural traditions in costume'),
    category: __('Culture'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-4.jpg',
    originalSrc: 'https://rdrschool.com/images/gallery/4.jpg',
    title: __('A stage for every talent'),
    alt: __('Children dancing during an RDR School cultural event'),
    category: __('Events'),
    position: 'center center',
  },
  {
    src: '/school-images/gallery-7.jpg',
    originalSrc: 'https://rdrschool.com/images/gallery/7.jpg',
    title: __('Little moments, big milestones'),
    alt: __('Students receiving awards on stage at an RDR School ceremony'),
    category: __('Achievements'),
    position: 'center center',
  },
] as const;
