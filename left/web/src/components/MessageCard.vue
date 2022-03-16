<script lang="ts">
import { defineComponent } from "vue";
import api, { IMessage } from "../api"

export default defineComponent({
  props: {
    message: {
      type: Object as () => IMessage,
      required: true,
    }
  },
  computed: {
    attachmentUrl() {
      if (this.message.attachments.length == 0) {
        return ""
      }
      return api.attachmentUrl(this.message.attachments[0])
    }
  },
})
</script>

<template>
  <el-card :body-style="{ padding: '0px' }">
    <el-image v-if="attachmentUrl" :src="attachmentUrl" />
    <div style="padding: 14px">
      <span>{{ message.subject }}</span>
      <div class="bottom">
        <time class="time">{{ message.created_at }}</time>
      </div>
    </div>
  </el-card>
</template>

<style lang="">
  
</style>