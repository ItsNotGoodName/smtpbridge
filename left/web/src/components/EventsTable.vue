<script lang="ts">
import { defineComponent } from "vue";
import { IEvent } from "../api";

export default defineComponent({
  props: {
    events: {
      type: Array as () => IEvent[],
      required: true,
    },
    loading: {
      type: Boolean,
    },
  }
})
</script>

<template>
  <el-table v-loading="loading" :data="events" class="w-full">
    <el-table-column prop="created_at" label="Date" width="180" />
    <el-table-column prop="name" label="Name" width="180" />
    <el-table-column prop="description" label="Description" />
  </el-table>
</template>

<style scoped>
</style>