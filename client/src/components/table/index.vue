<template lang="pug">
  div
    menu-component
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
