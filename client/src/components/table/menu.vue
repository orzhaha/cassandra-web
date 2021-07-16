<template lang="pug">
  el-menu.el-menu-demo(
    :default-active="defaultActive"
    mode="horizontal"
    @select="handleSelect")

    el-menu-item(
      v-for="(tab, index) in tabs"
      :index="tab.name"
      :key="tab.name") {{tab.name}}
</template>

<style>
</style>

<script>
import get from 'lodash/get'
import find from 'lodash/find'
import split from 'lodash/split'

export default {
  data() {
    const arr = split(this.$route.path, '/')

    return {
      defaultActive: get(arr, 5),

      tabs: this.$readonly ? [
        {
          name: 'rows',
          trigger: this.RowsToken
        },
        {
          name: 'find',
          trigger: this.Find
        },
        {
          name: 'columns',
          trigger: this.Columns
        },
        {
          name: 'definition',
          trigger: this.Definition
        },
        {
          name: 'export',
          trigger: this.Export
        }
      ] : [
        {
          name: 'rows',
          trigger: this.RowsToken
        },
        {
          name: 'find',
          trigger: this.Find
        },
        {
          name: 'columns',
          trigger: this.Columns
        },
        {
          name: 'definition',
          trigger: this.Definition
        },
        {
          name: 'export',
          trigger: this.Export
        },
        {
          name: 'import',
          trigger: this.Import
        },
        {
          name: 'truncate',
          trigger: this.Truncate
        }
      ]
    }
  },

  methods: {
    handleSelect(name) {
      const trigger = get(find(this.tabs, tab => tab.name === name), 'trigger')

      if (!trigger) {
        return
      }

      trigger()
    },

    showRows() {
      this.$router.push({
        name: 'rows',
        params: {
          page: 1,
          pagesize: 50
        }
      }).catch(err => err)
    },

    RowsToken() {
      this.$router.push({
        name: 'rowstoken',
        params: {
          pagesize: 50
        }
      }).catch(err => err)
    },

    Find() {
      this.$router.push({
        name: 'find'
      }).catch(err => err)
    },

    Columns() {
      this.$router.push({
        name: 'columns'
      }).catch(err => err)
    },

    Definition() {
      this.$router.push({
        name: 'definition'
      }).catch(err => err)
    },

    Export() {
      this.$router.push({
        name: 'export'
      }).catch(err => err)
    },

    Import() {
      this.$router.push({
        name: 'import'
      }).catch(err => err)
    },

    Truncate() {
      this.$router.push({
        name: 'truncate'
      }).catch(err => err)
    },
  }
}
</script>
