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
    el-divider(
      direction="vertical")
    el-button(
      icon="el-icon-copy-document"
      size="mini"
      @click.stop="handleCopyAllData()") copy all
    el-button(
      icon="el-icon-brush"
      size="mini"
      @click.stop="triggerFilterList()") filter
    div(v-show="showFilterList")
      el-divider(
        content-position="left")
        el-checkbox(
          :indeterminate="isIndeterminate"
          v-model="isFilterColumnAll"
          @change="handleFilterColumnAll") check all
      el-checkbox-group(
        v-model="customColumn"
        @change="handleFilterColumn")
        el-checkbox(
          v-for="name in columnNames"
          :label="name"
          :key="name") {{name}}
    el-table(
      v-if="column"
      :data="rowData"
      empty-text="empty data"
      stripe
      :height="wh"
      style="width: 100%")
        el-table-column(
          :show-overflow-tooltip="isShowOverflowTooltip"
          v-for="columnData in column.getColumnData()"
          v-if="filterColumn(columnData.column_name)"
          :key="columnData.column_name"
          :width="(isSetＷidth()) ? columnData.text_rect.width + 6 : undefined"
          :prop="columnData.column_name"
          :label="`${columnData.column_name}`")
        el-table-column(
          fixed="left"
          label="Tools"
          width="150")
          el-button-group(slot-scope="scope")

            el-button(
              type="primary"
              @click="handleOpenEditDialog(scope.row)"
              icon="el-icon-edit"
              size="mini")

            el-button(
              type="danger"
              icon="el-icon-delete"
              size="mini"
              @click="handleDeleteConfirm(scope.row)")

            el-button(
              type="info"
              icon="el-icon-copy-document"
              size="mini"
              @click="handleCopyData(scope.row)")
    el-pagination(:page-size="50"
      @next-click="handleNextClick"
      @prev-click="handlePrevClick"
      @size-change="handleSizeChange"
      background
      :page-sizes="[50, 100, 200, 300, 400, 500]"
      :current-page.sync="currentPage"
      layout="sizes, prev, next")

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
    margin: 10px 0px;
  }
</style>
<script>
import includes from 'lodash/includes'
import forEach from 'lodash/forEach'
import cloneDeep from 'lodash/cloneDeep'
import copy from 'copy-to-clipboard'
import stringify from 'stringify-object'
import JSONbig from 'json-bigint'
import Cookies from 'js-cookie'
import Column from '../../utils/column'
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Rows',

  data() {
    return {
      wh: window.innerHeight - 270,
      tokenRowData: [],
      rowData: [],
      itemData: {},
      prevNext: 'next',
      column: null,
      componentWidth: 0,
      currentPage: 1,
      editDialogVisible: false,
      editInputData: [],
      isShowOverflowTooltip: true,
      loading: true,
      customColumn: [],
      isFilterColumnAll: true,
      isIndeterminate: false,
      showFilterList: false,
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
  computed: {
    columnNames() {
      return this.column ? this.column.getColumnData().map(data => data.column_name) : []
    }
  },
  methods: {
    triggerFilterList() {
      this.showFilterList = !this.showFilterList
    },
    handleFilterColumn(val) {
      this.customColumn = val || []
      this.isFilterColumnAll = this.customColumn.length === this.columnNames.length
      this.isIndeterminate = this.customColumn.length > 0 && this.customColumn.length < this.columnNames.length
    },
    handleFilterColumnAll(val) {
      this.customColumn = val ? this.columnNames : [];
      this.isIndeterminate = false;
    },
    filterColumn(name) {
      return includes(this.customColumn, name)
    },
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

      this.fetch(this.tokenRowData[this.currentPage])
    },

    async fetch(findItemData, forceRender = false) {
      this.loading = true

      try {
        const res = await service.request('rowtoken', {
          data: {
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
            item: findItemData,
            prevnext: this.prevNext,
            pagesize: parseInt(this.$route.params.pagesize, 10)
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
        } else {
          this.$message({
            type: 'error',
            showClose: true,
            message: 'No data'
          });

          if (forceRender) {
            this.rowData = []
          }
          this.loading = false
          return false
        }
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error
        });

        this.loading = false
        return false
      }

      this.loading = false
      return true
    },

    async fetchColumn() {
      const column = new Column(this.$route.params.keyspace, this.$route.params.table)
      await column.init()
      this.column = column
      this.customColumn = this.columnNames
    },

    async handlePrevClick(page) {
      this.prevNext = 'next'
      this.fetch(this.tokenRowData[page])
    },

    async handleNextClick(page) {
      this.tokenRowData[page] = this.rowData[this.rowData.length - 1]
      this.prevNext = 'next'

      if (!await this.fetch(this.tokenRowData[page])) {
        this.currentPage = page - 1
      }
    },

    handleSizeChange(pagesize) {
      this.currentPage = 1
      this.$router.push({
        name: 'rowstoken',
        params: {
          pagesize
        }
      })
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
          message
        });
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error.body.message
        });
      }

      this.fetch(this.tokenRowData[this.currentPage], true)
    },

    handleCopyData(rowData) {
      copy(stringify(rowData));

      this.$message({
        type: 'success',
        showClose: true,
        message: 'copy to clipboard'
      });
    },

    handleCopyAllData() {
      copy(stringify(this.rowData));

      this.$message({
        type: 'success',
        showClose: true,
        message: 'copy to clipboard'
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

    changeIsShowOverflowTooltip(bool) {
      Cookies.set('isShowOverflowTooltip', bool)
    },

    jsonParams(jsonString) {
      try {
        return JSONbig.parse(jsonString)
      } catch (e) {
        return jsonString
      }
    },
  }
};
</script>
