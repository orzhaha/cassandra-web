<template lang="pug">
  div(
    class="w100"
    v-loading="loading")
    el-table(
      @cell-click="handlerCellClick"
      :data="tabledata"
      empty-text="empty data"
      stripe
      :row-style="rowStyle"
      style="width: 100%")
      el-table-column(
        column-key="table_name"
        label="table_name")
        template(slot-scope="scope")
          i(
            v-bind:class="[scope.row.views ? 'el-icon-view' : 'el-icon-document']"
            class="table-list-icon")
          span {{scope.row.table_name}}
      el-table-column(
        column-key="find"
        label="find")
        i(class="table-list-icon el-icon-search")
        span find
      el-table-column(
        column-key="columns"
        label="columns")
        i(class="table-list-icon el-icon-data-analysis")
        span columns
      el-table-column(
        column-key="definition"
        label="definition")
        i(class="table-list-icon el-icon-news")
        span definition
      el-table-column(
        column-key="export"
        label="export")
        i(class="table-list-icon el-icon-download")
        span export
      el-table-column(
        column-key="import"
        label="import")
        i(class="table-list-icon el-icon-upload2")
        span import
      el-table-column(
        column-key="truncate"
        label="truncate")
        i(class="table-list-icon el-icon-folder-delete")
        span truncate
      el-table-column(
        column-key="id"
        prop="id"
        :show-overflow-tooltip="true"
        label="id")
</template>
<style>
  .table-list-icon {
    margin-right: 4px;
    font-size: 1.4em;
  }

</style>
<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Table',
  data() {
    return {
      tabledata: [],
      loading: false,
    }
  },
  created() {
    document.title = `Keyspace-${this.$route.params.keyspace}`
    this.fetch()
  },
  watch: {
    $route() {
      document.title = `Keyspace-${this.$route.params.keyspace}`
      this.fetch()
    }
  },
  methods: {
    handlerCellClick(row, column) {
      switch (column.columnKey) {
        case 'table_name': {
          this.goToTable(row)
          break
        }
        case 'find': {
          this.goToFind(row)
          break
        }
        case 'columns': {
          this.goToColumns(row)
          break
        }
        case 'definition': {
          this.goToDefinition(row)
          break
        }
        case 'export': {
          this.goToExport(row)
          break
        }
        case 'import': {
          this.goToImport(row)
          break
        }
        case 'truncate': {
          this.goToTruncate(row)
          break
        }
        default: {
          break
        }
      }
    },
    goToTable(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/rowstoken/50`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    goToFind(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/find`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    goToColumns(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/columns`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    goToDefinition(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/definition`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    goToExport(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/export`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    goToImport(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/import`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    goToTruncate(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}/truncate`,
        params: {
          keyspace: row.keyspace_name,
          table: row.table_name
        }
      })
    },

    async fetch() {
      this.loading = true

      try {
        const res = await service.request('table', {
          query: {
            keyspace: this.$route.params.keyspace
          }
        })
        this.tabledata = res.get()
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error
        });
      }

      this.loading = false
    },
    rowStyle() {
      return { cursor: 'pointer' }
    },
  }
};
</script>
