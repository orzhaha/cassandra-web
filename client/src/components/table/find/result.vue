<template lang="pug">
  div(
    v-if="rowData.length !== 0"
    class="w100")
    el-table(
      :data="rowData"
      v-loading="loading"
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
          width="150")
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

            el-button(
              type="info"
              icon="el-icon-share"
              size="small"
              @click="handleCopyData(scope.row)")
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
</style>
<script>
import api from '@/api'
import forEach from 'lodash/forEach'
import cloneDeep from 'lodash/cloneDeep'
import copy from 'copy-to-clipboard'
import stringify from 'stringify-object'
import JSONbig from 'json-bigint'

const service = api.make('root')

export default {
  name: 'Result',
  props: [
    'rowData',
    'rowCount',
    'handleCurrentChange',
    'handleSizeChange',
    'page',
    'pagesize',
    'column',
    'find',
    'isShowOverflowTooltip',
    'componentWidth',
  ],
  data() {
    return {
      isRowEdit: null,
      editDialogVisible: false,
      editInputData: [],
      loading: false,
    }
  },
  created() {
  },
  watch: {
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
        row[column.column_name] = this.jsonParams(this.editInputData[index])
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

      this.loading = false

      this.find(false)
    },

    handleDeleteConfirm(rowData) {
      this.$confirm('Are you sure ?', '', {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        confirmButtonClass: 'el-button--danger'
      }).then(() => {
        JSON.stringify(rowData)
        this.deleteData(rowData)
      }).catch(() => {})
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
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error.body.message
        });
      }

      this.find(false)
    },

    handleCopyData(rowData) {
      copy(stringify(rowData));

      this.$message({
        type: 'success',
        showClose: true,
        duration: 0,
        message: `copy to clipboard`
      });
    },

    isSetＷidth() {
      if (this.column && this.componentWidth !== 0) {
        if ((this.column.getCloumnTextTotalWidth() + 110) >= this.componentWidth) {
          return true
        }
      }

      return false
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
