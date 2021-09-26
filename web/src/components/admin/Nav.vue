<template >
  <a-layout-sider breakpoint="lg" v-model="collapsed">
    <div class="log">
      <span>CMDB</span>
    </div>

    <a-menu theme="dark" mode="inline" @click="goToPage">

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
           
          <span>服务器列表</span>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu>
        <span slot="title">
          <a-icon type="cloud" />
          <span>K8S集群管理</span>
        </span>
        <a-menu-item key="addart">
          <span>Pod</span>
        </a-menu-item>
        <a-menu-item key="artlist">
          <span>SVC</span>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu>
        <span slot="title">
          <a-icon type="user" />
          <span>权限管理</span>
        </span>
        <a-menu-item key="userlist">
          <span>系统用户</span>
        </a-menu-item>
        <a-menu-item key="termuserlist">
          <span>终端用户</span>
        </a-menu-item>
        <a-menu-item key="termPermissionsList">
          <span>资产授权</span>
        </a-menu-item>
      </a-sub-menu>

      <a-sub-menu>
        <span slot="title">
          <a-icon type="file" />
          <span>运维操作</span>
        </span>
        <a-menu-item key="OpsRecords">
          <span>操作记录</span>
        </a-menu-item>
        <a-menu-item key="osinit">
          <span>系统初始化</span>
        </a-menu-item>
        <a-menu-item key="batchip">
          <span>修改ip</span>
        </a-menu-item>
        <a-menu-item key="storagemount">
          <span>挂载存储</span>
        </a-menu-item>
        <a-menu-item key="installagent">
          <span>安装监控agent</span>
        </a-menu-item>
        <a-menu-item key="UpdateCluster">
          <span>修改机器所属集群</span>
        </a-menu-item>
        <a-menu-item key="ExecShell">
          <a-icon type="code" />
          <span>执行shell命令</span>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu>
        <span slot="title">
          <a-icon type="code" />
          <span>终端操作</span>
        </span>
        <a-menu-item key="addserver">
          <span>命令记录</span>
        </a-menu-item>
      </a-sub-menu>
      <a-sub-menu>
        <span slot="title">
          <a-icon type="area-chart" />
          <span>监控管理</span>
        </span>
        <a-menu-item key="alert">
          <a-icon type="alert" />
          <span>业务告警展示</span>
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
