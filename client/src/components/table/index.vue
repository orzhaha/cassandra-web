<template lang="pug">
  div
    el-header
      menu-component

    el-container
      router-view
</template>
<style>
</style>

<script>
import Menu from './menu'

export default {
  created() {
    document.title = 'Index'
  },
  components: {
    'menu-component': Menu
  }
};
</script>
