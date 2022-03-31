<script lang="ts">
import { defineComponent } from "vue";
import { IMessage } from "../api"
import MessageDate from "./MessageDate.vue"
import AttachmentImage from "./AttachmentImage.vue";

export default defineComponent({
  props: {
    message: {
      type: Object as () => IMessage,
      required: true,
    }
  },
  methods: {
    onClick() {
      this.$router.push({
        name: "Message",
        params: {
          id: this.message.id,
        },
      });
    }
  },
  components: { MessageDate, AttachmentImage }
})
</script>

<template>
  <el-card :body-style="{ padding: '0px' }" @click="onClick">
    <div class="card">
      <attachment-image
        v-if="message.attachments.length"
        :attachment="message.attachments[0]"
        class="thumbnail"
      />
      <div class="body">
        <span>{{ message.subject }}</span>
        <div class="body-bottom">
          <message-date :message="message" />
          <el-tag v-if="message.attachments.length">{{ message.attachments.length }}</el-tag>
        </div>
      </div>
    </div>
  </el-card>
</template>

<style scoped>
.card {
  display: flex;
  height: 5rem;
}
.thumbnail {
  width: 20%;
}
.body {
  padding: 14px;
}
.body-bottom {
  margin-top: 13px;
  line-height: 12px;
  display: flex;
  align-items: center;
}
</style>