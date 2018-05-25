<template lang="pug">
  el-table(
    @row-click='click'
    :data="tabledata"
    stripe
    style="width: 100%"
  )
    el-table-column(
      prop="table_name"
      label="table_name"
    )
    el-table-column(
      prop="id"
      label="id"
    )
</template>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Table',
  data() {
    return {
      tabledata: []
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
    click(row) {
      this.$router.push({
        name: 'row',
        params: {
          table: `${row.keyspace_name}.${row.table_name}`
        }
      })
    },
    async fetch() {
      try {
        const res = await service.request('table', {
          query: {
            'keyspace': this.$route.params.keyspace
          }
        })
        this.tabledata = res.get()
      } catch (error) {
      }
    }
  }
};
</script>
