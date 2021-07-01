<template>
  <div>
    <a-card>
      <h3>修改机器ip</h3>

      <a-form-model
          :model="ops"
          ref="opsRef"
          :rules="opsRules"
          :hideRequiredMark="true"
      >
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-model-item label="源ip" prop="source_start_ip">
              <a-input style="width:300px" v-model="ops.source_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="源网关" prop="source_gateway">
              <a-input style="width:300px" v-model="ops.source_gateway"></a-input>
            </a-form-model-item>
            <a-form-model-item label="源结束ip位" prop="source_end_number">
              <a-input style="width:300px" v-model="ops.source_end_number"></a-input>
            </a-form-model-item>
<!--            <a-form-model-item label="源结束ip位" prop="source_end_number">-->
<!--              <a-input-number id="inputNumber"    :min="1" :max="255"  v-model="ops.source_end_number"/>-->
<!--            </a-form-model-item>-->
            <a-form-model-item label="目标ip" prop="target_start_ip">
              <a-input style="width:300px" v-model="ops.target_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="目标网关" prop="target_gateway">
              <a-input style="width:300px" v-model="ops.target_gateway"></a-input>
            </a-form-model-item>
            <a-form-model-item label="选择集群" prop="target_cluster_name">
              <a-input style="width:300px" v-model="ops.target_cluster_name"></a-input>
              <a-select placeholder="请选择集群" style="width:300px">
                <a-select-option v-for="item in ClusterList" :key="item.id" :value="item.cluster">{{item.cluster}}</a-select-option>
              </a-select>
            </a-form-model-item>

          </a-col>
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
      ops:{
        source_start_ip: '',
        source_gateway: '',
        source_end_number: '',
        target_start_ip: '',
        target_gateway: '',
        target_cluster_name: '',
      },
      ClusterList: [],
      headers: {},
      opsRules: {
        source_start_ip: [{ required: true, message: '请输入源ip', trigger: 'blur' }],
        source_gateway: [{ required: true, message: '请输入源网关', trigger: 'blur' }],
        target_start_ip: [{ required: true, message: '请输入目标ip', trigger: 'blur' }],
        target_gateway: [{ required: true, message: '请输入目标网关', trigger: 'blur' }],
        source_end_number: [
            { required: true, message: '请输入结束数,值为0-255', trigger: 'blur'},
            { max: 3, message: '最多可写3个字符', trigger: 'change' },
        ],
        target_cluster_name: [{ required: true, message: '请输入目标集群', trigger: 'blur' }],
      },
    }
  },
  created() {
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
    if (this.id) {
      this.getops(this.id)
    }
  },
  methods: {
    // 获取机器集群
    async getClusterList() {
      const { data: res } = await this.$http.get('idc/getclusters')
      if (res.status !== 200) return this.$message.error(res.message)
      this.ClusterList = res.data
      console.log(res.total)
    },
    // 提交任务
    artOk() {
      this.$refs.opsRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入内容')
          const { data: res } = await this.$http.put('idc/batchip', JSON.stringify(this.ops))
          if (res.status !== 200) return this.$message.error(res.message)
          this.$router.push('/OpsRecords')
          this.$message.success('添加任务成功')
      })
    },

    addCancel() {
      this.$refs.opsRef.resetFields()
    },
  },
}
</script>
