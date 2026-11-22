/**
 * Decorative, intentionally fictional school-campus illustration.
 * This is not a photograph or architectural representation of RDR School.
 * All geometry is inline so the login screen has no media/network dependency.
 */
export const AuthCampusIllustration = () => (
  <div className="rdr-campus-art" aria-hidden="true">
    <div className="rdr-campus-art__board">
      <div className="rdr-campus-art__overline"><span className="rdr-campus-art__dot"/> A PLACE TO BELONG</div>
      <svg className="rdr-campus-art__scene" viewBox="0 0 520 340" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M45 290c55-26 108-39 162-39 59 0 93 16 139 17 51 1 90-20 130-8v51H45v-21Z" fill="#dce4db"/>
        <circle cx="417" cy="73" r="23" fill="#e7cf9e"/>
        <path d="M58 108c27-17 52-17 79 0M378 125c21-12 41-12 62 0" stroke="#d4dcd8" strokeWidth="3" strokeLinecap="round"/>
        <path d="M119 159 256 81l143 78" fill="#344763"/>
        <path d="M121 159 256 83l142 76" stroke="#273852" strokeWidth="5" strokeLinejoin="round"/>
        <path d="M138 155h238v137H138z" fill="#fffdf9" stroke="#ccd5dc" strokeWidth="3"/>
        <path d="M138 155h238v15H138z" fill="#e5e6df"/>
        <path d="M233 112h47v43h-47z" fill="#f8f4e9" stroke="#cbd2d4" strokeWidth="3"/>
        <path d="M244 124h25v31h-25z" fill="#d8e6ed"/>
        <path d="M150 187h49v44h-49zM313 187h49v44h-49z" fill="#dce9ed" stroke="#a5bbc7" strokeWidth="3"/>
        <path d="M174 187v44M337 187v44M150 209h49M313 209h49" stroke="#a5bbc7" strokeWidth="2"/>
        <path d="M222 225h70v67h-70z" fill="#344763"/>
        <path d="M256 225v67" stroke="#8ea9b4" strokeWidth="2"/>
        <path d="M244 263h-4M270 263h4" stroke="#eccd92" strokeWidth="3" strokeLinecap="round"/>
        <path d="M202 291h108v9H202zM189 300h134v9H189z" fill="#b8b7aa"/>
        <path d="M137 292h-30M376 292h33" stroke="#bfc5bf" strokeWidth="3" strokeLinecap="round"/>
        <path d="M79 285v-73" stroke="#816f58" strokeWidth="9" strokeLinecap="round"/>
        <path d="M80 188c-30-28-23-56 1-68 25 4 40 27 28 52 29 13 29 48-4 53-19 3-29-12-29-12-29 15-55-9-43-35 9-18 29-18 46 10Z" fill="#a7beac"/>
        <path d="M440 286v-77" stroke="#816f58" strokeWidth="9" strokeLinecap="round"/>
        <path d="M439 204c-25-25-20-50 4-62 21 6 31 23 24 42 25 9 28 41 3 50-15 5-31-6-31-6-22 14-42-5-37-24 4-16 24-17 37 0Z" fill="#b3c6b4"/>
        <path d="M61 313h399" stroke="#9eac9e" strokeWidth="3" strokeLinecap="round"/>
        <path d="M179 275h27M306 275h27" stroke="#d4d9d2" strokeWidth="3" strokeLinecap="round"/>
        <path d="M133 244h19M361 244h16" stroke="#cdd5cd" strokeWidth="4" strokeLinecap="round"/>
      </svg>
      <div className="rdr-campus-art__caption">LEARN <span>·</span> DISCOVER <span>·</span> GROW</div>
    </div>
    <div className="rdr-campus-art__stamp">
      <span className="rdr-campus-art__stamp-icon" aria-hidden="true">✦</span>
      <span>Every day,<strong>together.</strong></span>
    </div>
  </div>
);
