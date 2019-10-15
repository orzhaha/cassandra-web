<template lang="pug">
  div(
    class="w100"
    v-loading="loading")
    codemirror(
      v-model="content"
      :options="cmOptions")
</template>

<style>
  .w100 {
    width: 100%;
  }
  .CodeMirror {
    border: 1px solid #eee;
    height: auto;
  }
</style>

<script>
/* eslint import/no-extraneous-dependencies: ["error", {"peerDependencies": true}] */
import 'codemirror/mode/sql/sql'
/* eslint import/no-extraneous-dependencies: ["error", {"peerDependencies": true}] */
import 'codemirror/theme/monokai.css'
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Definition',
  data() {
    return {
      content: null,
      loading: false,
      cmOptions: {
        // codemirror options
        tabSize: 4,
        mode: 'text/x-cassandra',
        theme: 'monokai',
        lineNumbers: true,
        line: true,
        // more codemirror options, 更多 codemirror 的高级配置...
      }
    }
  },
  created() {
    document.title = `Definition-${this.$route.params.table}`
    this.fetch()
  },
  watch: {
    $route() {
      this.fetch()
    }
  },
  methods: {
    async fetch() {
      this.loading = true

      try {
        const res = await service.request('describe', {
          query: {
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
          }
        })

        this.content = res.text
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error
        });
      }

      this.loading = false
    }
  }
};
</script>
