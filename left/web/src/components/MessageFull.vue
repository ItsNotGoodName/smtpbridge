<script lang="ts">
import { defineComponent } from "vue";
import { IMessage } from "../api"

export default defineComponent({
  props: {
    message: {
      type: Object as () => IMessage,
      required: true,
    }
  },
})
</script>

<template>
  <el-card :body-style="{ padding: '0px' }">{{ message }}</el-card>
</template>

<style scoped>
</style>