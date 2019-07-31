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
          :label="key")
          template(slot-scope="scope")
              span {{scope.row[key]}}
</template>

<style>
  .w100 {
    width: 100%;
  }
</style>
<script>
import { forEach } from 'lodash'
import JSONbig from 'json-bigint'
import api from '@/api'

const service = api.make('root')

export default {
  name: 'HostInfo',
  data() {
    return {
      keys: [],
      rowdata: [],
      isEdit: {},
      loading: false,
    }
  },
  created() {
    document.title = 'HostInfo'
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
        const res = await service.request('hostinfo', {})
        const rows = res.get()

        if (rows !== undefined && rows.length > 0) {
          this.keys = Object.keys(rows[0])
          this.rowdata = rows.map((row) => {
            const item = row
            forEach(item, (itemData, itemKey) => {
              if (typeof (itemData) === 'object') {
                item[itemKey] = JSONbig.stringify(itemData)
              } else {
                item[itemKey] = itemData
              }
            })
            return item
          })
        } else {
          this.keys = []
          this.rowdata = []
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
  }
};
</script>
