<script lang="ts" setup>
import { RouterView } from "vue-router";
import Navbar from "./components/Navbar.vue";
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <navbar />
      </el-header>
      <el-main>
        <router-view />
      </el-main>
      <el-backtop />
    </el-container>
  </div>
</template>

<style scoped></style>