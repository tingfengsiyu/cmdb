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
            <a-form-model-item label="机器开始ip" prop="StorageMount.init_start_ip" :rules="opsRules.init_start_ip">
              <a-input style="width:300px" v-model="ops.StorageMount.init_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="结束ip位" prop="StorageMount.init_end_number" :rules="opsRules.init_end_number">
              <a-input style="width:300px" v-model.number="ops.StorageMount.init_end_number"></a-input>
            </a-form-model-item>
            <a-form-model-item label="需挂载的存储开始ip" prop="StorageMount.storage_start_ip" :rules="opsRules.storage_start_ip">
              <a-input style="width:300px" v-model="ops.StorageMount.storage_start_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item label="存储结束位" prop="StorageMount.storage_stop_number" :rules="opsRules.storage_stop_number">
              <a-input style="width:300px" v-model.number="ops.StorageMount.storage_stop_number"></a-input>
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
    let checkNumber = (rule, value, callback) => {
      if (!value) {
        return callback(new Error('请输入整数值'));
      }else if (!Number.isInteger(value)) {
        callback(new Error('请输入整数'));
      } else {
        if (value < 1||value > 255) {
          callback(new Error('输入必须在0-255'));
        } else {
          callback();
        }
      }
    };
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
        init_start_ip: [{ required: true, message: '初始化机器开始ip', trigger: 'blur' }],
        init_end_number: [{ validator: checkNumber, trigger: 'change' }],
        storage_start_ip: [{ required: true, message: '存储机器开始ip', trigger: 'blur' }],
        storage_stop_number: [{ validator: checkNumber, trigger: 'change' }],
        init_user: [{ required: true, message: '初始化用户名', trigger: 'blur'}],
        init_pass: [{ required: true, message: '初始化用户名的密码', trigger: 'blur'}],
        role: [{ required: true, message: '机器角色 lotus-worker 或 lotus-storage ', trigger: 'blur'}],
      },
    }
  },
  created() {
    this.headers = { Authorization: `Bearer ${localStorage.getItem('token')}` }
  },
  methods: {
    // 提交任务
    artOk() {
      this.$refs.opsRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入内容')
        let ops = { ...this.ops }
        ops.StorageMount.init_end_number=String(ops.StorageMount.init_end_number)
        ops.StorageMount.storage_stop_number=String(ops.StorageMount.storage_stop_number)
        const { data: res } = await this.$http.post('idc/shellosinit', JSON.stringify(ops))
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
