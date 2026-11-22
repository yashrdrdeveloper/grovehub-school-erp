import { Link } from 'react-router-dom';
import { ArrowRight, ArrowUpRight } from '@Components/icons';
import { photographs } from '../../constants/school.content';
import { EnquiryLink } from '../siteActions';
import { ImageCard } from '../siteMedia';

const HomePage = () => (
  <div className="rdr-option-one">
    <section className="rdr-o1-hero">
      <div className="rdr-container rdr-o1-hero__grid">
        <div className="rdr-o1-hero__copy">
          <span className="rdr-eyebrow">A COMMUNITY OF LEARNERS · SINCE 2006</span>
          <h1>Rooted in Haluhera.<em>Ready for every tomorrow.</em></h1>
          <p>At Rao Deena Ram Vidya Vihar, learning goes beyond the classroom. We bring together academics, character and a caring school community so every child has room to grow.</p>
          <div className="rdr-o1-hero__actions">
            <EnquiryLink className="rdr-btn rdr-btn--solid">Enquire for admission <ArrowUpRight size={17}/></EnquiryLink>
            <Link className="rdr-btn rdr-o1-secondary" to="/about">Discover our school <ArrowRight size={17}/></Link>
          </div>
          <div className="rdr-o1-hero__proof"><span>CBSE affiliated</span><span>Senior Secondary</span><span>Rewari, Haryana</span></div>
        </div>
        <div className="rdr-o1-hero__images">
          <ImageCard className="rdr-o1-hero__main-photo" src={photographs[4].src} alt={photographs[4].alt} position={photographs[4].position} loading="eager">
            <span className="rdr-o1-hero__photo-caption">A school full of familiar faces.</span>
          </ImageCard>
          <div className="rdr-o1-hero__image-side">
            <ImageCard className="rdr-o1-hero__side-photo" src={photographs[0].src} alt={photographs[0].alt} position={photographs[0].position} loading="eager"/>
            <div className="rdr-o1-hero__note"><strong>Learning with purpose.</strong><small>Growing together since 2006.</small></div>
          </div>
        </div>
      </div>
    </section>
    <section className="rdr-o1-life rdr-container">
      <div className="rdr-o1-life__heading">
        <div><span className="rdr-eyebrow">LIFE AT RDR</span><h2>Every child has a story.<br/>Every day has a purpose.</h2></div>
        <p>From classrooms to celebrations, these are the experiences that shape school life.</p>
      </div>
      <div className="rdr-o1-life__grid">
        <Link to="/gallery" className="rdr-o1-life__card">
          <ImageCard src={photographs[0].src} alt={photographs[0].alt} position={photographs[0].position}/>
          <span>ACHIEVEMENT</span><h3>Celebrating every milestone.</h3><p>Recognising effort, confidence and achievement.</p>
        </Link>
        <Link to="/gallery" className="rdr-o1-life__card">
          <ImageCard src={photographs[1].src} alt={photographs[1].alt} position={photographs[1].position}/>
          <span>EXPRESSION</span><h3>More than textbooks.</h3><p>Music, movement and the joy of participating.</p>
        </Link>
        <Link to="/gallery" className="rdr-o1-life__card">
          <ImageCard src={photographs[2].src} alt={photographs[2].alt} position={photographs[2].position}/>
          <span>BELONGING</span><h3>Growing with confidence.</h3><p>Creating space for curiosity and friendships.</p>
        </Link>
      </div>
    </section>
    <section className="rdr-o1-cta">
      <div className="rdr-container rdr-o1-cta__inner">
        <div><h2>Come see what makes RDR special.</h2><p>Questions about admissions, classes or a campus visit?</p></div>
        <EnquiryLink className="rdr-btn rdr-btn--light">Start an enquiry <ArrowUpRight size={17}/></EnquiryLink>
      </div>
    </section>
    <section className="rdr-o1-disclosure rdr-container">
      <span>Looking for school documents?</span>
      <Link to="/mandatory-public-disclosure">Mandatory Public Disclosure <ArrowUpRight size={16}/></Link>
    </section>
  </div>
);

export { HomePage };
