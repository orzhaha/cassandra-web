<template lang="pug">
  router-link(to="/keyspace") KeySpace
</template>
<style>
</style>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Home',
  data() {
    return {}
  },
  created() {},
  methods: {}
};
</script>
