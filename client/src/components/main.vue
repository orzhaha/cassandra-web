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
import Cookies from 'js-cookie'
import ThemePicker from '@/components/theme-picker'
import SideBar from './side-bar'

export default {
  name: 'Main',
  created() {
    if (this.$route.query.theme === 'dark') {
      Cookies.set('theme', 'dark')
    } else if (this.$route.query.theme === 'white') {
      Cookies.set('theme', 'white')
    }

    const theme = Cookies.get('theme')

    if (theme === 'dark') {
      document.getElementById('caw-html').className = 'custom-dark';
    } else if (theme === 'white') {
      document.getElementById('caw-html').className = '';
    }
  },
  components: {
    'side-bar-component': SideBar,
    'theme-picker-component': ThemePicker
  }
};
</script>
