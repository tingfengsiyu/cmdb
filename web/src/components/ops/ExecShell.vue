<template>
  <div>
    <a-card>
      <h3>执行shell命令</h3>

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
            <a-form-model-item has-feedback label="结束ip位" prop="source_end_number">
              <a-input style="width:300px" v-model.number="ops.source_end_number" />
            </a-form-model-item>
            <a-form-model-item label="执行命令" >
              <a-input style="width:90%; height: 100%"  v-model="ops.cmd" type="textarea" ></a-input>
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
// import Editor from '../editor/index'
export default {
  // components: { Editor },
  data() {
    let checkNumber = (rule, value, callback) => {
      if (!value) {
        return callback(new Error('请输入'));
      }else if (!Number.isInteger(value)) {
        callback(new Error('请输入整数'));
      } else {
        if (value < 0||value > 255) {
          callback(new Error('输入必须在0-255'));
        } else {
          callback();
        }
      }
    };
    return {
      ops:{
        source_start_ip: '',
        source_end_number: '',
        cmd: '',
      },
      ClusterList: [],
      headers: {},
      opsRules: {
        source_start_ip: [{ required: true, message: '请输入源ip', trigger: 'blur' }],
        source_end_number: [{ validator: checkNumber, trigger: 'change' }],
        cmd: [{ required: true, message: '请输入要执行的命令', trigger: 'blur' }],
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
        ops.source_end_number=String(ops.source_end_number)
          const { data: res } = await this.$http.put('idc/execshell', JSON.stringify(ops))
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
