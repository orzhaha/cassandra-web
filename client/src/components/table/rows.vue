<template lang="pug">
  div(class="w100")
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
          :label="`${key} (${getCqlType(key)})`")
            template(slot-scope="scope")
              img.iconKey(
                v-if="isPartitionKey(key)"
                src="../../assets/icon-key.svg")
              template(v-if="inputType(key) === '' || inputType(key) === 'textarea'")
                span(
                  v-if="!isEdit(scope.$index, key)")  {{scope.row[key]}}
                el-input(
                  v-else
                  :type="inputType(key)"
                  :autosize="{ minRows: 1, maxRows: 10}"
                  v-model="scope.row[key]")

              template(v-else)
                span(
                  v-if="!isEdit(scope.$index, key)") {{scope.row[key]}}
                codemirror(
                  v-else
                  v-model="scope.row[key]"
                  :options="cmOptions")
        el-table-column(
          fixed="right"
          label="Tools"
          width="110")
          el-button-group(slot-scope="scope")
            el-button(
              v-if="!isRowEditActive(scope.$index)"
              type="primary"
              @click="activeRowEdit(scope.$index)"
              icon="el-icon-edit"
              size="small")

            template(v-else)
              el-button-group
                el-button(
                  type="success"
                  @click="handleEditRow(scope.$index, scope.row)"
                  icon="el-icon-check"
                  size="small")
                el-button(
                  type="info"
                  @click="handleCancelRow(scope.$index)"
                  icon="el-icon-close"
                  size="small")

            el-button(
              v-if="!isRowEditActive(scope.$index)"
              type="danger"
              icon="el-icon-delete"
              size="small"
              @click="handleDelete")

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
  .CodeMirror {
    border: 1px solid #eee;
    height: auto;
  }

  .iconKey {
    width: 15px;
    height: auto;
  }
</style>
<script>
import api from '@/api'
import get from 'lodash/get'
import forEach from 'lodash/forEach'
import cloneDeep from 'lodash/cloneDeep'
import includes from 'lodash/includes'
import JSONbig from 'json-bigint'
import getType from '@/config/data-type'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/theme/monokai.css'

const service = api.make('root')

export default {
  name: 'Rows',

  data() {
    return {
      keys: [],
      rowdata: [],
      rowcount: 0,
      pagesize: 50,
      isRowEdit: null,
      originalData: [],
      types: {},
      partitionKey: [],
      cmOptions: {
        mode: {
          name: 'javascript',
          json: true
        },
        theme: 'monokai',
        line: true,
        lineWrapping: true,
        autofocus: true,
      },
    }
  },
  created() {
    this.fetch()
    this.fetchType()
  },
  watch: {
    $route() {
      this.fetch()
    }
  },
  methods: {
    isEdit(index, rowKey) {
      return this.isRowEditActive(index) && !this.isPartitionKey(rowKey)
    },

    isRowEditActive(index) {
      return this.isRowEdit === index
    },

    async activeRowEdit(index) {
      if (this.isRowEdit === null) {
        this.isRowEdit = index

        return
      }

      if (this.isRowEditActive(index)) {
        this.isRowEdit = null

        return
      }

      if (this.isDataChange(
        JSON.stringify(this.originalData[this.isRowEdit]),
        JSON.stringify(this.rowdata[this.isRowEdit])
      )) {
        await this.$confirm('Do you want to save data on current change ?', '', {
          confirmButtonText: 'Save',
          cancelButtonText: 'Cancel',
        }).then(() => {
          this.handleEdit(this.isRowEdit, this.rowdata[this.isRowEdit])
        }).catch(() => {
          this.fetch()
        })
      }

      this.isRowEdit = index
    },

    handleDelete() {
      this.$confirm('Are you sure ?', '', {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonClass: 'el-button--danger'
      }).then(() => {
        this.$message({
          showClose: true,
          duration: 0,
          message: 'ken 還沒開 api'
        })
      }).catch(() => {})
    },

    handleEditRow(index, rowData) {
      if (this.isRowEditActive(index)) {
        this.handleEdit(index, rowData)
      }

      this.activeRowEdit(index)
    },

    handleCancelRow(index) {
      this.fetch()

      this.activeRowEdit(index)
    },

    isPartitionKey(key) {
      return includes(this.partitionKey, key)
    },
    getCqlType(key) {
      return get(this.types, `${key}.cql`, '')
    },

    getJSType(key) {
      return get(this.types, `${key}.js`, 'string')
    },

    inputType(key) {
      if (this.getJSType(key) === 'number') {
        return ''
      }

      if (this.getJSType(key) === 'string') {
        return 'textarea'
      }

      return null
    },

    async fetch() {
      try {
        const res = await service.request('row', {
          query: {
            limit: 1000,
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
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
                item[itemKey] = JSONbig.stringify(itemData)
              } else {
                item[itemKey] = itemData
              }
            })
            return item
          })

          this.originalData = cloneDeep(this.rowdata)

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

    async fetchType() {
      const res = await service.request('columns', {
        query: {
          keyspace: this.$route.params.keyspace,
          table: this.$route.params.table,
        }
      })

      const data = res.get()
      const types = {}

      forEach(data, (item) => {
        if (get(item, 'column_name')) {
          types[item.column_name] = {
            cql: get(item, 'type', ''),
            js: getType(get(item, 'type', ''))
          }

          if (get(item, 'kind') === 'partition_key') {
            this.partitionKey.push(get(item, 'column_name'))
          }
        }
      })

      this.types = types
    },

    handleCurrentChange(page) {
      this.$router.push({
        name: 'rows',
        params: {
          page
        }
      })
    },
    handleSizeChange(pagesize) {
      this.$router.push({
        name: 'rows',
        params: {
          pagesize
        }
      })
    },
    rowFormatter(row, column, cellValue) {
      return cellValue
    },

    isDataChange(data1, data2) {
      return JSON.stringify(data1) !== JSON.stringify(data2)
    },

    async handleEdit(index, newRowData) {
      if (!this.isDataChange(
        JSON.stringify(this.originalData[index]),
        JSON.stringify(newRowData)
      )) {
        return
      }

      // TODO: edit partition_key
      // if (this.isPartitionKey(rowKey)) {
      //   await this.$confirm('In the case, it will change the partition_key. are you sure？', '', {
      //     confirmButtonText: 'Do it',
      //     cancelButtonText: 'Cancel',
      //   }).then(() => {
      //     // this.updateData(row)
      //     this.$message({
      //       showClose: true,
      //       duration: 0,
      //       message: '尚未開放修改 partition_key 功能'
      //     })
      //   }).catch(() => {})

      //   this.fetch()

      //   return
      // }

      this.updateData(newRowData)
    },

    async updateData(row) {
      const cRow = cloneDeep(row)

      forEach(cRow, (itemData, itemKey) => {
        cRow[itemKey] = this.jsonParams(itemData)
      })

      try {
        const res = await service.request('save', {
          data: {
            item: JSONbig.stringify(cRow),
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
          }
        })

        const message = (res.get() === []) ? 'success' : res.get()

        this.$message({
          type: 'success',
          showClose: true,
          duration: 0,
          message
        });

        this.fetch()
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
        return JSONbig.parse(jsonString)
      } catch (e) {
        return jsonString
      }
    }
  }
};
</script>
