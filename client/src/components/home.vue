<template lang="pug">
  div
    div(class="ca-content")
      div(class="ca-bg")
      div(class="ca-box")
        div(class="ca-a") Cassandra-Web

    div(class="ca-content")
      router-link(
        class="ca-link"
        to="/main")
</template>
<style>
  .ca-bg {
    width: 100%;
    position: relative;
    display: inline-block;
    height: 100%;
  }

  .ca-link {
    width: 100%;
    position:relative;
    display: inline-block;
    height: 100%;
  }

  .ca-content {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: auto;
    width:100%;
    height: 100%;
  }

  .ca-box {
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

import Cookies from 'js-cookie'

export default {
  name: 'Home',
  data() {
    return {}
  },
  created() {
    document.title = 'Cassandra-Web'

    const ct = Cookies.get('theme')

    if (ct === 'dark') {
      this.theme = 'dark'
      document.getElementById('caw-html').className = 'custom-dark';
    } else if (ct === 'white') {
      this.theme = 'white'
      document.getElementById('caw-html').className = '';
    }
  },
  methods: {}
};
</script>
