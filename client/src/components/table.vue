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
  name: 'Home',
  data() {
    return {
      tabledata: []
    }
  },
  created() {
    this.fetch()
  },
  methods: {
    click: function(row, event, column) {
        this.$router.push('/row?table_name='+row.keyspace_name +"."+ row.table_name)
    },
    async fetch() {
      try {
        console.log(this.$route.query.keyspace)
        const res = await service.request('alltablebykeyspace', {
            query: {
                "keyspace": this.$route.query.keyspace
            }
        })
        this.tabledata = res.get()
        console.log(res.get())
      } catch (error) {
        console.error('[error]', error)
      }
    }
  }
};


</script>