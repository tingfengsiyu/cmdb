<template>
  <div style=" color: white" >{{ userName }}
      <a-button type="danger" @click="loginOut">退出</a-button>
  </div>
</template>

<script>
export default {
  methods: {
    loginOut() {
      localStorage.clear('token')
      this.$router.push('/login')
    },
  },
  computed: {
    userName() {
      let jwt=`Bearer ${localStorage.getItem('token')}`
      let arr=jwt.split('.')
      let data =JSON.parse(atob(arr[1]))
      return data.username
    }
  },
}
</script>
<style scoped>

</style>
