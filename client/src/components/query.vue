<template lang="pug">
  div(
    class="w100"
    v-loading="loading")
    el-table(
      :data="rowdata"
      :highlight-current-row="true"
      empty-text="empty data"
      stripe
      style="width: 100%")
        el-table-column(
        v-for="key in keys"
        :key="key"
        :formatter="rowFormatter"
        :label="key")
          template(slot-scope="scope")
            div
              span {{scope.row[key]}}
    el-pagination(
      background
      :total="rowcount"
      layout="total")
</template>

<style>
  .w100 {
    width: 100%;
  }
</style>
<script>
import api from '@/api'
import { forEach, isArray } from 'lodash'

const service = api.make('root')

export default {
  name: 'Query',
  data() {
    return {
      keys: [],
      rowdata: [],
      rowcount: 0,
      loading: false,
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
      this.loading = true

      try {
        const res = await service.request('query', {
          data: {
            query: this.$route.params.query
          }
        })

        const rows = res.get()[0]

        if (!isArray(rows)) {
          this.$message({
            type: 'success',
            showClose: true,
            duration: 0,
            message: rows
          });

          this.loading = false

          return;
        }

        if (rows !== undefined && rows.length > 0) {
          this.keys = Object.keys(rows[0])
          this.rowdata = rows.map((row) => {
            const item = row
            forEach(item, (itemData, itemKey) => {
              if (typeof (itemData) === 'object') {
                item[itemKey] = JSON.stringify(itemData)
              } else {
                item[itemKey] = itemData
              }
            })
            return item
          })
          this.rowcount = rows.length
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

      this.loading = false
    },
    rowFormatter(row, column, cellValue) {
      return cellValue
    }
  }
};
</script>
