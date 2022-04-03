<script lang="ts" setup>
import { RouterView } from "vue-router";
import Navbar from "./components/Navbar.vue";
</script>

<template>
  <el-container>
    <el-header>
      <navbar />
    </el-header>
    <el-main>
      <router-view />
    </el-main>
    <el-backtop />
  </el-container>
</template>

<style scoped></style>