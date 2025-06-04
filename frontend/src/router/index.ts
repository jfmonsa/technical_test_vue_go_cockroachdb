import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/HomeView.vue';
import Recommendations from '../views/RecommendationsView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/recommendations',
      name: 'recommendations',
      component: Recommendations,
    },
  ],
});

export default router;
