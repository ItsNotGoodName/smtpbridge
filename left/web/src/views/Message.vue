<script lang="ts">
import { defineComponent } from "vue"
import api, { IEvent, IMessage, IPage } from "../api"

export default defineComponent({
  data() {
    return {
      messageLoading: false,
      message: null as IMessage | null,
      messageError: null as string | null,
      eventsLoading: false,
      events: [] as IEvent[],
      eventsError: null as string | null,
      eventMaxPage: 1,
    };
  },
  created() {
    this.$watch(() => this.$route.params, () => {
      if (this.$route.name === "Message") {
        this.loadMessage();
        this.loadEvents();
      }
    }, { immediate: true });
  },
  computed: {
    id(): number {
      return this.$route.params.id as unknown as number;
    },
  },
  methods: {
    async loadMessage() {
      if (this.messageLoading) {
        return;
      }
      this.messageLoading = true;
      try {
        let res = await api.getMessage(this.id);
        this.message = null;
        this.messageError = null;
        if (res.ok) {
          this.message = res.data!;
        }
        else {
          this.messageError = res.error!.message;
        }
      }
      catch (error: any) {
        this.messageError = error.message;
      }
      finally {
        this.messageLoading = false;
      }
    },
    async loadEvents() {
      if (this.eventsLoading) {
        return;
      }
      this.eventsLoading = true;
      try {
        let res = await api.getMessageEvents(this.id, {} as IPage);
        this.events = [];
        this.eventsError = null;
        if (res.ok) {
          this.events = res.data!.events;
          this.eventMaxPage = res.data!.max_page;
        }
        else {
          this.eventsError = res.error!.message;
        }
      }
      catch (error: any) {
        this.eventsError = error.message;
      }
      finally {
        this.eventsLoading = false;
      }
    },
  },
})
</script>

<template>
  <el-space fill class="w-full">
    <el-alert v-if="messageLoading" title="loading..." type="info" effect="dark" :closable="false" />
    <el-alert
      v-if="messageError"
      :title="messageError"
      type="error"
      effect="dark"
      :closable="false"
    />
    <template v-if="message">
      <message-full :message="message" />
      <events-table :loading="eventsLoading" :events="events" />
    </template>
  </el-space>
</template>

<style scoped></style>