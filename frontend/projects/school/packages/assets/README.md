# SVG and image assets

Static illustrations belong to `packages/assets/<domain>/`. Standalone file assets use `.svg`, while SVGs needing inherited `currentColor` or CSS variables use a pure typed SVG-only TSX asset. Assets never own application logic, event handlers, or screens. Render through `@Components/components/SVGRenderer`, which supports inline components and file URLs. Every visual must be explicitly decorative or have accessible alt text. Do not duplicate or keep unused artwork.

## Display copy

Assets contain vectors, filenames and media resources, not the global translation function. Feature-owned wrappers provide accessible alt text and captions at render time through `__()` when meaningful. Decorative assets remain hidden from assistive technology; CSS/SVG technical identifiers and file paths are not localized. Do not import shell or `@Base/i18n` into pure assets.

## Asset-facing source-copy ownership

Media paths, SVG IDs and file keys are not translatable. Asset or gallery presentation `title`, `alt` and captions are literal `__('message')` calls in their owning content constants, not `__(photo.alt)` at the render call site. No mutable `let` declarations in frontend modules.
