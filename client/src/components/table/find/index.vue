<template lang="pug">
  div(
    class="w100"
    id="tablerows"
    ref="tablerows")
    div(class="w100")
      el-button(
        type="text"
        class="caret-bottom"
        v-bind:class="[isCollapse ? 'el-icon-caret-top' : 'el-icon-caret-bottom']"
        @click.stop="isCollapse = !isCollapse") Folding Find Column
      el-checkbox(
        @change="changeIsNotCollapse"
        v-model="isNotCollapse") Do Not Collapse
      el-checkbox(
        @change="changeIsShowOverflowTooltip"
        v-model="isShowOverflowTooltip") Show Overflow Tooltip
      template(
        v-if="isNotCollapse || isCollapse")
        el-table(
          v-if="column"
          :data="column.getColumnData()"
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
                  el-option(
                    v-for="item in partitionOperator"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value")
              template(v-else-if="column.isClusteringKey(scope.row.column_name)")
                el-select(
                  v-model="columnInput[scope.$index].operator"
                  placeholder='operator')
                  el-option(
                    v-for="item in clusteringOperator"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value")

              template(v-else)
          el-table-column(
            prop="kind"
            label="Value")
            template(slot-scope="scope")
              template(v-if="column.isPartitionKey(scope.row.column_name) || column.isClusteringKey(scope.row.column_name)")
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
        :column="column"
        :isShowOverflowTooltip="isShowOverflowTooltip"
        :componentWidth="componentWidth")
</template>

<style>
  .w100 {
    width: 100%;
  }
  .caret-bottom {
    margin: 10px 20px;
  }
</style>

<script>
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
      columnInput: [],
      column: null,
      componentWidth: 0,
      isCollapse: true,
      isNotCollapse: false,
      isShowOverflowTooltip: true,
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
    }
  },
  created() {
    this.fetch()

    const isNotCollapse = Cookies.get('isNotCollapse')

    if (isNotCollapse !== undefined) {
      this.isNotCollapse = isNotCollapse === 'true'
    }

    const isShowOverflowTooltip = Cookies.get('isShowOverflowTooltip')

    if (isShowOverflowTooltip !== undefined) {
      this.isShowOverflowTooltip = isShowOverflowTooltip === 'true'
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
  methods: {
    async fetch() {
      const column = new Column(this.$route.params.keyspace, this.$route.params.table)
      await column.init()

      forEach(column.getColumnData(), (item, index) => {
        this.$set(this.columnInput, index, {})
        this.$set(this.columnInput[index], 'operator', '=')
        this.$set(this.columnInput[index], 'column_name', item.column_name)
        this.$set(this.columnInput[index], 'value', '')
      })

      this.column = column
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

      try {
        const res = await service.request('find', {
          data: {
            table: `${this.$route.params.keyspace}.${this.$route.params.table}`,
            item: reqItem,
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

        if (isShowMsg) {
          const message = (res.get() === []) ? 'success' : `conut: ${rows.length}`

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
    },

    changeIsNotCollapse(bool) {
      Cookies.set('isNotCollapse', bool)
    },

    changeIsShowOverflowTooltip(bool) {
      Cookies.set('isShowOverflowTooltip', bool)
    },
  },
};
</script>
