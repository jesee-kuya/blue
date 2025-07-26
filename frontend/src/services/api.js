const API_BASE_URL = 'http://localhost:8080';

class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

const handleResponse = async (response) => {
  if (!response.ok) {
    if (response.status >= 500) {
      throw new ApiError('Service temporarily unavailable', response.status);
    } else if (response.status === 429) {
      throw new ApiError('Too many requests. Please try again later.', response.status);
    } else {
      throw new ApiError('Request failed', response.status);
    }
  }
  return response.json();
};

export const sendMessage = async (message, attachments = []) => {
  try {
    const formData = new FormData();
    formData.append('message', message);
    
    attachments.forEach((attachment, index) => {
      if (attachment.type === 'file') {
        formData.append(`file_${index}`, attachment.file);
      } else if (attachment.type === 'url') {
        formData.append(`url_${index}`, attachment.url);
      }
    });

    const response = await fetch(`${API_BASE_URL}/search`, {
      method: 'POST',
      body: formData,
    });

    return await handleResponse(response);
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError('Unable to connect to service', 0);
  }
};

export const getMarketingCopy = async (message) => {
  try {
    const response = await fetch(`${API_BASE_URL}/marketing`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message }),
    });

    return await handleResponse(response);
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError('Unable to connect to service', 0);
  }
};

export const checkHealth = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/health`);
    return await handleResponse(response);
  } catch (error) {
    throw new ApiError('Service unavailable', 0);
  }
};