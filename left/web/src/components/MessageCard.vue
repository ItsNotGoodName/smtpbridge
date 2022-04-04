<script lang="ts" setup>
import { computed } from "vue"
import { IMessage } from "../api"
import Date from "./Date.vue"

const { message } = defineProps({
  message: {
    type: Object as () => IMessage,
    required: true,
  },
})

const srcList = computed(() => {
  return message.attachments.map(a => a.url)
})
</script>

<template>
  <el-card shadow="hover" :body-style="{ padding: '0px' }">
    <div class="flex h-20">
      <el-image
        v-if="message.attachments.length"
        class="aspect-square"
        :preview-src-list="srcList"
        :src="message.attachments[0].url"
        fit="fill"
      />
      <router-link
        class="flex-1 p-4 text-left no-underline"
        :to="{ name: 'Message', params: { id: message.id } }"
      >
        <div>{{ message.subject }}</div>
        <el-space>
          <Date class="text-gray-400" :date="message.created_at" />
          <el-tag v-if="message.attachments.length">{{ message.attachments.length }}</el-tag>
        </el-space>
      </router-link>
    </div>
  </el-card>
</template>

<style scoped>
</style>