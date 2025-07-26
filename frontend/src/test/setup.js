import '@testing-library/jest-dom';
import { vi } from 'vitest';

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn(() => ({
  disconnect: vi.fn(),
  observe: vi.fn(),
  unobserve: vi.fn(),
}));

// Mock ResizeObserver
global.ResizeObserver = vi.fn(() => ({
  disconnect: vi.fn(),
  observe: vi.fn(),
  unobserve: vi.fn(),
}));

// Mock scrollIntoView
Element.prototype.scrollIntoView = vi.fn();

// Mock URL constructor for file upload tests
global.URL = global.URL || {
  createObjectURL: vi.fn(() => 'mocked-url'),
  revokeObjectURL: vi.fn(),
};

// Mock FormData for file upload tests
global.FormData = global.FormData || class FormData {
  constructor() {
    this.data = new Map();
  }
  
  append(key, value) {
    this.data.set(key, value);
  }
  
  get(key) {
    return this.data.get(key);
  }
  
  has(key) {
    return this.data.has(key);
  }
};