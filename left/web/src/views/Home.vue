<script lang="ts">
import { defineComponent } from "vue"
import api, { IResponse, IMessage, IMessages } from "../api"
import MessageCard from "../components/MessageCard.vue"

export default defineComponent({
  data() {
    return {
      next_cursor: 0,
      has_more: false,
      messages: [] as IMessage[],
      loading: false,
    };
  },
  async beforeMount() {
    this.handleResponse(await api.getMessages({}));
  },
  methods: {
    handleResponse(res: IResponse<IMessages>) {
      if (res.ok) {
        this.messages = this.messages.concat(res.data!.messages);
        this.has_more = res.data!.has_more;
        this.next_cursor = res.data!.next_cursor;
      }
    },
    async load() {
      if (this.has_more && !this.loading) {
        this.loading = true;
        this.handleResponse(await api.getMessages({ cursor: this.next_cursor }));
        this.loading = false;
      }
    }
  },
  components: { MessageCard }
})
</script>

<template>
  <el-space wrap fill>
    <MessageCard :key="message.id" v-for="message of messages" :message="message" />
    <el-button v-if="has_more" type="text" @click="load" :loading="loading">Load more</el-button>
  </el-space>
</template>

<style lang="scss" scoped></style>