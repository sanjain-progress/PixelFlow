import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
const AUTH_URL = process.env.REACT_APP_AUTH_URL || 'http://localhost:50051';

// Create axios instance
const api = axios.create({
    baseURL: API_URL,
});

// Add interceptor to add token to requests
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

export const authService = {
    login: async (email, password) => {
        const response = await axios.post(`${AUTH_URL}/login`, { email, password });
        return response.data;
    },
    register: async (email, password) => {
        const response = await axios.post(`${AUTH_URL}/register`, { email, password });
        return response.data;
    },
    validate: async (token) => {
        const response = await axios.get(`${AUTH_URL}/validate`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        return response.data;
    }
};

export const taskService = {
    upload: async (imageUrl) => {
        const response = await api.post('/api/upload', { image_url: imageUrl });
        return response.data;
    },
    getAll: async () => {
        const response = await api.get('/api/tasks');
        return response.data;
    }
};

export default api;
