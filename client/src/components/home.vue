<template lang="pug">
  div(class="container")
    div(class="content")
      div(class="title")
        router-link(to="/main") Cassandra-Web
</template>
<style>
  html {
    height: 100%;
  }
  body {
    margin: 0;
    padding: 0;
    width: 100%;
    height: 100%;
    color: #B0BEC5;
    display: table;
  }
  .container {
    text-align: center;
    display: table-cell;
    vertical-align: middle;
  }
  .content {
    text-align: center;
    display: inline-block;
  }
  .title {
    font-size: 96px;
    margin-bottom: 40px;
    font-weight: 100;
    font-family: 'Lato';
  }
  a {
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
