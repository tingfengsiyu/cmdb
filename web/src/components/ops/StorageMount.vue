<template>
  <div>
    <a-card>
      <h3>存储挂载</h3>

      <a-form-model
          :model="ops"
          ref="opsRef"
          :rules="opsRules"
          :hideRequiredMark="true"
      >
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-model-item label="源ip" prop="init_start_ip">
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
            <a-form-model-item label="执行操作" prop="operating">
              <a-input style="width:300px" v-model="ops.operating"></a-input>
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
        init_start_ip: '',
        init_end_number: '',
        storage_start_ip: '',
        storage_stop_number: '',
        operating: '',
      },
      ClusterList: [],
      opsRules: {
        init_start_ip: [{ required: true, message: '请输入挂载机器开始ip', trigger: 'blur' }],
        init_end_number: [{ required: true, message: '请输入挂载机器结束ip位,值为1-253', trigger: 'blur' },{
          max: 3, message: '最多可写3个字符', trigger: 'change'
        }],
        storage_start_ip: [{ required: true, message: '请输入存储机器开始ip', trigger: 'blur' }],
        storage_stop_number: [{ required: true, message: '请输入挂载机器结束ip位,值为1-253', trigger: 'blur' },{
          max: 3, message: '最多可写3个字符', trigger: 'change'
        }],
        operating: [
          { required: true, message: '请输入执行的操作,值为all|add  all，全新，add追加', trigger: 'blur'},
          { max: 3, message: '最多可写3个字符', trigger: 'change' },
        ],
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
        const { data: res } = await this.$http.post('idc/storagemount', JSON.stringify(this.ops))
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
