import { vi, beforeEach, afterEach } from 'vitest';
import { sendMessage, getMarketingCopy, checkHealth } from '../api';

// Mock fetch
global.fetch = vi.fn();

describe('API Service', () => {
  beforeEach(() => {
    fetch.mockClear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('sendMessage', () => {
    test('sends message successfully', async () => {
      const mockResponse = { search_results: [] };
      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await sendMessage('test message');
      
      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/search', {
        method: 'POST',
        body: expect.any(FormData),
      });
      expect(result).toEqual(mockResponse);
    });

    test('handles network errors', async () => {
      fetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(sendMessage('test')).rejects.toThrow('Unable to connect to service');
    });

    test('handles server errors', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
      });

      await expect(sendMessage('test')).rejects.toThrow('Service temporarily unavailable');
    });

    test('handles rate limiting', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 429,
      });

      await expect(sendMessage('test')).rejects.toThrow('Too many requests. Please try again later.');
    });
  });

  describe('getMarketingCopy', () => {
    test('gets marketing copy successfully', async () => {
      const mockResponse = { marketing_copy: {} };
      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await getMarketingCopy('test message');
      
      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/marketing', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message: 'test message' }),
      });
      expect(result).toEqual(mockResponse);
    });
  });

  describe('checkHealth', () => {
    test('checks health successfully', async () => {
      const mockResponse = { status: 'ok' };
      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await checkHealth();
      
      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/health');
      expect(result).toEqual(mockResponse);
    });
  });
});