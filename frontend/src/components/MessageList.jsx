import { useState, useEffect, useRef } from 'react';
import ProductCard from './ProductCard';
import AdCopyCard from './AdCopyCard';
import './MessageList.css';

const MessageList = ({ messages, isLoading }) => {
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const formatTimestamp = (timestamp) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const renderMessageContent = (message) => {
    if (message.type === 'search_results' && message.data?.products) {
      return (
        <div className="search-results">
          <p>{message.text}</p>
          <div className="product-grid">
            {message.data.products.map((product, index) => (
              <ProductCard key={index} product={product} />
            ))}
          </div>
        </div>
      );
    }

    if (message.type === 'marketing_copy' && message.data?.marketing) {
      return (
        <div className="marketing-results">
          <p>{message.text}</p>
          <AdCopyCard adCopy={message.data.marketing} />
        </div>
      );
    }

    return <p>{message.text}</p>;
  };

  return (
    <div className="message-list" role="log" aria-live="polite">
      {messages.map((message) => (
        <div
          key={message.id}
          className={`message ${message.sender}`}
          role="article"
          aria-label={`Message from ${message.sender}`}
        >
          <div className="message-header">
            <span className="sender">{message.sender === 'user' ? 'You' : 'Assistant'}</span>
            <span className="timestamp">{formatTimestamp(message.timestamp)}</span>
          </div>
          <div className="message-content">
            {renderMessageContent(message)}
          </div>
        </div>
      ))}
      {isLoading && (
        <div className="message assistant loading" role="status" aria-label="Assistant is typing">
          <div className="message-header">
            <span className="sender">Assistant</span>
          </div>
          <div className="message-content">
            <div className="typing-indicator">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
        </div>
      )}
      <div ref={messagesEndRef} />
    </div>
  );
};

export default MessageList;