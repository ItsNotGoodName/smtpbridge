<script lang="ts" setup>
import { ref, watch, computed } from "vue"
import { useRoute } from "vue-router";
import api from "../api"
import { useFetch } from "../fetch"

const route = useRoute()

const id = ref(0)
const eventPage = ref(1)
const eventLimit = ref(10)
const eventAscending = ref(false)
const eventMaxPage = ref(1)

const {
  data: message,
  error: messageError,
  loading: messageLoading,
  fetch: fetchMessage
} = useFetch(computed(() => api.messageGet(id.value)), { skip: true })
const {
  data: events,
  error: eventsError,
  loading: eventsLoading,
  fetch: fetchEvents
} = useFetch(computed(() => api.messageEventsGet(id.value, { page: eventPage.value, limit: eventLimit.value, ascending: eventAscending.value })), { skip: true })

watch(() => route.params, () => {
  if (route.name === "Message") {
    id.value = parseInt(route.params.id as string)
    fetchMessage()
    fetchEvents()
  }
}, { immediate: true })

watch(() => events.value, () => {
  if (events.value) {
    eventPage.value = events.value.page
    eventMaxPage.value = events.value.max_page
  }
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
      <el-alert
        v-if="eventsError"
        :title="eventsError"
        type="error"
        effect="dark"
        :closable="false"
      />
      <el-card :body-style="{ padding: '0px' }">
        <template #header>
          <div class="text-md font-bold">Events</div>
        </template>
        <events-table :loading="eventsLoading" v-if="events" :events="events.events" />
        <el-pagination
          layout="sizes, prev, pager, next"
          v-model:currentPage="eventPage"
          v-model:page-size="eventLimit"
          :page-count="eventMaxPage"
          @current-change="fetchEvents"
          @size-change="fetchEvents"
        />
      </el-card>
    </template>
  </el-space>
</template>

<style scoped></style>