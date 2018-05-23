<template lang="pug">
  el-table(
    @row-click='click'
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
  name: 'Home',
  data() {
    return {
      keys: [],
      rowdata: []
    }
  },
  created() {
    this.fetch()
  },
  methods: {
    click: function(row, event, column) {
        // this.$router.push('/row?table='+row.table_name)
    },
    async fetch() {
      try {
        const res = await service.request('allrowbytable', {
            query: {
                "limit": 200,
                "table": this.$route.query.table_name
            }
        })

        const rows = res.get('row')

        this.keys = Object.keys(rows[0])
        this.rowdata = rows
      } catch (error) {
        console.error('[error]', error)
      }
    }
  }
};


</script>