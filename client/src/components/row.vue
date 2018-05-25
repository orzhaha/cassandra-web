<template lang="pug">
  el-table(
    :data="rowdata"
    stripe
    style="width: 100%")
      el-table-column(
      v-for="key in keys"
      :prop="key"
      :label="key"
      )
</template>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Row',
  data() {
    return {
      keys: [],
      rowdata: []
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
        const res = await service.request('row', {
          query: {
            'limit': 200,
            'table': this.$route.params.table
          }
        })

        const rows = res.get('row')

        this.keys = Object.keys(rows[0])
        this.rowdata = rows
      } catch (error) {
      }
    }
  }
};
</script>
