<template lang="pug">
  el-container(style="height: 600px; border: 1px solid #eee")
    el-aside(style="width: auto")
      //- el-radio-group(v-model="isCollapse"
      //-   style="margin-bottom: 20px;"
      //-   size="mini")
      //-   el-radio-button(:label="false") 展开
      //-   el-radio-button(:label="true") 收起
      el-menu(default-active="2"
        class="el-menu-vertical-demo"
        :collapse-transition="false"
        :collapse="isCollapse")
        el-menu-item(
        v-for="v,i in keyspace"
        @click="getTable(v.keyspace_name)"
        :index="i.toString()")
          i(class="el-icon-tickets")
          span(lot="title") {{v.keyspace_name}}
        el-menu-item(
        @click="openQuery"
        index="-1")
          i(class="el-icon-edit")
          span(lot="title") Query
        el-menu-item(
        @click="toggleMenu"
        index="-2")
          i(v-bind:class="[isCollapse ? 'el-icon-caret-right' : 'el-icon-caret-left']")
          span(lot="title") Folding menu
    el-container
      router-view(v-loading="false")
</template>
<style>
  .el-header {
    background-color: #B3C0D1;
    color: #333;
    line-height: 60px;
  }

  .el-aside {
    color: #333;
  }
</style>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Keyspace',
  data() {
    return {
      isCollapse: false,
      keyspace: [],
    }
  },
  created() {
    this.fetch()
  },
  methods: {
    async getTable(keyspace) {
      try {
        this.$router.push({
          name: 'table',
          params: {
            keyspace
          }
        })
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error
        });
      }
    },
    async fetch() {
      try {
        const res = await service.request('keyspace')
        this.keyspace = res.get()
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error
        });
      }
    },
    toggleMenu() {
      this.isCollapse = !this.isCollapse
    },
    openQuery() {
      this.$prompt('Enter Query', 'CQL Query', {
        confirmButtonText: 'Execute',
        cancelButtonText: 'Cancel',
      }).then(async ({ value }) => {
        try {
          const res = await service.request('query', {
            data: {
              query: value
            }
          })
          this.$message({
            type: 'success',
            showClose: true,
            duration: 0,
            message: res.get()
          });
        } catch (error) {
          this.$message({
            type: 'error',
            showClose: true,
            duration: 0,
            message: error
          });
        }
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Cancel Query'
        });
      });
    }
  }
};
</script>
