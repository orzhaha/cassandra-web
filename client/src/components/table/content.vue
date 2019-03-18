<template lang="pug">
    div.content-title
      span welcome to
      span.keyspaceName {{` ${$route.params.keyspace}`}}
      span.tableName {{` ${$route.params.table}`}}
      span {{'! :)'}}
</template>
<style>
  .content-title {
    margin:70px auto;
    font-size: 40px;
    font-weight: 100;
    font-family: 'Lato';
  }
  .keyspaceName {
    color: #409EFF;
  }
  .tableName {
    color: red;
  }
</style>
<script>
export default {}
</script>
