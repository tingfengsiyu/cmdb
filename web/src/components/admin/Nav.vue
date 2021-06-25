<template >
  <a-layout-sider breakpoint="lg" v-model="collapsed">
    <div class="log">
      <span>CMDB</span>
    </div>

    <a-menu theme="dark" mode="inline" @click="goToPage">
      <a-menu-item key="index">
        <a-icon type="dashboard"></a-icon>
        <span>仪表盘</span>
      </a-menu-item>

      <a-sub-menu>
        <span slot="title">
          <a-icon type="file" />
          <span>服务器管理</span>
        </span>
        <a-menu-item key="addserver">
          <a-icon type="form" />
          <span>新增服务器</span>
        </a-menu-item>
        <a-menu-item key="listserver">
          <a-icon type="snippets" />
          <span>服务器列表</span>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu>
        <span slot="title">
          <a-icon type="file" />
          <span>kubernetes集群管理</span>
        </span>
        <a-menu-item key="addart">
          <a-icon type="form" />
          <span>Pod</span>
        </a-menu-item>
        <a-menu-item key="artlist">
          <a-icon type="snippets" />
          <span>SVC</span>
        </a-menu-item>
      </a-sub-menu>
      <a-menu-item key="userlist">
        <a-icon type="user" />
        <span>用户列表</span>
      </a-menu-item>
      <a-sub-menu>
        <span slot="title">
          <a-icon type="file" />
          <span>运维操作</span>
        </span>
        <a-menu-item key="addart">
          <a-icon type="form" />
          <span>操作记录</span>
        </a-menu-item>
        <a-menu-item key="artlist">
          <a-icon type="snippets" />
          <span>系统初始化</span>
        </a-menu-item>
      </a-sub-menu>
    </a-menu>
  </a-layout-sider>
</template>

<script>
export default {
  data() {
    return {
      collapsed: false,
    }
  },
  methods: {
    goToPage(item) {
      this.$router.push('/' + item.key).catch((err) => err)
    },
  },
}
</script>

<style scoped>
.log {
  height: 32px;
  margin: 16px;
  background-color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 17px;
}
</style>
