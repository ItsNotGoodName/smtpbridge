<script lang="ts">
import { defineComponent } from "vue";
import api, { IAttachment } from "../api"

export default defineComponent({
  props: {
    attachment: {
      type: Object as () => IAttachment,
      required: true,
    },
  },
  computed: {
    url() {
      return api.attachmentUrl(this.attachment)
    },
  }
})
</script>

<template>
  <el-image :src="url"></el-image>
</template>

<style scoped>
</style>