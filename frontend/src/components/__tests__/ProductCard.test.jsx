import { render, screen } from '@testing-library/react';
import ProductCard from '../ProductCard';

describe('ProductCard', () => {
  const mockProduct = {
    title: 'Test Product',
    price: 29.99,
    link: 'https://example.com/product'
  };

  test('renders product information correctly', () => {
    render(<ProductCard product={mockProduct} />);
    
    expect(screen.getByText('Test Product')).toBeInTheDocument();
    expect(screen.getByText('$29.99')).toBeInTheDocument();
    expect(screen.getByText('View Product →')).toBeInTheDocument();
  });

  test('formats price correctly', () => {
    const product = { ...mockProduct, price: 1234.56 };
    render(<ProductCard product={product} />);
    
    expect(screen.getByText('$1,234.56')).toBeInTheDocument();
  });

  test('has correct link attributes', () => {
    render(<ProductCard product={mockProduct} />);
    
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', 'https://example.com/product');
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  test('has proper accessibility attributes', () => {
    render(<ProductCard product={mockProduct} />);
    
    expect(screen.getByRole('article')).toHaveAttribute('aria-label', 'Product: Test Product');
    expect(screen.getByRole('link')).toHaveAttribute('aria-label', 'View Test Product on marketplace');
  });
});