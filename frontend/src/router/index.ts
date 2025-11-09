import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ArticleListView from '../views/ArticleListView.vue'
import ArticleDetailView from '../views/ArticleDetailView.vue'
import ArticleUploadView from '../views/ArticleUploadView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/articles',
      name: 'article-list',
      component: ArticleListView,
    },
    {
      path: '/articles/upload',
      name: 'article-upload',
      component: ArticleUploadView,
    },
    {
      path: '/articles/:id',
      name: 'article-detail',
      component: ArticleDetailView,
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

export default router
