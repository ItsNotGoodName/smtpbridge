<script lang="ts">
import { defineComponent } from "vue"
import api, { IResponse, IMessage, IMessages } from "../api"
import MessageCard from "../components/MessageCard.vue"

export default defineComponent({
  data() {
    return {
      limit: 10,
      limits: [10, 20, 50, 100],
      ascending: false,
      cursor: 0,
      back_cursor: 0,
      next_cursor: 0,
      loading: false,
      messages: [] as IMessage[],
    };
  },
  beforeMount() {
    this.firstPage()
  },
  methods: {
    async load(cursor: number) {
      if (this.loading) {
        return
      }

      this.loading = true;
      try {
        let res = await api.getMessages({ cursor: cursor, ascending: this.ascending, limit: this.limit })
        if (res.ok) {
          this.cursor = cursor
          this.messages = res.data!.messages;
          this.back_cursor = res.data!.back_cursor;
          this.next_cursor = res.data!.next_cursor;
        }
      } finally {
        this.loading = false;
      }
    },
    refreshPage() {
      return this.load(this.cursor);
    },
    firstPage() {
      return this.load(0);
    },
    backPage() {
      return this.load(this.back_cursor);
    },
    nextPage() {
      return this.load(this.next_cursor);
    },
  },
  components: { MessageCard }
})
</script>

<template>
  <el-space fill style="width: 100%">
    <el-space>
      <span>Limit</span>
      <el-select v-model="limit" placeholder="Limit" @change="refreshPage" :disabled="loading">
        <el-option v-for="item in limits" :key="item" :label="item" :value="item" />
      </el-select>
      <el-switch
        v-model="ascending"
        inactive-text="Newest"
        active-text="Oldest"
        :loading="loading"
        @click="firstPage"
      />
    </el-space>
    <el-space fill wrap :fill-ratio="20">
      <MessageCard :key="message.id" v-for="message of messages" :message="message" />
    </el-space>
    <el-button-group>
      <el-button
        type="primary"
        :disabled="cursor == back_cursor"
        @click="firstPage"
        :loading="loading"
      >First Page</el-button>
      <el-button
        type="primary"
        :disabled="cursor == back_cursor"
        @click="backPage"
        :loading="loading"
      >Previous Page</el-button>
      <el-button
        type="primary"
        :disabled="cursor == next_cursor"
        @click="nextPage"
        :loading="loading"
      >Next Page</el-button>
    </el-button-group>
  </el-space>
</template>

<style lang="scss" scoped></style>