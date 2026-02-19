import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import PropertySelector from "../views/PropertySelector.vue";
import UnifiedLayout from "../layouts/UnifiedLayout.vue";
import DashboardOverview from "../views/DashboardOverview.vue";
import Rooms from "../views/Rooms.vue";
import Tenants from "../views/Tenants.vue";
import Complaints from "../views/Complaints.vue";
import Admission from "../views/Admission.vue";
// ADD THIS IMPORT:
import PublicComplaints from "../views/PublicComplaints.vue"; 

const routes = [
  { path: '/login', name: 'Login', component: Login },
  {
    path: '/',
    component: UnifiedLayout,
    children: [
      { path: 'property-selector', name: 'PropertySelector', component: PropertySelector },
      { path: 'property/:id/dashboard', name: 'Dashboard', component: DashboardOverview },
      { path: 'property/:id/inventory', name: 'Rooms', component: Rooms },
      { path: 'property/:id/tenants', name: 'Tenants', component: Tenants },
      { path: 'property/:id/complaints', name: 'Complaints', component: Complaints },
      { path: 'property/:id/admission', name: 'Admission', component: Admission },
      { path: 'property/:id/expenditure', name: 'Expenditure', component: () => import('../views/Expenditure.vue') },
      { path: 'property/:id/payments', name: 'PaymentHistory', component: () => import('../views/PaymentHistory.vue') },
      { path: 'property/:id/archives', name: 'Archives', component: () => import('../views/ArchivedTenants.vue') },
      { path: 'profile-settings', name: 'ProfileSettings', component: () => import('../views/ProfileSettings.vue') },
    ]
  },
  { path: '/help/:id', name: 'PublicHelp', component: PublicComplaints },
  // Move redirect to the bottom to avoid catching the /help route
  { path: '/:pathMatch(.*)*', redirect: '/login' } 
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// CRITICAL FIX: Route Guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  
  // Allow access to Login AND the Public Help page without a token
  if (to.name !== 'Login' && to.name !== 'PublicHelp' && !token) {
    next({ name: 'Login' });
  } else {
    next();
  }
});

export default router;