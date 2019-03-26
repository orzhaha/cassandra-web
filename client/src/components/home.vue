<template lang="pug">
  div(class="ca-content")
    router-link(
      class="ca-a"
      to="/main") Cassandra-Web
</template>
<style>
  .ca-content {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: auto;
    width:702px;
    height: 135px;
  }

  .ca-a {
    width: 100%;
    font-size: 96px;
    margin-bottom: 40px;
    font-weight: 100;
    font-family: 'Lato';
    color: #B0BEC5;
    text-decoration:none;
  }
</style>

<script>

export default {
  name: 'Home',
  data() {
    return {}
  },
  created() {},
  methods: {}
};
</script>
