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
  <el-space fill class="w-full">
    <el-descriptions size="small" :title="'Message ' + message.id" :column="1" :border="true">
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">Subject</div>
        </template>
        {{ message.subject }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">From</div>
        </template>
        {{ message.from }}
      </el-descriptions-item>
      <el-descriptions-item v-for="to, index of message.to" :key="to">
        <template #label>
          <div class="cell-item">To #{{ index + 1 }}</div>
        </template>
        <div>{{ to }}</div>
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">Time</div>
        </template>
        {{ new Date(message.created_at).toLocaleString() }}
      </el-descriptions-item>
    </el-descriptions>
    <el-card>
      <template #header>
        <div class="text-md font-bold">Body</div>
      </template>
      <el-scrollbar>
        <code v-if="message.text">
          <pre class="m-0">{{ message.text }}</pre>
        </code>
      </el-scrollbar>
    </el-card>
    <el-card
      :body-style="{ padding: '0px' }"
      v-for="attachment in message.attachments"
      :key="attachment.id"
    >
      <el-image :src="attachment.file" :preview-src-list="[attachment.file]" />
      <el-scrollbar>
        <div class="p-2">{{ attachment.name }}</div>
      </el-scrollbar>
    </el-card>
  </el-space>
</template>

<style scoped>
</style>