import axios from 'axios';

const api = axios.create({
    // This matches our Go server address
    baseURL: "http://localhost:8080/api"
});

// This interceptor automatically attaches our JWT token to every moment
api.interceptors.request.use((config) =>{
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error)
}
);

export default api;