import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 second timeout
});

// Add request interceptor for debugging
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API Response Error:', error.response?.status, error.message);
    return Promise.reject(error);
  }
);

// Public APIs
export const getSessions = () => api.get('/api/sessions');
export const getSpeakers = () => api.get('/api/speakers');
export const getAttendeeCount = () => api.get('/api/attendees/count');
export const registerAttendee = (data) => api.post('/api/attendees/register', data);

// Admin APIs
export const adminLogin = (password) => api.post('/api/admin/login', { password });
export const getAttendees = () => api.get('/api/admin/attendees');
export const getAttendeeStats = () => api.get('/api/admin/stats');
export const createOrUpdateSpeaker = (data) => api.post('/api/admin/speakers', data);
export const createOrUpdateSession = (data) => api.post('/api/admin/sessions', data);

export default api;

