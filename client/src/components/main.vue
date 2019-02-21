<template lang="pug">
  el-container(style="height: 100%; border: 1px solid #eee")

    el-aside(style="width: auto")
      //- NOTE: side bar
      side-bar-component

    el-main
       router-view
</template>
<style>
</style>

<script>
import SideBar from './side-bar'

export default {
  name: 'Main',

  components: {
    'side-bar-component': SideBar
  }
};
</script>
