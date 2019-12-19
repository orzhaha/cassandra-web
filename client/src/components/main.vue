<template lang="pug">
el-container
  el-container(:style="`height: ${wh}px; border: 1px solid #eee`")
    el-aside(style="width: auto")

      //- NOTE: side bar
      side-bar-component(:style="`min-height: ${wh}px;`")

    el-main
       router-view(:key="$route.fullPath")
  el-footer
    //- NOTE: 主題
    //- theme-picker-component
    br
    el-switch(
      style="display: block"
      v-model="theme"
      @change="changeThemeColor"
      active-color="#FFF"
      inactive-color="#222933"
      active-value="dark"
      inactive-value="white"
      active-text="dark"
      inactive-text="white")
</el-switch>
</template>
<style>
</style>

<script>
import Cookies from 'js-cookie'
import ThemePicker from '@/components/theme-picker'
import SideBar from './side-bar'

export default {
  name: 'Main',
  data() {
    return {
      wh: window.innerHeight - 70,
      theme: 'white'
    }
  },
  created() {
    const theme = Cookies.get('theme')

    if (theme === 'dark') {
      this.theme = 'dark'
      document.getElementById('caw-html').className = 'custom-dark';
    } else if (theme === 'white') {
      this.theme = 'white'
      document.getElementById('caw-html').className = '';
    }
  },
  components: {
    'side-bar-component': SideBar,
    'theme-picker-component': ThemePicker
  },
  methods: {
    changeThemeColor(theme) {
      if (theme === 'dark') {
        Cookies.set('theme', 'dark')
      } else if (theme === 'white') {
        Cookies.set('theme', 'white')
      }

      const ct = Cookies.get('theme')

      if (ct === 'dark') {
        this.theme = 'dark'
        document.getElementById('caw-html').className = 'custom-dark';
      } else if (ct === 'white') {
        this.theme = 'white'
        document.getElementById('caw-html').className = '';
      }
    },
  }
};
</script>
