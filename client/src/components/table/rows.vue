<template lang="pug">
  div(
    class="w100"
    id="tablerows"
    ref="tablerows"
    v-loading="loading")
    el-checkbox(
      class="is-show-overflow-tooltip-checkbox"
      @change="changeIsShowOverflowTooltip"
      v-model="isShowOverflowTooltip") Show Overflow Tooltip
    el-table(
      v-if="column"
      :data="rowData"
      empty-text="empty data"
      stripe
      style="width: 100%")
        el-table-column(
          :show-overflow-tooltip="isShowOverflowTooltip"
          v-for="columnData in column.getColumnData()"
          :key="columnData['column_name']"
          :width="(isSetＷidth()) ? columnData['text_rect']['width'] + 6 : undefined"
          :prop="columnData['column_name']"
          :label="`${columnData['column_name']}`")
        el-table-column(
          fixed="left"
          label="Tools"
          width="110")
          el-button-group(slot-scope="scope")
            el-button(
              type="primary"
              @click="handleOpenEditDialog(scope.row)"
              icon="el-icon-edit"
              size="small")

            el-button(
              type="danger"
              icon="el-icon-delete"
              size="small"
              @click="handleDeleteConfirm(scope.row)")

    el-pagination(:page-size="20"
      @current-change="handleCurrentChange"
      @prev-click="handleCurrentChange"
      @next-click="handleCurrentChange"
      @size-change="handleSizeChange"
      background
      :pageSize="pagesize"
      :page-sizes="[50, 100, 200, 300, 400, 500]"
      :total="rowCount"
      layout="total, sizes, prev, pager, next")

    el-dialog(
      :title="`${$route.params.keyspace}.${$route.params.table}`"
      :visible.sync="editDialogVisible")
      el-table(
        id="editTable"
        ref="editTable"
        v-if="column"
        :data="column.getColumnData()"
        empty-text="empty data"
        stripe)
        el-table-column(
          prop="column_name"
          label="Field")
        el-table-column(
          prop="kind"
          label="Kind")
        el-table-column(
          prop="type"
          label="Type")
        el-table-column(
          prop="kind"
          label="Value")
          template(slot-scope="scope")
            template(v-if="column.isPartitionKey(scope.row.column_name) || column.isClusteringKey(scope.row.column_name)")
              el-input(
                v-model="editInputData[scope.$index]"
                :type="column.inputType(scope.row.column_name)"
                :disabled="true")
            template(v-else)
              el-input(
                v-model="editInputData[scope.$index]"
                :type="column.inputType(scope.row.column_name)"
                :autosize="{ minRows: 1, maxRows: 10}")
      br
      el-button(
        type="danger"
        icon="el-icon-check"
        @click="handleSubmitEditDialog") Submit
      el-button(
        type="primary"
        @click="handleCloseEditDialog"
        icon="el-icon-close") Cancel

</template>

<style>
  .w100 {
    width: 100%;
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
import forEach from 'lodash/forEach'
import cloneDeep from 'lodash/cloneDeep'
import JSONbig from 'json-bigint'
import Cookies from 'js-cookie'
import Column from '../../utils/column'
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Rows',

  data() {
    return {
      rowData: [],
      rowCount: 0,
      column: null,
      componentWidth: 0,
      pagesize: 50,
      editDialogVisible: false,
      editInputData: [],
      isShowOverflowTooltip: true,
      loading: true,
    }
  },
  created() {
    document.title = `Rows-${this.$route.params.table}`
    this.$nextTick(() => {
      this.componentWidth = this.$refs.tablerows.clientWidth
    });

    this.fetch()
    this.fetchColumn()

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
      document.title = `Rows-${this.$route.params.table}`
      this.fetch()
    }
  },
  methods: {
    handleOpenEditDialog(row) {
      this.editDialogVisible = true

      forEach(this.column.getColumnData(), (column, index) => {
        if (this.column.getJSType(column.column_name) === 'boolean') {
          this.editInputData[index] = (row[column.column_name]) ? 'true' : 'false'
        } else {
          this.editInputData[index] = row[column.column_name];
        }

        this.$nextTick(() => {
          this.$refs.editTable.toggleRowExpansion(index, false);
          this.$refs.editTable.toggleRowExpansion(index, true);
        })
      })
    },

    handleCloseEditDialog() {
      this.editDialogVisible = false
    },

    async handleSubmitEditDialog() {
      const row = {}

      forEach(this.column.getColumnData(), (column, index) => {
        row[column.column_name] = this.jsonParams(this.editInputData[index], this.column.getJSType(column.column_name))
      })

      this.editDialogVisible = false
      this.loading = true

      try {
        const res = await service.request('save', {
          data: {
            item: JSONbig.stringify(row),
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
          }
        })

        const message = (res.get() === []) ? 'success' : res.get()

        this.$message({
          type: 'success',
          showClose: true,
          message
        });
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error.body.message
        });
      }

      this.loading = false

      this.fetch()
    },

    async fetch() {
      this.loading = true

      try {
        const res = await service.request('row', {
          query: {
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
            page: this.$route.params.page,
            pagesize: this.$route.params.pagesize
          }
        })
        const rows = res.get('row')

        if (rows !== undefined && rows.length > 0) {
          this.rowData = rows.map((row) => {
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

          this.rowCount = res.get('count')
        } else {
          this.rowData = []
          this.rowCount = 0
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

    async fetchColumn() {
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
      }).catch(err => err)
    },

    handleSizeChange(pagesize) {
      this.$router.push({
        name: 'rows',
        params: {
          pagesize
        }
      }).catch(err => err)
    },

    handleDeleteConfirm(rowData) {
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

    async deleteData(row) {
      const cRow = cloneDeep(row)

      forEach(cRow, (itemData, itemKey) => {
        cRow[itemKey] = this.jsonParams(itemData, this.column.getJSType(itemKey))
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
          message
        });
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error.body.message
        });
      }

      this.fetch()
    },

    isSetＷidth() {
      if (this.column && this.componentWidth !== 0) {
        if ((this.column.getCloumnTextTotalWidth() + 110) >= this.componentWidth) {
          return true
        }
      }

      return false
    },

    changeIsShowOverflowTooltip(bool) {
      Cookies.set('isShowOverflowTooltip', bool)
    },

    jsonParams(jsonString, columnType) {
      try {
        if (columnType === 'string') {
          return jsonString
        }

        return JSONbig.parse(jsonString)
      } catch (e) {
        return jsonString
      }
    },
  }
};
</script>
