<template>
  <div>
    <a-card>
      <h3>机器初始化</h3>

      <a-form-model
          :model="ops"
          ref="opsRef"
          :rules="opsRules"
          :hideRequiredMark="true"
      >
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-model-item label="机器开始ip" prop="init_start_ip">
              <a-input style="width:300px" v-model="ops.init_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="结束ip位" prop="init_end_number">
              <a-input style="width:300px" v-model="ops.init_end_number"></a-input>
            </a-form-model-item>
            <!--            <a-form-model-item label="源结束ip位" prop="source_end_number">-->
            <!--              <a-input-number id="inputNumber"    :min="1" :max="255"  v-model="ops.source_end_number"/>-->
            <!--            </a-form-model-item>-->
            <a-form-model-item label="存储开始ip" prop="storage_start_ip">
              <a-input style="width:300px" v-model="ops.storage_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="存储结束位" prop="storage_stop_number">
              <a-input style="width:300px" v-model="ops.storage_stop_number"></a-input>
            </a-form-model-item>
            <a-form-model-item label="初始化用户名" prop="init_user">
              <a-input style="width:300px" v-model="ops.init_user"></a-input>
            </a-form-model-item>
            <a-form-model-item label="初始化用户密码" prop="init_pass">
              <a-input style="width:300px" v-model="ops.init_pass"></a-input>
            </a-form-model-item>
            <a-form-model-item label="初始化机器角色" prop="role">
              <a-input style="width:300px" v-model="ops.role"></a-input>
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
        StorageMount: {
          init_start_ip: "",
          init_end_number: "",
          storage_start_ip:"",
          storage_stop_number:""
        },
        init_user:"",
        init_pass:"",
        role: ""
      },
      ClusterList: [],
      opsRules: {
        init_start_ip: [{ required: true, message: '挂载机器开始ip', trigger: 'blur' }],
        init_end_number: [{ required: true, message: '挂载机器结束ip位,值为1-253', trigger: 'blur' }],
        storage_start_ip: [{ required: true, message: '存储机器开始ip', trigger: 'blur' }],
        storage_stop_number: [{ required: true, message: '挂载机器结束ip位,值为1-253', trigger: 'blur' }],
        init_user: [{ required: true, message: '初始化用户名', trigger: 'blur'}],
        init_pass: [{ required: true, message: '初始化用户名的密码', trigger: 'blur'}],
        role: [{ required: true, message: '机器角色 lotus-worker 或 lotus-storage ', trigger: 'blur'}],
      },
    }
  },
  created() {
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
  },
  methods: {
    // 提交任务
    artOk() {
      this.$refs.opsRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入内容')
        const { data: res } = await this.$http.put('idc/shellosinit', JSON.stringify(this.ops))
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
