import type { PageHeadingComponent } from './types/pageHeading.props';

export const PageHeading: PageHeadingComponent = ({
  eyebrow, title, description, action, dataTestId,
}) => (
  <header className="page-heading" data-testid={dataTestId}>
    <div>
      {eyebrow && <span className="eyebrow">{eyebrow}</span>}
      <h1>{title}</h1>
      {description && <p>{description}</p>}
    </div>
    {action}
  </header>
);
