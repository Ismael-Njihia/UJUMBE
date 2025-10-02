import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add API key to requests
export function setApiKey(apiKey) {
  if (apiKey) {
    api.defaults.headers.common['X-API-Key'] = apiKey;
  } else {
    delete api.defaults.headers.common['X-API-Key'];
  }
}

// Auth API
export const authAPI = {
  register: (email, password) => api.post('/auth/register', { email, password }),
  login: (email, password) => api.post('/auth/login', { email, password }),
};

// Email API
export const emailAPI = {
  send: (data) => api.post('/emails/send', data),
  getStatus: (id) => api.get(`/emails/${id}`),
  getLogs: (id) => api.get(`/emails/${id}/logs`),
  getAnalytics: () => api.get('/analytics'),
  getQuota: () => api.get('/quota'),
};

// Template API
export const templateAPI = {
  create: (data) => api.post('/templates', data),
  list: () => api.get('/templates'),
  get: (id) => api.get(`/templates/${id}`),
  delete: (id) => api.delete(`/templates/${id}`),
};

// Domain API
export const domainAPI = {
  add: (domain) => api.post('/domains', { domain }),
  list: () => api.get('/domains'),
  verify: (id) => api.post(`/domains/${id}/verify`),
  delete: (id) => api.delete(`/domains/${id}`),
};
