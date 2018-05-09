<template lang="pug">
  div#app
    router-view(v-loading="true")
</template>

<script>
export default {
  name: 'App',
};
</script>

<style lang="sass">
body
  height: 100%
#app
  font-family: 'Avenir', Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  color: #2c3e50
</style>
