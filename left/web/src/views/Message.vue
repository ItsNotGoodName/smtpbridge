<script lang="ts" setup>
import { ref, watch, computed } from "vue"
import { useRoute, useRouter } from "vue-router";
import api from "../api"
import { useFetch } from "../fetch"
import EventsTable from "../components/EventsTable.vue";

const route = useRoute()

const id = ref(0)
const page = ref(1)
const limit = ref(1)
const ascending = ref(false)
const maxPage = ref(1)

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
} = useFetch(computed(() => api.messageEventsGet(id.value, { page: page.value, limit: limit.value, ascending: ascending.value })), { skip: true })

watch(() => route.params, () => {
  if (route.name === "Message") {
    id.value = parseInt(route.params.id as string)
    fetchMessage()
    fetchEvents()
  }
}, { immediate: true })

watch(() => events.value, () => {
  if (events.value) {
    page.value = events.value.page
    maxPage.value = events.value.max_page
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
          layout="prev, pager, next"
          v-model:currentPage="page"
          v-model:page-size="limit"
          :page-count="maxPage"
          @current-change="fetchEvents"
        />
      </el-card>
    </template>
  </el-space>
</template>

<style scoped></style>