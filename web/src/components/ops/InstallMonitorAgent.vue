<template>
  <div>
    <a-card>
      <h3>安装集群监控agent</h3>

      <a-form-model
          ref="opsRef"
          :hideRequiredMark="true"
      >
        <a-select placeholder="请选择集群"  style="width:200px" @change="artOk" >
          <a-select-option v-for="item in ClusterList" :key="item.id" :value="item.cluster" >{{item.cluster}}</a-select-option>
        </a-select>

        <a-row :gutter="24">
            <a-form-model-item label="集群名" >
              <a-select placeholder="请选择集群"  style="width:200px" @change="handleChange" >
                <a-select-option v-for="item in ClusterList" :key="item.id" :value="item.cluster" >{{item.cluster}}</a-select-option>
              </a-select>
            </a-form-model-item>
        </a-row>
        <a-form-model-item>
          <a-button
              type="danger"
              style="margin-right:15px"
              @click.once="artOk"
          >"提交"</a-button>
          <a-button type="primary" @click.once="addCancel">取消</a-button>
        </a-form-model-item>
      </a-form-model>
    </a-card>
  </div>
</template>

<script>
export default {
  data() {
    return {
      ClusterList: [],
      queryParam: {
        cluster:'',
      },
    }
  },
  created() {
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
    this.getClusterList()
  },
  methods: {
    //获取集群
    async getClusterList() {
      const { data: res } = await this.$http.get('idc/getclusters')
      if (res.status !== 200) return this.$message.error(res.message)
      this.ClusterList = res.data
    },

    handleChange(value) {
      this.queryParam.cluster = value;
      this.$http.get(`idc/installmointoragent`, {
            params: {
              clustername: value
            },})
    },
    // 提交任务
    artOk() {
      this.$refs.opsRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入内容')
        const { data: res } = await this.$http.post(`idc/installagent`, {
          params: {
            clustername: this.queryParam.cluster
          },

        })
          console.log(this.queryParam.cluster)
          if (res.status !== 200) return this.$message.error(res.message)
          this.$router.push('/OpsRecords')
          this.$message.success(res.message)
      })
    },

    addCancel() {
      this.$refs.opsRef.resetFields()
    },
  },
}
</script>
