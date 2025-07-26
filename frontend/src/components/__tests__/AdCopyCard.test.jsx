import { render, screen } from '@testing-library/react';
import AdCopyCard from '../AdCopyCard';

describe('AdCopyCard', () => {
  const mockAdCopy = {
    headlines: ['Great Product!', 'Amazing Deal!'],
    descriptions: ['High quality product', 'Best value for money'],
    call_to_action: 'Buy Now',
    target_segments: ['Tech Enthusiasts', 'Budget Shoppers']
  };

  test('renders all ad copy sections', () => {
    render(<AdCopyCard adCopy={mockAdCopy} />);
    
    expect(screen.getByText('Headlines')).toBeInTheDocument();
    expect(screen.getByText('Descriptions')).toBeInTheDocument();
    expect(screen.getByText('Call to Action')).toBeInTheDocument();
    expect(screen.getByText('Target Segments')).toBeInTheDocument();
  });

  test('renders headlines correctly', () => {
    render(<AdCopyCard adCopy={mockAdCopy} />);
    
    expect(screen.getByText('Great Product!')).toBeInTheDocument();
    expect(screen.getByText('Amazing Deal!')).toBeInTheDocument();
  });

  test('renders descriptions correctly', () => {
    render(<AdCopyCard adCopy={mockAdCopy} />);
    
    expect(screen.getByText('High quality product')).toBeInTheDocument();
    expect(screen.getByText('Best value for money')).toBeInTheDocument();
  });

  test('renders call to action', () => {
    render(<AdCopyCard adCopy={mockAdCopy} />);
    
    expect(screen.getByText('Buy Now')).toBeInTheDocument();
  });

  test('renders target segments', () => {
    render(<AdCopyCard adCopy={mockAdCopy} />);
    
    expect(screen.getByText('Tech Enthusiasts')).toBeInTheDocument();
    expect(screen.getByText('Budget Shoppers')).toBeInTheDocument();
  });

  test('handles missing optional fields', () => {
    const minimalAdCopy = {
      headlines: ['Test Headline'],
      descriptions: ['Test Description']
    };
    
    render(<AdCopyCard adCopy={minimalAdCopy} />);
    
    expect(screen.getByText('Test Headline')).toBeInTheDocument();
    expect(screen.getByText('Test Description')).toBeInTheDocument();
    expect(screen.queryByText('Call to Action')).not.toBeInTheDocument();
    expect(screen.queryByText('Target Segments')).not.toBeInTheDocument();
  });

  test('has proper accessibility attributes', () => {
    render(<AdCopyCard adCopy={mockAdCopy} />);
    
    expect(screen.getByRole('article')).toHaveAttribute('aria-label', 'Generated marketing copy');
  });
});