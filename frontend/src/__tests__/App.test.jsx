import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import App from '../App';
import * as api from '../services/api';

// Mock the API module
vi.mock('../services/api');

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test('renders chat interface', () => {
    render(<App />);
    
    expect(screen.getByText('Blue Marketplace Assistant')).toBeInTheDocument();
    expect(screen.getByText(/Hello! I'm your Blue marketplace assistant/)).toBeInTheDocument();
    expect(screen.getByLabelText('Message input')).toBeInTheDocument();
  });

  test('sends message and displays response', async () => {
    const user = userEvent.setup();
    const mockResponse = {
      search_results: [
        { title: 'Test Product', price: 29.99, link: 'https://example.com' }
      ]
    };
    
    api.sendMessage.mockResolvedValueOnce(mockResponse);
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'search for laptops');
    await user.click(sendButton);
    
    expect(screen.getByText('search for laptops')).toBeInTheDocument();
    
    await waitFor(() => {
      expect(screen.getByText('Here are the products I found:')).toBeInTheDocument();
      expect(screen.getByText('Test Product')).toBeInTheDocument();
    });
  });

  test('handles marketing copy requests', async () => {
    const user = userEvent.setup();
    const mockResponse = {
      marketing_copy: {
        headlines: ['Great Product!'],
        descriptions: ['Amazing features'],
        call_to_action: 'Buy Now'
      }
    };
    
    api.getMarketingCopy.mockResolvedValueOnce(mockResponse);
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'create marketing copy for laptops');
    await user.click(sendButton);
    
    expect(screen.getByText('create marketing copy for laptops')).toBeInTheDocument();
    
    await waitFor(() => {
      expect(screen.getByText('Here\'s the marketing copy I generated:')).toBeInTheDocument();
      expect(screen.getByText('Great Product!')).toBeInTheDocument();
      expect(screen.getByText('Amazing features')).toBeInTheDocument();
      expect(screen.getByText('Buy Now')).toBeInTheDocument();
    });
  });

  test('shows loading state during API calls', async () => {
    const user = userEvent.setup();
    
    // Mock a delayed response
    api.sendMessage.mockImplementation(() => 
      new Promise(resolve => setTimeout(() => resolve({ message: 'Done' }), 100))
    );
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'test message');
    await user.click(sendButton);
    
    expect(screen.getByLabelText('Assistant is typing')).toBeInTheDocument();
    
    await waitFor(() => {
      expect(screen.queryByLabelText('Assistant is typing')).not.toBeInTheDocument();
    });
  });

  test('handles API errors gracefully', async () => {
    const user = userEvent.setup();
    
    api.sendMessage.mockRejectedValueOnce(new Error('Network error'));
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'test message');
    await user.click(sendButton);
    
    await waitFor(() => {
      expect(screen.getByText(/Sorry, I encountered an error/)).toBeInTheDocument();
      expect(screen.getByText('Network error')).toBeInTheDocument();
    });
  });

  test('displays error banner for API errors', async () => {
    const user = userEvent.setup();
    
    api.sendMessage.mockRejectedValueOnce(new Error('Service unavailable'));
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'test message');
    await user.click(sendButton);
    
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText('Service unavailable')).toBeInTheDocument();
    });
  });

  test('dismisses error banner when close button clicked', async () => {
    const user = userEvent.setup();
    
    api.sendMessage.mockRejectedValueOnce(new Error('Test error'));
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    const sendButton = screen.getByLabelText('Send message');
    
    await user.type(input, 'test message');
    await user.click(sendButton);
    
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
    
    const dismissButton = screen.getByLabelText('Dismiss error');
    await user.click(dismissButton);
    
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  test('handles file attachments', async () => {
    const user = userEvent.setup();
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    
    api.sendMessage.mockResolvedValueOnce({ message: 'File processed' });
    
    render(<App />);
    
    const fileInput = screen.getByLabelText('Select files');
    await user.upload(fileInput, file);
    
    expect(screen.getByText('test.jpg')).toBeInTheDocument();
    
    const sendButton = screen.getByLabelText('Send message');
    await user.click(sendButton);
    
    expect(api.sendMessage).toHaveBeenCalledWith('', [{ type: 'file', file }]);
  });

  test('handles URL attachments', async () => {
    const user = userEvent.setup();
    
    api.sendMessage.mockResolvedValueOnce({ message: 'URL processed' });
    
    render(<App />);
    
    const urlInput = screen.getByLabelText('URL input');
    const addUrlButton = screen.getByText('Add URL');
    
    await user.type(urlInput, 'https://example.com');
    await user.click(addUrlButton);
    
    expect(screen.getByText('https://example.com')).toBeInTheDocument();
    
    const sendButton = screen.getByLabelText('Send message');
    await user.click(sendButton);
    
    expect(api.sendMessage).toHaveBeenCalledWith('', [{ type: 'url', url: 'https://example.com' }]);
  });

  test('removes attachments when remove button clicked', async () => {
    const user = userEvent.setup();
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    
    render(<App />);
    
    const fileInput = screen.getByLabelText('Select files');
    await user.upload(fileInput, file);
    
    expect(screen.getByText('test.jpg')).toBeInTheDocument();
    
    const removeButton = screen.getByLabelText('Remove file');
    await user.click(removeButton);
    
    expect(screen.queryByText('test.jpg')).not.toBeInTheDocument();
  });

  test('clears attachments after sending message', async () => {
    const user = userEvent.setup();
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
    
    api.sendMessage.mockResolvedValueOnce({ message: 'Success' });
    
    render(<App />);
    
    const fileInput = screen.getByLabelText('Select files');
    await user.upload(fileInput, file);
    
    expect(screen.getByText('test.jpg')).toBeInTheDocument();
    
    const sendButton = screen.getByLabelText('Send message');
    await user.click(sendButton);
    
    await waitFor(() => {
      expect(screen.queryByText('test.jpg')).not.toBeInTheDocument();
    });
  });

  test('has proper accessibility structure', () => {
    render(<App />);
    
    expect(screen.getByRole('main')).toBeInTheDocument();
    expect(screen.getByRole('banner')).toBeInTheDocument();
    expect(screen.getByRole('log')).toBeInTheDocument();
  });

  test('supports keyboard navigation', async () => {
    const user = userEvent.setup();
    
    api.sendMessage.mockResolvedValueOnce({ message: 'Success' });
    
    render(<App />);
    
    const input = screen.getByLabelText('Message input');
    
    await user.type(input, 'test message');
    await user.keyboard('{Enter}');
    
    expect(api.sendMessage).toHaveBeenCalledWith('test message', []);
  });
});
