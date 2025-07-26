import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import FileUploader from '../FileUploader';

// Mock alert
global.alert = vi.fn();

describe('FileUploader', () => {
  const mockOnAddAttachment = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('renders file upload button and URL input', () => {
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={false} />);
    
    expect(screen.getByText('📎 Attach Files')).toBeInTheDocument();
    expect(screen.getByLabelText('URL input')).toBeInTheDocument();
    expect(screen.getByText('Add URL')).toBeInTheDocument();
  });

  test('handles file selection', async () => {
    const user = userEvent.setup();
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={false} />);
    
    const fileInput = screen.getByLabelText('Select files');
    await user.upload(fileInput, file);
    
    expect(mockOnAddAttachment).toHaveBeenCalledWith({
      type: 'file',
      file: file
    });
  });

  test('validates file types', async () => {
    const user = userEvent.setup();
    const invalidFile = new File(['test'], 'test.exe', { type: 'application/exe' });
    
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={false} />);
    
    const fileInput = screen.getByLabelText('Select files');
    await user.upload(fileInput, invalidFile);
    
    expect(global.alert).toHaveBeenCalledWith('Unsupported file type: application/exe');
    expect(mockOnAddAttachment).not.toHaveBeenCalled();
  });

  test('handles URL submission', async () => {
    const user = userEvent.setup();
    
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={false} />);
    
    const urlInput = screen.getByLabelText('URL input');
    const addButton = screen.getByText('Add URL');
    
    await user.type(urlInput, 'https://example.com');
    await user.click(addButton);
    
    expect(mockOnAddAttachment).toHaveBeenCalledWith({
      type: 'url',
      url: 'https://example.com'
    });
  });

  test('validates URLs', async () => {
    const user = userEvent.setup();
    
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={false} />);
    
    const urlInput = screen.getByLabelText('URL input');
    const addButton = screen.getByText('Add URL');
    
    await user.type(urlInput, 'invalid-url');
    await user.click(addButton);
    
    expect(global.alert).toHaveBeenCalledWith('Please enter a valid URL');
    expect(mockOnAddAttachment).not.toHaveBeenCalled();
  });

  test('handles drag and drop', async () => {
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={false} />);
    
    const dropZone = screen.getByRole('button', { name: 'File upload area' });
    
    fireEvent.dragEnter(dropZone);
    expect(dropZone).toHaveClass('active');
    
    fireEvent.drop(dropZone, {
      dataTransfer: { files: [file] }
    });
    
    expect(mockOnAddAttachment).toHaveBeenCalledWith({
      type: 'file',
      file: file
    });
  });

  test('disables components when disabled prop is true', () => {
    render(<FileUploader onAddAttachment={mockOnAddAttachment} disabled={true} />);
    
    expect(screen.getByText('📎 Attach Files')).toBeDisabled();
    expect(screen.getByLabelText('URL input')).toBeDisabled();
    expect(screen.getByText('Add URL')).toBeDisabled();
  });
});