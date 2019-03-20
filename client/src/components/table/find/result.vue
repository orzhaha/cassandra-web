<template lang="pug">
  div(
    v-if="rowData.length === 0")
  div(
    v-else
    class="w100")
    el-table(
      :data="rowData"
      :highlight-current-row="true"
      empty-text="empty data"
      stripe
      :row-style="rowStyle"
      style="width: 100%")
        el-table-column(
          show-overflow-tooltip="true"
          v-for="columnData in column.getColumnData()"
          :key="columnData['column_name']"
          :formatter="rowFormatter"
          :label="`${columnData['column_name']} (${column.getCqlType(columnData['column_name'])})`")
            template(slot-scope="scope")
              img.iconKey(
                v-if="column.isPartitionKey(columnData['column_name'])"
                src="../../../assets/icon-key.svg"
                title="PartitionKey")
              img.iconKey(
                v-if="column.isClusteringKey(columnData['column_name'])"
                src="../../../assets/key-ring.svg"
                title="Clustering")
              template(v-if="column.inputType(columnData['column_name']) === ''|| column.inputType(columnData['column_name']) === 'textarea'")
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
import forEach from 'lodash/forEach'
import cloneDeep from 'lodash/cloneDeep'
import JSONbig from 'json-bigint'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/theme/monokai.css'

const service = api.make('root')

export default {
  name: 'Result',
  props: [
    'rowData',
    'originalData',
    'column',
    'find',
  ],
  data() {
    return {
      isRowEdit: null,
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
  },
  watch: {
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
        JSON.stringify(this.rowData[this.isRowEdit])
      )) {
        await this.$confirm('Do you want to save data on current change ?', '', {
          confirmButtonText: 'Save',
          cancelButtonText: 'Cancel',
        }).then(() => {
          this.handleEdit(this.isRowEdit, this.rowData[this.isRowEdit])
        }).catch(() => {
          this.find(false)
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
        this.find(false)
      }).catch(() => {})
    },

    handleEditRow(index, rowData) {
      if (this.isRowEditActive(index)) {
        this.handleEdit(index, rowData)
      }

      this.activeRowEdit(index)
    },

    handleCancelRow(index) {
      this.find(false)

      this.activeRowEdit(index)
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

        this.find(false)
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error.body.message
        });
      }
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

        this.find(false)
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
