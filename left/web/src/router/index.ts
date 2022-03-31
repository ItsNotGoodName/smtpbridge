import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router"

import Home from "../views/Home.vue"

const routes = [
  { path: "/", name: "Home", component: Home },
  { path: "/message/:id(\\d+)", name: "Message", component: () => import("../views/Message.vue") },
] as RouteRecordRaw[]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router