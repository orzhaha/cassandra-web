<template lang="pug">
  div
    br
    el-upload(
      class="upload-demo"
      ref="upload"
      v-bind:action="href"
      :file-list="fileList"
      :auto-upload="false")
      el-button(
          slot="trigger"
          type="primary"
          icon="el-icon-upload2") Select File

      el-button(
          type="danger"
          icon="el-icon-upload2"
          class="ml10"
          @click="submitUpload") Upload
</template>

<style>
.ml10 {
  margin-left: 10px;
}
.CodeMirror {
  border: 1px solid #eee;
  height: auto;
}
</style>

<script>
import { domain } from '@/api/config'

export default {
  name: 'Import',
  data() {
    return {
      fileList: [],
      href: `${domain}/import?table=${this.$route.params.keyspace}.${this.$route.params.table}`,
    }
  },
  created() {},
  watch: {
    $route() {}
  },
  methods: {
    submitUpload() {
      this.$refs.upload.submit();
    },
  },
};
</script>
