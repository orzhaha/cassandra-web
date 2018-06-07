<template lang="pug">
  div(class="w100")
    el-table(
      :data="rowdata"
      :highlight-current-row="true"
      empty-text="empty data"
      stripe
      style="width: 100%")
        el-table-column(
        v-for="key in keys"
        :formatter="rowFormatter"
        :prop="key"
        :label="key"
        )
    el-pagination( :page-size="20"
      @current-change="handleCurrentChange"
      @prev-click="handleCurrentChange"
      @next-click="handleCurrentChange"
      @size-change="handleSizeChange"
      background
      :pageSize="pagesize"
      :page-sizes="[50, 100, 200, 300, 400]"
      :total="rowcount"
      layout="total, sizes, prev, pager, next")
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
      rowdata: [],
      rowcount: 0,
      pagesize: 50
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
            table: this.$route.params.table,
            page: this.$route.params.page,
            pagesize: this.$route.params.pagesize
          }
        })
        const rows = res.get('row')

        if (rows !== undefined && rows.length > 0) {
          this.keys = Object.keys(rows[0])
          this.rowdata = rows
          this.rowcount = res.get('count')
        } else {
          this.keys = []
          this.rowdata = []
          this.rowcount = 0
        }
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error
        });
      }
    },
    handleCurrentChange(page) {
      this.$router.push({
        name: 'row',
        params: {
          page
        }
      })
    },
    handleSizeChange(pagesize) {
      this.$router.push({
        name: 'row',
        params: {
          pagesize
        }
      })
    },
    rowFormatter(row, column, cellValue) {
      if (typeof (cellValue) === 'object') {
        return JSON.stringify(cellValue)
      }
      return cellValue
    }
  }
};
</script>
