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

      tabs: [
        {
          name: 'rows',
          trigger: this.showRows
        },
        {
          name: 'rowstoken',
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
      })
    },

    RowsToken() {
      this.$router.push({
        name: 'rowstoken',
        params: {
          pagesize: 50
        }
      })
    },

    Find() {
      this.$router.push({
        name: 'find'
      })
    },

    Columns() {
      this.$router.push({
        name: 'columns'
      })
    },

    Definition() {
      this.$router.push({
        name: 'definition'
      })
    },

    Export() {
      this.$router.push({
        name: 'export'
      })
    },

    Import() {
      this.$router.push({
        name: 'import'
      })
    },
  }
}
</script>
