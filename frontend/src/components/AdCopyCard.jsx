import './AdCopyCard.css';

const AdCopyCard = ({ adCopy }) => {
  return (
    <div className="ad-copy-card" role="article" aria-label="Generated marketing copy">
      <div className="ad-section">
        <h4>Headlines</h4>
        <ul className="headlines-list">
          {adCopy.headlines?.map((headline, index) => (
            <li key={index} className="headline-item">{headline}</li>
          ))}
        </ul>
      </div>
      
      <div className="ad-section">
        <h4>Descriptions</h4>
        <ul className="descriptions-list">
          {adCopy.descriptions?.map((description, index) => (
            <li key={index} className="description-item">{description}</li>
          ))}
        </ul>
      </div>
      
      {adCopy.call_to_action && (
        <div className="ad-section">
          <h4>Call to Action</h4>
          <div className="cta-item">{adCopy.call_to_action}</div>
        </div>
      )}
      
      {adCopy.target_segments && adCopy.target_segments.length > 0 && (
        <div className="ad-section">
          <h4>Target Segments</h4>
          <div className="segments-list">
            {adCopy.target_segments.map((segment, index) => (
              <span key={index} className="segment-tag">{segment}</span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default AdCopyCard;