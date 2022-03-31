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
    subject() {
      return this.message.subject ? this.message.subject : "No Subject";
    },
    attachmentUrls() {
      let urls = []
      for (let attachment of this.message.attachments) {
        urls.push(api.attachmentUrl(attachment))
      }
      return urls
    }
  },
})
</script>

<template>
  <el-card :body-style="{ padding: '0px' }">
    <el-carousel v-if="attachmentUrls.length > 0">
      <el-carousel-item v-for="attachmentUrl, idx of attachmentUrls" :key="idx">
        <el-image :src="attachmentUrl" />
      </el-carousel-item>
    </el-carousel>
    <div style="padding: 14px">
      <span>{{ subject }}</span>
      <div class="bottom">
        <time class="time">{{ message.created_at }}</time>
      </div>
    </div>
  </el-card>
</template>

<style scoped>
.bottom {
  margin-top: 13px;
  line-height: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.time {
  font-size: 13px;
  color: #999;
}
</style>