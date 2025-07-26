import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import MessageInput from '../MessageInput';

describe('MessageInput', () => {
  const mockOnSendMessage = vi.fn();
  const mockOnRemoveAttachment = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('renders input field and send button', () => {
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={false}
        attachments={[]}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    expect(screen.getByLabelText('Message input')).toBeInTheDocument();
    expect(screen.getByLabelText('Send message')).toBeInTheDocument();
  });

  test('sends message on form submit', async () => {
    const user = userEvent.setup();
    
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={false}
        attachments={[]}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'Test message');
    await user.click(sendButton);
    
    expect(mockOnSendMessage).toHaveBeenCalledWith('Test message', []);
  });

  test('sends message on Enter key press', async () => {
    const user = userEvent.setup();
    
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={false}
        attachments={[]}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    const input = screen.getByLabelText('Message input');
    
    await user.type(input, 'Test message');
    await user.keyboard('{Enter}');
    
    expect(mockOnSendMessage).toHaveBeenCalledWith('Test message', []);
  });

  test('does not send empty message', async () => {
    const user = userEvent.setup();
    
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={false}
        attachments={[]}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    const sendButton = screen.getByLabelText('Send message');
    
    await user.click(sendButton);
    
    expect(mockOnSendMessage).not.toHaveBeenCalled();
  });

  test('displays attachments preview', () => {
    const attachments = [
      { type: 'file', file: { name: 'test.jpg' } },
      { type: 'url', url: 'https://example.com' }
    ];
    
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={false}
        attachments={attachments}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    expect(screen.getByText('test.jpg')).toBeInTheDocument();
    expect(screen.getByText('https://example.com')).toBeInTheDocument();
  });

  test('removes attachment when remove button clicked', async () => {
    const user = userEvent.setup();
    const attachments = [{ type: 'file', file: { name: 'test.jpg' } }];
    
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={false}
        attachments={attachments}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    const removeButton = screen.getByLabelText('Remove file');
    await user.click(removeButton);
    
    expect(mockOnRemoveAttachment).toHaveBeenCalledWith(0);
  });

  test('disables input when disabled prop is true', () => {
    render(
      <MessageInput 
        onSendMessage={mockOnSendMessage}
        disabled={true}
        attachments={[]}
        onRemoveAttachment={mockOnRemoveAttachment}
      />
    );
    
    expect(screen.getByLabelText('Message input')).toBeDisabled();
    expect(screen.getByLabelText('Send message')).toBeDisabled();
  });
});