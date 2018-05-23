<template lang="pug">
  el-table(
    @row-click='click'
    :data="keyspace"
    stripe
    style="width: 100%"
  )
    el-table-column(
      prop="keyspace_name"
      label="keyspace_name"
    )
    el-table-column(
      prop="replication.class"
      label="class"
    )
</template>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Home',
  data() {
    return {
      keyspace: []
    }
  },
  created() {
    this.fetch()
  },
  methods: {
    click: function(row, event, column) {
      this.$router.push('/table?keyspace='+row.keyspace_name)
    },
    async fetch() {
      try {
        const res = await service.request('allKeySpace')
        this.keyspace = res.get()
      } catch (error) {
        console.error('[error]', error)
      }
    }
  }
};


</script>