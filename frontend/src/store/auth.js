import { defineStore } from 'pinia' ;
import api from '../api';
export const useAuthStore = defineStore('auth', {
    state: () => ({
        token: localStorage.getItem('token') || null,
        owner: null,
    }),
    actions: {
        async login(credentials){
            try{
                const response = await api.post('/auth/login', credentials);
                this.token = response.data.token;
                localStorage.setItem('token',this.token);
                return true;
            }
            catch (error){
                console.error("Login failed",error);
                return false;
            }
        },
        logout(){
            this.token=null;
            localStorage.removeItem('token');
        }
    }
}); 