<template lang="pug">
el-container
  el-container(style="height: 100%; border: 1px solid #eee")
    el-aside(style="width: auto")

      //- NOTE: side bar
      side-bar-component

    el-main
       router-view
  el-footer
    //- NOTE: 主題
    theme-picker-component
</template>
<style>
</style>

<script>
import ThemePicker from '@/components/theme-picker'
import SideBar from './side-bar'

export default {
  name: 'Main',

  components: {
    'side-bar-component': SideBar,
    'theme-picker-component': ThemePicker
  }
};
</script>
