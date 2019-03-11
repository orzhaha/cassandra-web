<template lang="pug">
    div
      br
      el-button(
        type="danger"
        icon="el-icon-download")
        a(
          v-bind:href="href")
          span Download
</template>

<style>
a {
  color:#FFFFFF;
  text-decoration:none;
}
.CodeMirror {
  border: 1px solid #eee;
  height: auto;
}
</style>

<script>
import { domain } from '@/api/config'

export default {
  name: 'Export',
  data() {
    return {
      href: `${domain}/export?table=${this.$route.params.keyspace}.${this.$route.params.table}`,
    }
  },
};
</script>
