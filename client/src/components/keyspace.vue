<template lang="pug">
  el-container(style="height: 600px; border: 1px solid #eee")
    el-aside(style="width: auto")
      el-radio-group(v-model="isCollapse"
        style="margin-bottom: 20px;"
        size="mini")
        el-radio-button(:label="false") 展开
        el-radio-button(:label="true") 收起
      el-menu(default-active="2"
        class="el-menu-vertical-demo"
        @open="handleOpen"
        @close="handleClose"
        :collapse="isCollapse")
        el-menu-item(
        v-for="v,i in keyspace"
        @click="getTable(v.keyspace_name)"
        :index="i.toString()")
          i(class="el-icon-tickets")
          span(lot="title") {{v.keyspace_name}}
    el-container
      router-view(v-loading="false")
      //- Table(:table_data="table_data"
      //-   :table_field="table_field")
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
      }
    },
    async fetch() {
      try {
        const res = await service.request('keyspace')
        this.keyspace = res.get()
      } catch (error) {
      }
    },
    handleOpen(key, keyPath) {
      console.log(key, keyPath);
    },
    handleClose(key, keyPath) {
      console.log(key, keyPath);
    }
  }
};
</script>
