<script lang="ts">
import { defineComponent } from "vue"
import api, { IMessage } from "../api"

export default defineComponent({
  data() {
    return {
      loading: false,
      message: null as IMessage | null,
      error: null as string | null,
    };
  },
  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        if (this.$route.name === "Message") {
          this.load()
        }
      },
      { immediate: true }
    )
  },
  methods: {
    async load() {
      if (this.loading) {
        return
      }

      this.loading = true;
      try {
        let res = await api.getMessage(this.$route.params.id as unknown as number)
        this.message = null
        this.error = null
        if (res.ok) {
          this.message = res.data!;
        } else {
          this.error = res.error!.message;
        }
      } catch (error: any) {
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
  },
})
</script>

<template>
  <el-alert v-if="loading" title="loading..." type="info" effect="dark" :closable="false" />
  <el-alert v-if="error" :title="error" type="error" effect="dark" :closable="false" />
  <message-full v-if="message" :message="message" />
</template>

<style scoped></style>