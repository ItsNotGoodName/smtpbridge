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
  computed: {
    srcList() {
      return this.message.attachments.map(a => a.file)
    }
  },
  components: { MessageDate }
})
</script>

<template>
  <el-card shadow="hover" :body-style="{ padding: '0px' }">
    <div class="flex h-20">
      <el-image
        v-if="message.attachments.length"
        class="aspect-square"
        :preview-src-list="srcList"
        :src="message.attachments[0].file"
        fit="fill"
      />
      <router-link
        class="flex-1 p-4 text-left no-underline"
        :to="{ name: 'Message', params: { id: message.id } }"
      >
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