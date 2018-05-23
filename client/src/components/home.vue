<template lang="pug">
  div#home
    h1 HOME
</template>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Home',
  async created() {
    try {
      const res = await service.request('allKeySpace')

      console.log(res)
    } catch (error) {
      console.log('[error]', error)
    }
  }
};
</script>
