<template lang="pug">
  el-table(
    :data="tabledata"
    empty-text="empty data"
    stripe
    :row-style="rowStyle"
    style="width: 100%")
    el-table-column(
      prop="table_name"
      label="table_name")
      template(slot-scope="scope")
        span(
          @click="goToTable(scope.row)")
          i(
            v-bind:class="[scope.row.views ? 'el-icon-view' : 'el-icon-document']"
            class="table-list-icon")
          span {{scope.row.table_name}}
    el-table-column(
      prop="table_name"
      label="find")
      template(slot-scope="scope")
        span(
          @click="goToFind(scope.row)")
          i(class="table-list-icon el-icon-search")
          span find
    el-table-column(
      prop="table_name"
      label="columns")
      template(slot-scope="scope")
        span(
          @click="goToColumns(scope.row)")
          i(class="table-list-icon el-icon-info")
          span columns
    el-table-column(
      prop="table_name"
      label="definition")
      template(slot-scope="scope")
        span(
          @click="goToDefinition(scope.row)")
          i(class="table-list-icon el-icon-news")
          span definition
    el-table-column(
      prop="table_name"
      label="export")
      template(slot-scope="scope")
        span(
          @click="goToExport(scope.row)")
          i(class="table-list-icon el-icon-download")
          span export
    el-table-column(
      prop="table_name"
      label="import")
      template(slot-scope="scope")
        span(
          @click="goToImport(scope.row)")
          i(class="table-list-icon el-icon-upload2")
          span import
    el-table-column(
      prop="id"
      show-overflow-tooltip="true"
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
      tabledata: []
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
    goToTable(row) {
      this.$router.push({
        path: `${row.keyspace_name}/${row.table_name}`,
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

    async fetch() {
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
    },
    rowStyle() {
      return { cursor: 'pointer' }
    },
  }
};
</script>
