<script lang="ts">
import { defineComponent } from "vue"
import api, { IMessage } from "../api"
import MessageCard from "../components/MessageCard.vue"

export default defineComponent({
  data() {
    return {
      loading: false,
      message: null as IMessage | null,
      error: null as String | null,
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
      } finally {
        this.loading = false;
      }
    },
  },
  components: { MessageCard }
})
</script>

<template>
  <div v-if="loading">Loading...</div>
  <div v-if="error">{{ error }}</div>
  <MessageCard v-if="message" :message="message" />
</template>

<style lang="scss" scoped></style>