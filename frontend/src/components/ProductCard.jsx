import './ProductCard.css';

const ProductCard = ({ product }) => {
  const formatPrice = (price) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(price);
  };

  return (
    <div className="product-card" role="article" aria-label={`Product: ${product.title}`}>
      <h3 className="product-title">{product.title}</h3>
      <div className="product-price">{formatPrice(product.price)}</div>
      <a
        href={product.link}
        target="_blank"
        rel="noopener noreferrer"
        className="product-link"
        aria-label={`View ${product.title} on marketplace`}
      >
        View Product →
      </a>
    </div>
  );
};

export default ProductCard;