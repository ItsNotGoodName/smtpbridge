<script lang="ts">
import { defineComponent } from "vue"
import api, { IMessage } from "../api"
import MessageCard from "../components/MessageCard.vue"

export default defineComponent({
  data() {
    return {
      limits: [10, 20, 50, 100],
      back_cursor: 0,
      next_cursor: 0,
      loading: false,
      messages: [] as IMessage[],
    };
  },
  created() {
    this.$watch(
      () => this.$route.query,
      () => {
        if (this.$route.name === "Home") {
          this.load();
        }
      },
      { immediate: true }
    )
  },
  computed: {
    limit: {
      get(): number {
        let limit = parseInt(this.$route.query.limit as string);
        return isNaN(limit) ? this.limits[0] : limit;
      },
      set(value) {
        this.push({ limit: value })
      },
    },
    cursor: {
      get() {
        let cursor = parseInt(this.$route.query.cursor as string);
        return isNaN(cursor) ? 0 : cursor;
      },
      set(value) {
        this.push({ cursor: value })
      },
    },
    ascending: {
      get() {
        return (this.$route.query.ascending as string) == "true"
      },
      set(value) {
        this.push({ ascending: value, cursor: 0 })
      },
    },
    backPageDisabled(): boolean {
      return this.back_cursor == this.cursor || this.back_cursor == 0
    },
    nextPageDisabled(): boolean {
      return this.next_cursor == this.cursor || this.next_cursor == 9223372036854776000
    },
  },
  methods: {
    async load() {
      if (this.loading) {
        return
      }

      this.loading = true;
      try {
        let res = await api.getMessages({ cursor: this.cursor, ascending: this.ascending, limit: this.limit })
        if (res.ok) {
          this.messages = res.data!.messages;
          this.back_cursor = res.data!.back_cursor;
          this.next_cursor = res.data!.next_cursor;
        }
      } finally {
        this.loading = false;
      }
    },
    push(query: {}) {
      return this.$router.push({ name: "Home", query: { ...this.$route.query, ...query, } })
    },
    firstPage() {
      return this.push({ cursor: 0 })
    },
    backPage() {
      return this.push({ cursor: this.back_cursor })
    },
    nextPage() {
      return this.push({ cursor: this.next_cursor })
    },
  },
  components: { MessageCard }
})
</script>

<template>
  <el-space fill>
    <el-space wrap>
      <el-space>
        <span>Limit</span>
        <el-select v-model="limit" placeholder="Limit" :disabled="loading">
          <el-option v-for="item in limits" :key="item" :label="item" :value="item" />
        </el-select>
      </el-space>
      <el-switch v-model="ascending" inactive-text="Newest" active-text="Oldest" :loading="loading" />
    </el-space>
    <el-space fill wrap :fill-ratio="20">
      <MessageCard :key="message.id" v-for="message of messages" :message="message" />
    </el-space>
    <el-button-group>
      <el-button type="primary" :disabled="backPageDisabled" @click="firstPage" :loading="loading">First</el-button>
      <el-button type="primary" :disabled="backPageDisabled" @click="backPage" :loading="loading">Previous</el-button>
      <el-button type="primary" :disabled="nextPageDisabled" @click="nextPage" :loading="loading">Next</el-button>
    </el-button-group>
  </el-space>
</template>

<style scoped></style>