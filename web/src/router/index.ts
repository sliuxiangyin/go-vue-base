import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import AboutView from '../views/AboutView.vue'
import ApiTestView from '../views/ApiTestView.vue'
import LoginView from '../views/LoginView.vue'
import DbAccountView from '../views/DbAccountView.vue'
import TestView from "@/views/TestView.vue";
import DatabaseChatView from "@/views/database/IndexView.vue";

 const routes = [
   {
     path: '/',
     redirect: '/test'
   },
  {
    path: "/test",
    component: TestView,
    children: [
      {
        path: '',
        name: 'Home',
        component: HomeView
      },
      {
        path: 'about',
        name: 'About',
        component: AboutView
      },
      {
        path: 'api-test',
        name: 'ApiTest',
        component: ApiTestView
      },
      {
        path: 'db-account',
        name: 'DbAccount',
        component: DbAccountView
      },
      {
        path: 'login',
        name: 'Login',
        component: LoginView
      }
    ]
  },
  {
    path: "/database-chat/:id",
    name: "databaseChat",
    component: DatabaseChatView,
  }
]


const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router