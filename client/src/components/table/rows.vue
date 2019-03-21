<template lang="pug">
  div(class="w100")
    el-checkbox(
      class="is-show-overflow-tooltip-checkbox"
      @change="changeIsShowOverflowTooltip"
      v-model="isShowOverflowTooltip") Show Overflow Tooltip
    el-table(
      v-if="column"
      :data="rowdata"
      :highlight-current-row="true"
      empty-text="empty data"
      stripe
      :row-style="rowStyle"
      style="width: 100%")
        el-table-column(
          :show-overflow-tooltip="isShowOverflowTooltip"
          v-for="columnData in column.getColumnData()"
          :key="columnData['column_name']"
          :formatter="rowFormatter"
          :label="`${columnData['column_name']}`")
            template(slot-scope="scope")
              img.iconKey(
                v-if="column.isPartitionKey(columnData['column_name'])"
                src="../../assets/icon-key.svg"
                title="PartitionKey")
              img.iconKey(
                v-if="column.isClusteringKey(columnData['column_name'])"
                src="../../assets/key-ring.svg"
                title="Clustering")
              template(v-if="column.inputType(columnData['column_name']) === '' || column.inputType(columnData['column_name']) === 'textarea'")
                span(
                  v-if="!isEdit(scope.$index, columnData['column_name'])")  {{scope.row[columnData['column_name']]}}
                el-input(
                  v-else
                  :type="column.inputType(columnData['column_name'])"
                  :autosize="{ minRows: 1, maxRows: 10}"
                  v-model="scope.row[columnData['column_name']]")

              template(v-else)
                span(
                  v-if="!isEdit(scope.$index, columnData['column_name'])") {{scope.row[columnData['column_name']]}}
                codemirror(
                  v-else
                  v-model="scope.row[columnData['column_name']]"
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
              @click="handleDelete(scope.row)")

    el-pagination(:page-size="20"
      @current-change="handleCurrentChange"
      @prev-click="handleCurrentChange"
      @next-click="handleCurrentChange"
      @size-change="handleSizeChange"
      background
      :pageSize="pagesize"
      :page-sizes="[50, 100, 200, 300, 400, 500]"
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
  .is-show-overflow-tooltip-checkbox {
    margin: 10px 20px;
    padding: 12px 0px;
  }
</style>
<script>
import api from '@/api'
import forEach from 'lodash/forEach'
import cloneDeep from 'lodash/cloneDeep'
import JSONbig from 'json-bigint'
import Cookies from 'js-cookie'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/theme/monokai.css'
import Column from '../../utils/column'

const service = api.make('root')

export default {
  name: 'Rows',

  data() {
    return {
      rowdata: [],
      rowcount: 0,
      column: null,
      pagesize: 50,
      isRowEdit: null,
      originalData: [],
      isShowOverflowTooltip: true,
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

    const isNotCollapse = Cookies.get('isNotCollapse')

    if (isNotCollapse !== undefined) {
      this.isNotCollapse = isNotCollapse === 'true'
    }

    const isShowOverflowTooltip = Cookies.get('isShowOverflowTooltip')

    if (isShowOverflowTooltip !== undefined) {
      this.isShowOverflowTooltip = isShowOverflowTooltip === 'true'
    }
  },
  watch: {
    $route() {
      this.fetch()
    }
  },
  methods: {
    isEdit(index, rowKey) {
      return this.isRowEditActive(index) && (!this.column.isPartitionKey(rowKey) && !this.column.isClusteringKey(rowKey))
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

    handleDelete(rowData) {
      this.$confirm('Are you sure ?', '', {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonClass: 'el-button--danger'
      }).then(() => {
        JSON.stringify(rowData)
        this.deleteData(rowData)
        this.fetch()
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
      const column = new Column(this.$route.params.keyspace, this.$route.params.table)
      await column.init()
      this.column = column
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
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error.body.message
        });
      }

      this.fetch()
    },

    async deleteData(row) {
      const cRow = cloneDeep(row)

      forEach(cRow, (itemData, itemKey) => {
        cRow[itemKey] = this.jsonParams(itemData)
      })

      try {
        const res = await service.request('delete', {
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
    },
    changeIsShowOverflowTooltip(bool) {
      Cookies.set('isShowOverflowTooltip', bool)
    },
  }
};
</script>
