import { useState } from 'react';
import { ArrowRight, Check, ExternalLink } from 'lucide-react';
import { legacyStories, storyNavigation, type StoryKey } from '../../config/sheetCopy';
import { PanelPhoto } from '../panelPhoto';

const StoryPanel = ({initial='history',navigate}:{initial?:StoryKey;navigate:(path:string)=>void}) => {
  const [selected,setSelected]=useState<StoryKey>(initial);
  const story=legacyStories[selected];
  return <div className="rdr-sheet-story">
    <div className="rdr-sheet-story__tabs" role="tablist" aria-label="School story topics">
      {storyNavigation.map(item=>{
        const Icon=item.icon;
        return <button type="button" role="tab" aria-selected={selected===item.id}
          key={item.id} className={selected===item.id?'active':''} onClick={()=>setSelected(item.id)}>
          <Icon size={16}/>{item.label}
        </button>;
      })}
    </div>
    <div className="rdr-sheet-story__layout" role="tabpanel">
      <div className="rdr-sheet-story__copy">
        <span className="rdr-sheet-kicker">{story.eyebrow}</span>
        <h3>{story.title}</h3>
        <p>{story.description}</p>
        <div className="rdr-sheet-story__points">
          {story.highlights.map(label=><span key={label}><Check size={15}/>{label}</span>)}
        </div>
        <div className="rdr-sheet-story__footer">
          <a href={story.link} target="_blank" rel="noreferrer">{story.linkLabel}<ExternalLink size={15}/></a>
          <button type="button" onClick={()=>navigate('enquiry')}>Ask us a question <ArrowRight size={16}/></button>
        </div>
        {(selected==='principal'||selected==='director')&&<small>Role and message are from the school's historical website. Please confirm current leadership before official publication.</small>}
      </div>
      <PanelPhoto index={selected==='history'?4:selected==='philosophy'?2:0}
        alt="Students participating in activities at RDR School"/>
    </div>
  </div>;
};

export { StoryPanel };
