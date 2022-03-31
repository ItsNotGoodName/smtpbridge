<script lang="ts">
import { defineComponent } from "vue";
import { IMessage } from "../api"

export default defineComponent({
  props: {
    message: {
      type: Object as () => IMessage,
      required: true,
    }
  },
  computed: {
    date() {
      return new Date(this.message.created_at).toLocaleString();
    },
  }
})
</script>

<template>
  <time>{{ date }}</time>
</template>

<style scoped>
time {
  color: #999;
}
</style>