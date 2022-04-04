<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    date: {
      type: String,
      required: true,
    }
  },
  computed: {
    parsedDate() {
      return new Date(this.date).toLocaleString();
    },
  }
})
</script>

<template>
  <time>{{ parsedDate }}</time>
</template>

<style scoped>
</style>