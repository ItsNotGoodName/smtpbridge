<script lang="ts">
import { defineComponent } from "vue";
import { IMessage } from "../api"
import MessageDate from "./MessageDate.vue"

export default defineComponent({
  props: {
    message: {
      type: Object as () => IMessage,
      required: true,
    }
  },
  components: { MessageDate }
})
</script>

<template>
  <el-card :body-style="{ padding: '0px' }">
    <div class="flex h-20">
      <el-image
        v-if="message.attachments.length"
        class="aspect-square"
        :src="message.attachments[0].file"
      />
      <router-link class="flex-1 p-4" :to="{ name: 'Message', params: { id: message.id } }">
        <div>{{ message.subject }}</div>
        <el-space>
          <message-date :message="message" />
          <el-tag v-if="message.attachments.length">{{ message.attachments.length }}</el-tag>
        </el-space>
      </router-link>
    </div>
  </el-card>
</template>

<style scoped>
</style>