<template lang="pug">
  div(
    class="w100"
    id="tablerows"
    ref="tablerows"
    v-loading="loading")
    div(class="w100"
        @keyup.enter="find()")
      el-button(
        type="text"
        class="caret-bottom"
        v-bind:class="[isCollapse ? 'el-icon-caret-top' : 'el-icon-caret-bottom']"
        @click.stop="isCollapse = !isCollapse") Folding Find Column
      el-divider(
        direction="vertical")
      el-checkbox(
        @change="changeIsNotCollapse"
        v-model="isNotCollapse") Do Not Collapse
      el-divider(
        direction="vertical")
      el-checkbox(
        @change="changeIsShowOverflowTooltip"
        v-model="isShowOverflowTooltip") Show Overflow Tooltip
      el-divider(
        direction="vertical")
      el-checkbox(
        @change="changeIsAllowFilter"
        v-model="isAllowFilter") Allow Filter
      el-divider(
        direction="vertical")
      el-button(
        icon="el-icon-copy-document"
        size="small"
        @click.stop="handleCopyAllData()") copy all
      el-button(
        icon="el-icon-brush"
        size="small"
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
      template(
        v-if="isNotCollapse || isCollapse")
        el-table(
          v-if="column"
          :data="allowFilterColumn"
          :highlight-current-row="true"
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
            label="Operator")
            template(slot-scope="scope")
              template(v-if="column.isPartitionKey(scope.row.column_name)")
                el-select(
                  v-model="columnInput[scope.$index].operator"
                  placeholder='operator')
                  template(v-if="isAllowFilter")
                    el-option(
                      v-for="item in clusteringOperator"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value")
                  template(v-else="isAllowFilter")
                    el-option(
                      v-for="item in partitionOperator"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value")
              template(v-else-if="isAllowFilter || column.isClusteringKey(scope.row.column_name)")
                el-select(
                  v-model="columnInput[scope.$index].operator"
                  placeholder='operator')
                  el-option(
                    v-for="item in clusteringOperator"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value")
          el-table-column(
            prop="kind"
            label="Value")
            template(slot-scope="scope")
              template(v-if="isAllowFilter || column.isPartitionKey(scope.row.column_name) || column.isClusteringKey(scope.row.column_name)")
                el-input(
                  v-model="columnInput[scope.$index].value"
                  :type="column.inputType(scope.row.column_name)"
                  :autosize="{ minRows: 1, maxRows: 10}")

              template(v-else)
        br
        el-button(
          type="danger"
          icon="el-icon-search"
          @click.stop="find()")
    div(class="w100")
      Result(
        :find="find"
        :rowData="rowData"
        :rowCount="rowCount"
        :handleCurrentChange="handleCurrentChange"
        :handleSizeChange="handleSizeChange"
        :page="page"
        :pagesize="pagesize"
        :column="column"
        :isShowOverflowTooltip="isShowOverflowTooltip"
        :componentWidth="componentWidth"
        :customColumn="customColumn")
</template>

<style>
  .w100 {
    width: 100%;
  }
  .caret-bottom {
    margin: 10px 0px;
  }
</style>

<script>
import copy from 'copy-to-clipboard'
import stringify from 'stringify-object'
import JSONbig from 'json-bigint'
import forEach from 'lodash/forEach'
import Cookies from 'js-cookie'
import api from '@/api'
import Result from './result'
import Column from '../../../utils/column'

const service = api.make('root')

export default {
  name: 'Find',
  components: {
    Result,
  },
  data() {
    return {
      rowData: [],
      rowCount: 0,
      columnInput: [],
      column: null,
      componentWidth: 0,
      page: 0,
      pagesize: 50,
      isCollapse: true,
      isNotCollapse: false,
      isShowOverflowTooltip: true,
      isAllowFilter: false,
      loading: false,
      partitionOperator: [
        {
          value: '=',
          label: '=',
        },
      ],
      clusteringOperator: [
        {
          value: '=',
          label: '=',
        },
        {
          value: '>',
          label: '>',
        },
        {
          value: '>=',
          label: '>=',
        },
        {
          value: '<',
          label: '<',
        },
        {
          value: '<=',
          label: '<=',
        },
        {
          value: 'in',
          label: 'in',
        },
      ],
      customColumn: [],
      isFilterColumnAll: true,
      isIndeterminate: false,
      showFilterList: false,
    }
  },
  created() {
    document.title = `Find-${this.$route.params.table}`
    this.fetch()

    const isNotCollapse = Cookies.get('isNotCollapse')

    if (isNotCollapse !== undefined) {
      this.isNotCollapse = isNotCollapse === 'true'
    }

    const isShowOverflowTooltip = Cookies.get('isShowOverflowTooltip')

    if (isShowOverflowTooltip !== undefined) {
      this.isShowOverflowTooltip = isShowOverflowTooltip === 'true'
    }

    const isAllowFilter = Cookies.get('isAllowFilter')

    if (isAllowFilter !== undefined) {
      this.isAllowFilter = isAllowFilter === 'true'
    }

    this.$nextTick(() => {
      this.componentWidth = this.$refs.tablerows.clientWidth
    });
  },
  watch: {
    $route() {
      this.fetch()
    },
  },
  computed: {
    columnNames() {
      return this.column ? this.column.getColumnData().map(data => data.column_name) : []
    },
    allowFilterColumn() {
      return this.column
        ? this.column.getColumnData()
          .filter(data => this.isAllowFilter || (data.kind === 'partition_key' || data.kind === 'clustering'))
        : []
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
    async fetch() {
      this.loading = true

      const column = new Column(this.$route.params.keyspace, this.$route.params.table)
      await column.init()

      forEach(column.getColumnData(), (item, index) => {
        this.$set(this.columnInput, index, {})
        this.$set(this.columnInput[index], 'operator', '=')
        this.$set(this.columnInput[index], 'column_name', item.column_name)
        this.$set(this.columnInput[index], 'value', '')
      })

      this.column = column
      this.customColumn = this.columnNames
      this.loading = false
    },

    async find(isShowMsg = true) {
      this.isCollapse = false

      const reqItem = {}

      forEach(this.columnInput, (item) => {
        if (item.value !== '') {
          reqItem[item.column_name] = {
            operator: item.operator,
            value: item.value,
          }
        }
      })

      this.loading = true

      try {
        const res = await service.request('find', {
          data: {
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
            item: reqItem,
            page: parseInt(this.page, 10),
            pagesize: parseInt(this.pagesize, 10),
            isallowfilter: this.isAllowFilter,
          }
        })

        const rows = res.get('row')

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

        if (isShowMsg) {
          const message = (res.get() === []) ? 'success' : `conut: ${this.rowCount}`

          this.$message({
            type: 'success',
            showClose: true,
            duration: 0,
            message
          });
        }
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error.body.message
        });
      }

      this.loading = false
    },

    handleCurrentChange(page) {
      this.page = page
      this.find(false)
    },

    handleSizeChange(pagesize) {
      this.pagesize = pagesize
      this.find(false)
    },

    handleCopyAllData() {
      copy(stringify(this.rowData));

      this.$message({
        type: 'success',
        showClose: true,
        duration: 0,
        message: 'copy to clipboard'
      });
    },

    changeIsNotCollapse(bool) {
      Cookies.set('isNotCollapse', bool)
    },

    changeIsShowOverflowTooltip(bool) {
      Cookies.set('isShowOverflowTooltip', bool)
    },

    changeIsAllowFilter(bool) {
      Cookies.set('isAllowFilter', bool)
    },
  },
};
</script>
