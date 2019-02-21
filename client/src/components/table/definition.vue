<template lang="pug">
  codemirror(
    v-model="content"
    :options="cmOptions")
</template>

<style>
.CodeMirror {
  border: 1px solid #eee;
  height: auto;
}
</style>

<script>
import 'codemirror/mode/sql/sql'
import 'codemirror/theme/monokai.css'
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Row',
  data() {
    return {
      content: null,
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
    this.fetch()
  },
  watch: {
    $route() {
      this.fetch()
    }
  },
  methods: {
    async fetch() {
      try {
        const res = await service.request('describe', {
          query: {
            kind: 'table',
            item: `${this.$route.params.keyspace}.${this.$route.params.table}`,
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
    }
  }
};
</script>
