<template>
  <a-layout class="container">
    <a-layout-sider>
<!--      <div class="log">-->
<!--        <span>cmdb</span>-->
<!--      </div>-->

    </a-layout-sider>
    <Nav></Nav>
    <a-layout>
      <a-layout-header class="headerBtn">
        <Header></Header>
      </a-layout-header>
      <a-layout-content>
        <router-view :key="$route.fullPath"></router-view>
      </a-layout-content>
      <a-layout-footer>
        <Footer> </Footer>
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>
<script>
import Nav from '../components/admin/Nav'
import Footer from '../components/admin/Footer'
import Header from '../components/admin/Header'
export default {
  components: {Nav,Footer,Header}
}
</script>
<style>
.container{
  height:  100%;
}
.headerBtn {
 display: flex;
  justify-content: flex-end;
  align-items: center;
}
.log {
  height: 32px;
  margin: 16px;
  background-color: #fff;
}
</style>
