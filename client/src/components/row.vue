<template lang="pug">
  div(class="w100")
    el-table(
      :data="rowdata"
      empty-text="empty data"
      stripe
      style="width: 100%")
        el-table-column(
        v-for="key in keys"
        :formatter="rowFormatter"
        :prop="key"
        :label="key"
        )
    //- el-pagination( :page-size="20"
    //-   @current-change="handleCurrentChange"
    //-   @prev-click="handleCurrentChange"
    //-   @next-click="handleCurrentChange"
    //-   :pager-count="11"
    //-   layout="prev, pager, next"
    //-   :total="1000")
</template>

<style>
  .w100 {
    width: 100%;
  }
</style>
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
            limit: 1000,
            table: this.$route.params.table
          }
        })
        const rows = res.get('row')

        if (rows !== undefined && rows.length > 0) {
          this.keys = Object.keys(rows[0])
          this.rowdata = rows
        } else {
          this.keys = []
          this.rowdata = []
        }
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error
        });
      }
    },
    handleCurrentChange(row) {
      this.$message({
        type: 'error',
        showClose: true,
        duration: 0,
        message: row
      });
    },
    rowFormatter(row, column, cellValue) {
      if (typeof(cellValue) === 'object') {
        return JSON.stringify(cellValue)
      }
      return cellValue
    }
  }
};
</script>
