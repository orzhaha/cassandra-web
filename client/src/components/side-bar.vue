<template lang="pug">
el-menu(
  :default-active="$route.params.keyspace"
  class="el-menu-vertical-demo"
  v-loading="loading"
  :collapse-transition="false"
  @open="openSubmenu"
  @close="closeSubmenu"
  :collapse="isCollapse")

  //- NOTE: keyspace
  el-submenu(
    v-for="v,i in keyspace"
    :key="v.keyspace_name"
    :collapse="isCollapse"
    :index="v.keyspace_name")
    template(slot="title")
      div(@click="goToKeyspace(v.keyspace_name, $event)")
        i(class="el-icon-tickets")
        span(
          style="padding-right: 36px;"
          lot="title") {{v.keyspace_name}}
    el-menu-item-group()
      el-menu-item(
        v-for="vv in v.table"
        @click="goToTable(v.keyspace_name, vv.table_name)"
        :key="`${v.keyspace_name}.${vv.table_name}`"
        :index="`${v.keyspace_name}.${vv.table_name}`")
        template(slot="title")
          i(v-bind:class="[vv.views ? 'el-icon-view' : 'el-icon-document']")
          span(
          lot="title") {{vv.table_name}}

  //- NOTE: query
  el-menu-item(
    @click="openQuery"
    index="-1")
      i(class="el-icon-edit")
      span(lot="title") Query

  //- NOTE: query
  el-menu-item(
    @click="goHostInfo"
    index="-2")
      i(class="el-icon-s-platform")
      span(lot="title") HostInfo

  //- NOTE: folding menu
  el-menu-item(
    @click="toggleMenu"
    index="-3")
      i(v-bind:class="[isCollapse ? 'el-icon-caret-right' : 'el-icon-caret-left']")
      span(lot="title") Folding menu
</template>
<style>
</style>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'side-bar',

  data() {
    return {
      isCollapse: false,
      keyspace: [],
      keyspaceOpen: {},
      loading: false,
    }
  },

  created() {
    this.fetch()
  },

  methods: {
    async fetch() {
      this.loading = true

      try {
        const res = await service.request('keyspace')

        this.keyspace = res.get()
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          message: error
        });
      }

      this.loading = false
    },

    goToTable(keyspaceName, tableName) {
      this.$router.push({
        path: `/main/keyspace/${keyspaceName}/${tableName}/rowstoken/50`,
      }).catch(err => err)
    },

    async openSubmenu(keyspace) {
      if (!this.keyspaceOpen[keyspace]) {
        this.keyspaceOpen[keyspace] = true
      }
    },

    async closeSubmenu(keyspace) {
      if (this.keyspaceOpen[keyspace]) {
        this.keyspaceOpen[keyspace] = false
      }
    },

    async goToKeyspace(keyspace, event) {
      this.$router.push({
        name: 'table-list',
        params: {
          keyspace
        }
      }).catch(err => err)

      if (this.keyspaceOpen[keyspace]) {
        event.stopPropagation()
      }
    },

    async goHostInfo() {
      this.$router.push({
        name: 'hostinfo'
      }).catch(err => err)
    },

    openQuery() {
      this.$prompt('Enter Query', 'CQL Query', {
        confirmButtonText: 'Execute',
        cancelButtonText: 'Cancel',
        inputType: 'textarea',
      }).then(async ({ value }) => {
        this.$router.push({
          name: 'query',
          params: {
            query: `${value}`
          }
        }).catch(err => err)
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Cancel Query'
        })
      })
    },

    toggleMenu() {
      this.isCollapse = !this.isCollapse
    },
  }
};
</script>
