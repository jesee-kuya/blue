# Blue Marketplace Assistant - Frontend

A React-based chat interface for the Blue marketplace assistant that provides product search, marketing copy generation, and file/URL attachment capabilities.

## Features

### Core Components
- **MessageList**: Displays conversation history with support for text, product listings, and ad copy
- **MessageInput**: Text input with send functionality and keyboard shortcuts
- **FileUploader**: Drag-and-drop file upload and URL attachment support
- **ProductCard**: Displays search results in card format
- **AdCopyCard**: Shows generated marketing copy with organized sections
- **ErrorBoundary**: Graceful error handling for component failures

### Functionality
- Real-time chat interface with auto-scroll
- File upload with type validation (images, PDFs, text files)
- URL attachment support
- Responsive design for mobile and desktop
- Accessibility features (screen reader support, keyboard navigation)
- Loading states and error handling
- Integration with Go backend API

## Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Run tests with coverage
npm run coverage

# Build for production
npm run build
```

### Environment Setup

The frontend expects the backend API to be running on `http://localhost:8080`. Update the `API_BASE_URL` in `src/services/api.js` if your backend runs on a different port.

## Project Structure

```
frontend/
├── src/
│   ├── components/           # React components
│   │   ├── __tests__/       # Component tests
│   │   ├── MessageList.jsx
│   │   ├── MessageInput.jsx
│   │   ├── FileUploader.jsx
│   │   ├── ProductCard.jsx
│   │   ├── AdCopyCard.jsx
│   │   └── ErrorBoundary.jsx
│   ├── services/            # API services
│   │   ├── __tests__/
│   │   └── api.js
│   ├── test/               # Test setup
│   │   └── setup.js
│   ├── __tests__/          # App-level tests
│   ├── App.jsx
│   ├── App.css
│   └── main.jsx
├── vitest.config.js        # Test configuration
└── package.json
```

## API Integration

The frontend integrates with the following backend endpoints:

- `POST /search` - Product search with file/URL attachments
- `POST /marketing` - Marketing copy generation
- `GET /health` - Health check

### Error Handling

The application handles various error scenarios:
- Network connectivity issues
- API server errors (5xx)
- Rate limiting (429)
- Invalid file uploads
- Service unavailability

## Testing

The project uses Vitest and React Testing Library for testing:

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:ui

# Generate coverage report
npm run coverage
```

### Test Coverage Requirements
- Minimum 80% coverage for all metrics (lines, functions, branches, statements)
- All interactive components tested
- API integration mocked and tested
- Accessibility features tested
- Error scenarios covered

### Test Structure
- Unit tests for individual components
- Integration tests for component interactions
- API service tests with mocked responses
- Accessibility tests using screen reader queries

## Accessibility Features

- Semantic HTML structure with proper ARIA labels
- Keyboard navigation support
- Screen reader compatibility
- Focus management
- High contrast support
- Responsive design for various screen sizes

## File Upload Support

Supported file types:
- Images: JPEG, PNG, GIF, WebP
- Documents: PDF, TXT
- Data: JSON, CSV

File size limit: 10MB per file

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Development

### Code Style
- ESLint configuration in `eslint.config.js`
- Prettier for code formatting
- React hooks for state management
- CSS modules for styling

### Performance Considerations
- Lazy loading for large file previews
- Debounced API calls
- Optimized re-renders with React.memo where appropriate
- Efficient message list virtualization for large conversations

## Deployment

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

The build output will be in the `dist/` directory, ready for deployment to any static hosting service.

## Contributing

1. Follow the existing code style and patterns
2. Write tests for new features
3. Ensure accessibility compliance
4. Update documentation as needed
5. Test across different browsers and screen sizes

## Troubleshooting

### Common Issues

**Tests failing with "Cannot find module" errors:**
- Ensure all dependencies are installed: `npm install`
- Check that test setup file is properly configured

**File uploads not working:**
- Verify file types are in the allowed list
- Check file size limits
- Ensure backend API is running and accessible

**API calls failing:**
- Check backend server is running on correct port
- Verify CORS configuration on backend
- Check network connectivity

**Styling issues:**
- Ensure CSS files are properly imported
- Check for conflicting styles
- Verify responsive breakpoints
