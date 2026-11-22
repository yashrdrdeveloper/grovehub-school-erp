import { LandingPage } from '../../LandingPage';
import { useLandingPage } from '../../hooks/useLandingPage';

export const LandingContainer = () => <LandingPage model={useLandingPage()}/>;
