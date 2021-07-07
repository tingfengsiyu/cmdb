<template>
  <a-button type="danger" @click="loginOut">退出</a-button>
</template>

<script>
export default {
  methods:{
    loginOut(){
      localStorage.clear('token')
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>

</style>
