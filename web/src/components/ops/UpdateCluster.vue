<template>
  <div>
    <a-card>
      <h3>更新机器所属集群</h3>

      <a-form-model
          :model="ops"
          ref="opsRef"
          :rules="opsRules"
          :hideRequiredMark="true"
      >
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-model-item label="起始ip" prop="source_start_ip">
              <a-input style="width:300px" v-model="ops.source_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="结束ip位" prop="source_end_number">
              <a-input style="width:300px" v-model="ops.source_end_number"></a-input>
            </a-form-model-item>
            <a-form-model-item label="目标集群" prop="target_cluster_name">
              <a-input style="width:300px" v-model="ops.target_cluster_name"></a-input>
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
        source_end_number: '',
        target_cluster_name: '',
      },
      Catelist: [],
      headers: {},
      fileList: [],
      opsRules: {
        source_start_ip: [{ required: true, message: '请输入源ip', trigger: 'blur' }],
        source_end_number: [{ required: true, message: '请输入源结束ip位，1-255', trigger: 'blur' }],
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
    // 提交任务
    artOk() {
      this.$refs.opsRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入内容')
          const { data: res } = await this.$http.put('idc/updatecluster', JSON.stringify(this.ops))
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
