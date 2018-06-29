<template lang="pug">
  div(class="w100")
    el-table(
      :data="rowdata"
      :highlight-current-row="true"
      empty-text="empty data"
      stripe
      :row-style="rowStyle")
        el-table-column(
        v-for="key in keys"
        :key="key"
        :formatter="rowFormatter"
        :label="key")
          template(slot-scope="scope")
            div(@dblclick.stop="handleActive(`${scope.$index}${key}`, scope.row[key])")
              span(v-if="!isEdit[`${scope.$index}${key}`]") {{scope.row[key]}}
              el-input(v-else
              @blur.stop="handleEdit(`${scope.$index}${key}`, scope.row[key], scope.row)"
              v-model="scope.row[key]"
              )
    el-pagination(:page-size="20"
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
import { forEach, cloneDeep } from 'lodash'

const service = api.make('root')

export default {
  name: 'Row',
  data() {
    return {
      keys: [],
      rowdata: [],
      rowcount: 0,
      pagesize: 50,
      isEdit: {},
      originalData: {},
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
      return cellValue
    },
    handleActive(key, colData) {
      this.originalData[key] = colData
      this.$set(this.isEdit, key, !this.isEdit[key])
    },
    async handleEdit(key, colData, row) {
      this.$set(this.isEdit, key, !this.isEdit[key])

      if (this.originalData[key] === colData) {
        return
      }
      const cRow = cloneDeep(row)

      forEach(cRow, (itemData, itemKey) => {
        cRow[itemKey] = this.jsonParams(itemData)
      })

      try {
        const res = await service.request('save', {
          data: {
            item: JSON.stringify(cRow),
            table: this.$route.params.table,
          }
        })

        const message = (res.get() === []) ? 'success' : res.get()

        this.$message({
          type: 'success',
          showClose: true,
          duration: 0,
          message
        });
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error.body.message
        });
      }
    },
    rowStyle() {
      return { cursor: 'pointer' }
    },
    jsonParams(jsonString) {
      try {
        return JSON.parse(jsonString)
      } catch (e) {
        return jsonString
      }
    }
  }
};
</script>
