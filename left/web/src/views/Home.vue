<script lang="ts">
import { defineComponent } from "vue"
import api, { IMessage } from "../api"
import MessageCard from "../components/MessageCard.vue"

export default defineComponent({
  data() {
    return {
      loading: false,
      messages: [] as IMessage[],
      error: null as string | null,
      limits: [10, 20, 50, 100],
      has_back: false,
      back_cursor: 0,
      next_cursor: 0,
      has_next: false,
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
          this.has_back = res.data!.has_back;
          this.back_cursor = res.data!.back_cursor;
          this.next_cursor = res.data!.next_cursor;
          this.has_next = res.data!.has_next;
        } else {
          this.error = res.error!.message;
        }
      }
      catch (error: any) {
        this.error = error.message;
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
    <el-alert v-if="error" :title="error" type="error" effect="dark" :closable="false" />
    <el-space wrap>
      <el-space>
        <span>Limit</span>
        <el-select v-model="limit" placeholder="Limit" :disabled="loading">
          <el-option v-for="item in limits" :key="item" :label="item" :value="item" />
        </el-select>
      </el-space>
      <el-switch
        v-model="ascending"
        inactive-text="Newest"
        active-text="Oldest"
        :loading="loading"
      />
    </el-space>
    <el-space fill wrap :fill-ratio="20">
      <MessageCard :key="message.id" v-for="message of messages" :message="message" />
    </el-space>
    <el-button-group>
      <el-button type="primary" :disabled="!has_back" @click="firstPage" :loading="loading">First</el-button>
      <el-button type="primary" :disabled="!has_back" @click="backPage" :loading="loading">Previous</el-button>
      <el-button type="primary" :disabled="!has_next" @click="nextPage" :loading="loading">Next</el-button>
    </el-button-group>
    <el-backtop />
  </el-space>
</template>

<style scoped></style>