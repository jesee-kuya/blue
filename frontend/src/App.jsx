import { useState, useEffect } from 'react';
import MessageList from './components/MessageList';
import MessageInput from './components/MessageInput';
import FileUploader from './components/FileUploader';
import ErrorBoundary from './components/ErrorBoundary';
import { sendMessage, getMarketingCopy } from './services/api';
import './App.css';

function App() {
  const [messages, setMessages] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [attachments, setAttachments] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Add welcome message
    setMessages([{
      id: Date.now(),
      sender: 'assistant',
      text: 'Hello! I\'m your Blue marketplace assistant. I can help you search for products, generate marketing copy, and more. How can I assist you today?',
      timestamp: Date.now(),
      type: 'text'
    }]);
  }, []);

  const addMessage = (sender, text, type = 'text', data = null) => {
    const message = {
      id: Date.now() + Math.random(),
      sender,
      text,
      timestamp: Date.now(),
      type,
      data
    };
    setMessages(prev => [...prev, message]);
    return message;
  };

  const handleSendMessage = async (messageText, messageAttachments) => {
    setError(null);
    
    // Add user message
    addMessage('user', messageText || 'Sent attachments');
    
    setIsLoading(true);
    
    try {
      // Determine if this is a marketing request
      const isMarketingRequest = messageText.toLowerCase().includes('marketing') || 
                                messageText.toLowerCase().includes('ad copy') ||
                                messageText.toLowerCase().includes('campaign');
      
      let response;
      if (isMarketingRequest) {
        response = await getMarketingCopy(messageText);
      } else {
        response = await sendMessage(messageText, messageAttachments);
      }
      
      // Handle different response types
      if (response.search_results && response.search_results.length > 0) {
        addMessage('assistant', 'Here are the products I found:', 'search_results', {
          products: response.search_results
        });
      } else if (response.marketing_copy) {
        addMessage('assistant', 'Here\'s the marketing copy I generated:', 'marketing_copy', {
          marketing: response.marketing_copy
        });
      } else if (response.message) {
        addMessage('assistant', response.message);
      } else {
        addMessage('assistant', 'I received your request but couldn\'t process it properly. Please try again.');
      }
      
    } catch (error) {
      console.error('API Error:', error);
      setError(error.message);
      addMessage('assistant', `Sorry, I encountered an error: ${error.message}`);
    } finally {
      setIsLoading(false);
      setAttachments([]);
    }
  };

  const handleAddAttachment = (attachment) => {
    setAttachments(prev => [...prev, attachment]);
  };

  const handleRemoveAttachment = (index) => {
    setAttachments(prev => prev.filter((_, i) => i !== index));
  };

  return (
    <ErrorBoundary>
      <div className="app">
        <header className="app-header">
          <h1>Blue Marketplace Assistant</h1>
          {error && (
            <div className="error-banner" role="alert">
              {error}
              <button onClick={() => setError(null)} aria-label="Dismiss error">×</button>
            </div>
          )}
        </header>
        
        <main className="chat-container">
          <MessageList messages={messages} isLoading={isLoading} />
          
          <div className="input-section">
            <FileUploader 
              onAddAttachment={handleAddAttachment}
              disabled={isLoading}
            />
            <MessageInput
              onSendMessage={handleSendMessage}
              disabled={isLoading}
              attachments={attachments}
              onRemoveAttachment={handleRemoveAttachment}
            />
          </div>
        </main>
      </div>
    </ErrorBoundary>
  );
}

export default App;
