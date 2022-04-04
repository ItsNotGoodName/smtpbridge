<script lang="ts" setup>
import { IEvent } from "../api";
import Date from "./Date.vue"

const { events, loading } = defineProps({
  events: {
    type: Object as () => IEvent[],
    required: true,
  },
  loading: {
    type: Boolean,
  },
});
</script>

<template>
  <el-table v-loading="loading" :data="events" class="w-full">
    <el-table-column prop="created_at" label="Date" width="180">
      <template #default="scope">
        <Date :date="scope.row.created_at" />
      </template>
    </el-table-column>
    <el-table-column prop="name" label="Name" width="180" />
    <el-table-column prop="description" label="Description" />
  </el-table>
</template>

<style scoped>
</style>