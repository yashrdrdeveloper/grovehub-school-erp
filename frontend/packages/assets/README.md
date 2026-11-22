# SVG and image assets

Static illustrations belong to `packages/assets/<domain>/`. Standalone file assets use `.svg`, while SVGs needing inherited `currentColor` or CSS variables use a pure typed SVG-only TSX asset. Assets never own application logic, event handlers, or screens. Render through `@Components/components/SVGRenderer`, which supports inline components and file URLs. Every visual must be explicitly decorative or have accessible alt text. Do not duplicate or keep unused artwork.
