<template lang="pug">
  div.w100
    el-table(
      :data="rowdata"
      :highlight-current-row="true"
      empty-text="empty data"
      stripe
      :row-style="rowStyle"
      style="width: 100%")
        el-table-column(
          v-for="key in keys"
          :key="key"
          :formatter="rowFormatter"
          :label="key")
          template(slot-scope="scope")
            div(
              @dblclick.stop="handleActive(`${scope.$index}${key}`, scope.row[key])")
              span(v-if="!isEdit[`${scope.$index}${key}`]") {{scope.row[key]}}
              el-input(v-else
                @blur.stop="handleEdit(`${scope.$index}${key}`, scope.row[key], scope.row)"
                v-model="scope.row[key]")
</template>

<style>
  .w100 {
    width: 100%;
  }
</style>
<script>
import api from '@/api'
import { forEach } from 'lodash'
import JSONbig from 'json-bigint'

const service = api.make('root')

export default {
  name: 'Columns',
  data() {
    return {
      keys: [],
      rowdata: [],
      isEdit: {},
      originalData: {}
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
        const res = await service.request('columns', {
          query: {
            keyspace: this.$route.params.keyspace,
            table: this.$route.params.table,
          }
        })

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
    },
    rowFormatter(row, column, cellValue) {
      return cellValue
    },
    handleActive() {
      return null
    },
    handleEdit() {
      return null
    },
    // handleActive(key, colData) {
    //   this.originalData[key] = colData
    //   this.$set(this.isEdit, key, !this.isEdit[key])
    // },
    // async handleEdit(key, colData, row) {
    //   this.$set(this.isEdit, key, !this.isEdit[key])

    //   if (this.originalData[key] === colData) {
    //     return
    //   }
    //   const cRow = cloneDeep(row)

    //   forEach(cRow, (itemData, itemKey) => {
    //     cRow[itemKey] = this.jsonParams(itemData)
    //   })

    //   try {
    //     const res = await service.request('save', {
    //       data: {
    //         item: JSONbig.stringify(cRow),
    //         table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
    //       }
    //     })

    //     const message = (res.get() === []) ? 'success' : res.get()

    //     this.$message({
    //       // type: 'success',
    //       showClose: true,
    //       duration: 0,
    //       message
    //     });
    //   } catch (error) {
    //     this.$message({
    //       type: 'error',
    //       showClose: true,
    //       duration: 0,
    //       message: error.body.message
    //     });
    //   }
    // },
    rowStyle() {
      return { cursor: 'pointer' }
    },
    jsonParams(jsonString) {
      try {
        return JSONbig.parse(jsonString)
      } catch (e) {
        return jsonString
      }
    }
  }
};
</script>
