import { useState, useRef } from 'react';
import './FileUploader.css';

const FileUploader = ({ onAddAttachment, disabled }) => {
  const [dragActive, setDragActive] = useState(false);
  const [urlInput, setUrlInput] = useState('');
  const fileInputRef = useRef(null);

  const allowedFileTypes = [
    'image/jpeg', 'image/png', 'image/gif', 'image/webp',
    'text/plain', 'application/pdf',
    'application/json', 'text/csv'
  ];

  const validateFile = (file) => {
    if (!allowedFileTypes.includes(file.type)) {
      throw new Error(`Unsupported file type: ${file.type}`);
    }
    if (file.size > 10 * 1024 * 1024) { // 10MB limit
      throw new Error('File size must be less than 10MB');
    }
    return true;
  };

  const validateUrl = (url) => {
    try {
      new URL(url);
      return true;
    } catch {
      throw new Error('Please enter a valid URL');
    }
  };

  const handleFiles = (files) => {
    Array.from(files).forEach(file => {
      try {
        validateFile(file);
        onAddAttachment({ type: 'file', file });
      } catch (error) {
        alert(error.message);
      }
    });
  };

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (disabled) return;
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFiles(e.dataTransfer.files);
    }
  };

  const handleFileSelect = (e) => {
    if (e.target.files && e.target.files[0]) {
      handleFiles(e.target.files);
    }
  };

  const handleUrlSubmit = (e) => {
    e.preventDefault();
    if (!urlInput.trim()) return;
    
    try {
      validateUrl(urlInput.trim());
      onAddAttachment({ type: 'url', url: urlInput.trim() });
      setUrlInput('');
    } catch (error) {
      alert(error.message);
    }
  };

  return (
    <div className="file-uploader">
      <div
        className={`drop-zone ${dragActive ? 'active' : ''} ${disabled ? 'disabled' : ''}`}
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        role="button"
        tabIndex={disabled ? -1 : 0}
        aria-label="File upload area"
      >
        <input
          ref={fileInputRef}
          type="file"
          multiple
          onChange={handleFileSelect}
          disabled={disabled}
          className="file-input"
          accept={allowedFileTypes.join(',')}
          aria-label="Select files"
        />
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          disabled={disabled}
          className="upload-button"
        >
          📎 Attach Files
        </button>
        <p className="upload-hint">or drag and drop files here</p>
      </div>
      
      <form onSubmit={handleUrlSubmit} className="url-form">
        <input
          type="url"
          value={urlInput}
          onChange={(e) => setUrlInput(e.target.value)}
          placeholder="Paste a URL here..."
          disabled={disabled}
          className="url-input"
          aria-label="URL input"
        />
        <button
          type="submit"
          disabled={disabled || !urlInput.trim()}
          className="url-button"
        >
          Add URL
        </button>
      </form>
    </div>
  );
};

export default FileUploader;