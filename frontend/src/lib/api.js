import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default {
  // Auth
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  getUser: () => api.get('/user'),
  getDashboardStats: () => api.get('/user/stats'),

  // API Keys
  getAPIKeys: () => api.get('/api-keys'),
  createAPIKey: (data) => api.post('/api-keys', data),
  revokeAPIKey: (id) => api.delete(`/api-keys/${id}`),

  // Domains
  getDomains: () => api.get('/domains'),
  addDomain: (data) => api.post('/domains', data),
  getDomain: (id) => api.get(`/domains/${id}`),
  verifyDomain: (id) => api.post(`/domains/${id}/verify`),
  deleteDomain: (id) => api.delete(`/domains/${id}`),

  // Templates
  getTemplates: () => api.get('/templates'),
  createTemplate: (data) => api.post('/templates', data),
  getTemplate: (id) => api.get(`/templates/${id}`),
  updateTemplate: (id, data) => api.put(`/templates/${id}`, data),
  deleteTemplate: (id) => api.delete(`/templates/${id}`),

  // Emails
  sendEmail: (data) => api.post('/emails/send', data),
  getEmailLogs: (params) => api.get('/emails/logs', { params }),

  // Analytics
  getAnalytics: (params) => api.get('/analytics', { params }),

  // Billing
  getTransactions: (params) => api.get('/billing/transactions', { params }),
  initiateTopup: (data) => api.post('/billing/topup', data),
};
