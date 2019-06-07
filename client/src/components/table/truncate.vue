<template lang="pug">
    div
      br
      el-button(
        type="danger"
        icon="el-icon-folder-delete"
        @click="handleTruncateConfirm()") Truncate
</template>

<style>
a.download_btn {
  color:#FFFFFF;
  text-decoration:none;
}
.CodeMirror {
  border: 1px solid #eee;
  height: auto;
}
</style>

<script>
import api from '@/api'

const service = api.make('root')

export default {
  name: 'Truncate',
  methods: {
    handleTruncateConfirm() {
      const table = `${this.$route.params.keyspace}.${this.$route.params.table}`
      this.$confirm(`Truncate Table "${table}" Are you sure ?`, '', {
        confirmButtonText: 'Submit',
        cancelButtonText: 'Cancel',
        confirmButtonClass: 'el-button--danger'
      }).then(() => {
        this.TruncateTable(table)
      }).catch(() => {})
    },

    async TruncateTable(table) {
      try {
        const res = await service.request('truncate', {
          data: {
            table: table,
          }
        })

        const message = (res.get() === []) ? 'success' : res.get()

        this.$message({
          type: 'success',
          showClose: true,
          duration: 0,
          message
        });
      } catch (error) {
        this.$message({
          type: 'error',
          showClose: true,
          duration: 0,
          message: error.body.message
        });
      }
    },
  },
};
</script>
