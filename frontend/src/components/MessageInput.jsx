import { useState } from 'react';
import './MessageInput.css';

const MessageInput = ({ onSendMessage, disabled, attachments, onRemoveAttachment }) => {
  const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (message.trim() || attachments.length > 0) {
      onSendMessage(message.trim(), attachments);
      setMessage('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const isSubmitDisabled = disabled || (!message.trim() && attachments.length === 0);

  return (
    <form className="message-input-form" onSubmit={handleSubmit}>
      {attachments.length > 0 && (
        <div className="attachments-preview">
          {attachments.map((attachment, index) => (
            <div key={index} className="attachment-item">
              <span className="attachment-name">
                {attachment.type === 'file' ? attachment.file.name : attachment.url}
              </span>
              <button
                type="button"
                className="remove-attachment"
                onClick={() => onRemoveAttachment(index)}
                aria-label={`Remove ${attachment.type === 'file' ? 'file' : 'URL'}`}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}
      <div className="input-container">
        <textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Ask about products, get marketing copy, or search marketplaces..."
          disabled={disabled}
          rows={1}
          className="message-input"
          aria-label="Message input"
        />
        <button
          type="submit"
          disabled={isSubmitDisabled}
          className="send-button"
          aria-label="Send message"
        >
          Send
        </button>
      </div>
    </form>
  );
};

export default MessageInput;