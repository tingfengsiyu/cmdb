<template>
  <div>
    <a-card>
      <h3>运行ansible playbook</h3>

      <a-form-model
          :model="ops"
          ref="opsRef"
          :rules="opsRules"
          :hideRequiredMark="true"
      >
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-model-item label="连续ip">
              <a-input style="width:300px" v-model="ops.continuous_ip"></a-input>
            </a-form-model-item>
            <a-form-model-item has-feedback label="不连续ip" >
              <a-input style="width:300px" v-model="ops.discontinuous_ip" />
            </a-form-model-item>
            <a-form-model-item label="目标主机组" >
              <a-select placeholder="请选择组名"  mode="multiple"  style="width:100%" @change="handleChange" >
                <a-select-option v-for="item in ClusterList" :key="item.id" :value="item.cluster" >{{item.cluster}}</a-select-option>
              </a-select>
            </a-form-model-item>
            <a-form-model-item has-feedback label="标签" prop="tag">
              <a-input style="width:300px" v-model="ops.tag" />
            </a-form-model-item>
            <a-form-model-item has-feedback label="变量" prop="variable">
              <a-input style="width:300px" v-model="ops.variable" />
            </a-form-model-item>
            <a-form-model-item label="playbook文件名" prop="filename">
              <a-select placeholder="选择playbook文件名"  style="width:100%" @change="fileChange" >
                <a-select-option v-for="item in FileList" :key="item.id" :value="item.filename" >{{item.filename}}</a-select-option>
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
        group: [],
        continuous_ip: '',
        discontinuous_ip:'',
        variable:'',
        filename: '',
        tag: ''
      },
      ClusterList: [],
      FileList:[],
      headers: {},
      opsRules: {
        tag: [{ required: true, message: '请输入tags', trigger: 'blur' }],
        variable: [{ required: true, message: '请输入ansible变量', trigger: 'blur' }],
        filename: [{ required: true, message: '请选择ansiblePlaybbok文件', trigger: 'blur' }],
      },
    }
  },
  created() {
    this.headers = { Authorization: `Bearer ${localStorage.getItem('token')}` }
    this.getClusterList()
    this.getFileList()
  },
  methods: {
    //获取集群
    async getClusterList() {
      const { data: res } = await this.$http.get('idc/getansiblehosts')
      if (res.status !== 200) return this.$message.error(res.message)
      this.ClusterList = res.data
    },
    async getFileList() {
      const { data: res } = await this.$http.get('idc/getplaybookfiles')
      if (res.status !== 200) return this.$message.error(res.message)
      this.FileList = res.data
    },

    handleChange(value) {
      this.ops.group = value;
    },
    fileChange(value) {
      this.ops.filename=value
    },
    // 提交任务
    artOk() {
      this.$refs.opsRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入内容')
          console.log(this.ops)
          const { data: res } = await this.$http.post('idc/ansibleplaybook', this.ops)
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
