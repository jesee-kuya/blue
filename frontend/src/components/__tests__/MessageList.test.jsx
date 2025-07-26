import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import MessageList from '../MessageList';

const mockMessages = [
  {
    id: 1,
    sender: 'user',
    text: 'Hello',
    timestamp: Date.now(),
    type: 'text'
  },
  {
    id: 2,
    sender: 'assistant',
    text: 'Hi there!',
    timestamp: Date.now(),
    type: 'text'
  }
];

describe('MessageList', () => {
  test('renders messages correctly', () => {
    render(<MessageList messages={mockMessages} isLoading={false} />);
    
    expect(screen.getByText('Hello')).toBeInTheDocument();
    expect(screen.getByText('Hi there!')).toBeInTheDocument();
    expect(screen.getByText('You')).toBeInTheDocument();
    expect(screen.getByText('Assistant')).toBeInTheDocument();
  });

  test('shows loading indicator when loading', () => {
    render(<MessageList messages={[]} isLoading={true} />);
    
    expect(screen.getByRole('status')).toBeInTheDocument();
    expect(screen.getByLabelText('Assistant is typing')).toBeInTheDocument();
  });

  test('renders product cards for search results', () => {
    const searchMessage = {
      id: 3,
      sender: 'assistant',
      text: 'Here are the products:',
      timestamp: Date.now(),
      type: 'search_results',
      data: {
        products: [
          { title: 'Test Product', price: 29.99, link: 'https://example.com' }
        ]
      }
    };

    render(<MessageList messages={[searchMessage]} isLoading={false} />);
    
    expect(screen.getByText('Test Product')).toBeInTheDocument();
    expect(screen.getByText('$29.99')).toBeInTheDocument();
  });

  test('renders ad copy cards for marketing results', () => {
    const marketingMessage = {
      id: 4,
      sender: 'assistant',
      text: 'Here is your ad copy:',
      timestamp: Date.now(),
      type: 'marketing_copy',
      data: {
        marketing: {
          headlines: ['Great Product!'],
          descriptions: ['Amazing features'],
          call_to_action: 'Buy Now'
        }
      }
    };

    render(<MessageList messages={[marketingMessage]} isLoading={false} />);
    
    expect(screen.getByText('Great Product!')).toBeInTheDocument();
    expect(screen.getByText('Amazing features')).toBeInTheDocument();
    expect(screen.getByText('Buy Now')).toBeInTheDocument();
  });

  test('has proper accessibility attributes', () => {
    render(<MessageList messages={mockMessages} isLoading={false} />);
    
    expect(screen.getByRole('log')).toBeInTheDocument();
    expect(screen.getAllByRole('article')).toHaveLength(2);
  });
});