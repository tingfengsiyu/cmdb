<template>
  <div>
    <a-card>
      <h3>{{ id ? '编辑服务器':'新增服务器'}}</h3>

      <a-form-model
          :model="serverInfo"
          ref="serverInfoRef"
          :rules="serverInfoRules"
          :hideRequiredMark="true"
      >
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-model-item label="服务器名称" prop="name">
              <a-input style="width:300px" v-model="serverInfo.name"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器型号" >
              <a-input style="width:300px" v-model="serverInfo.models"></a-input>
            </a-form-model-item>
            <a-form-model-item label="机架位置" >
              <a-input style="width:300px" v-model="serverInfo.location"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器ip" prop="name">
              <a-input style="width:300px" v-model="serverInfo.private_ip_address"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器公网ip" >
              <a-input style="width:300px" v-model="serverInfo.public_ip_address"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器标签ip" prop="name">
              <a-input style="width:300px" v-model="serverInfo.label_ip_address"></a-input>
            </a-form-model-item>
            <a-form-model-item label="cpu">
              <a-input style="width:300px" v-model="serverInfo.cpu"></a-input>
            </a-form-model-item>
            <a-form-model-item label="内存">
              <a-input style="width:300px" v-model="serverInfo.memory"></a-input>
            </a-form-model-item>
          </a-col>
          <a-col :span="4">
            <a-form-model-item label="角色标签"  prop="name">
              <a-input style="width:300px" v-model="serverInfo.label"></a-input>
            </a-form-model-item>
            <a-form-model-item label="磁盘">
              <a-input style="width:300px" v-model="serverInfo.disk"></a-input>
            </a-form-model-item>
            <a-form-model-item label="所属用户" prop="user">
              <a-input style="width:300px" v-model="serverInfo.user"></a-input>
            </a-form-model-item>
            <a-form-model-item label="所属集群" prop="name">
              <a-input style="width:300px" v-model="serverInfo.cluster"></a-input>
            </a-form-model-item>
            <a-form-model-item label="上架状态" prop="name">
              <a-input style="width:300px" v-model="serverInfo.state"></a-input>
            </a-form-model-item>
            <a-form-model-item label="城市" prop="name" >
              <a-input style="width:300px" v-model="serverInfo.city"></a-input>
            </a-form-model-item>
            <a-form-model-item label="机房名" prop="name" >
              <a-input style="width:300px" v-model="serverInfo.idc_name"></a-input>
            </a-form-model-item>
            <a-form-model-item label="机柜号" prop="name" >
              <a-input style="width:300px" v-model="serverInfo.cabinet_number"></a-input>
            </a-form-model-item>

          </a-col>
        </a-row>
        <a-form-model-item>
          <a-button
              type="danger"
              style="margin-right:15px"
              @click.once="artOk(serverInfo.id)"
          >{{serverInfo.id ? '更新':"提交"}}</a-button>
          <a-button type="primary" @click.once="addCancel">取消</a-button>
        </a-form-model-item>
      </a-form-model>
    </a-card>
  </div>
</template>

<script>
import { Url } from '../../plugin/http'
export default {
  props: ['id'],
  data() {
    return {
      serverInfo:{
        id: 0,
        name: '',
        models: '',
        location: '',
        private_ip_address: '',
        public_ip_address: '',
        label: '',
        cluster: '',
        label_ip_address: '',
        cpu: '',
        memory: '',
        disk: '',
        user: '',
        state: '',
        city: '',
        idc_name: '',
        cabinet_number: ''
      },
      Catelist: [],
      upUrl: Url + 'upload',
      headers: {},
      fileList: [],
      serverInfoRules: {
        name: [{ required: true, message: '请输入服务器信息', trigger: 'blur' }],
        user: [{ required: true, message: '请输入机器所属用户', trigger: 'blur' }]
      },
    }
  },
  created() {
    this.headers = { Authorization: `Bearer ${localStorage.getItem('token')}` }
    if (this.id) {
      this.getserverInfo(this.id)
    }
  },
  methods: {
    // 查询服务器信息
    async getserverInfo(id) {
        let { data: res }  = await this.$http.get(`idc/getserver/${id}`)
        if (res.status !== 200) return this.$message.error(res.message)
        let data = {}
        if(Array.isArray(res.data)){
          data = res.data[0]
        }else{
          data = res.data
        }
        this.serverInfo = data
        this.serverInfo.id = data.id
    },
    artOk(id) {
      this.$refs.serverInfoRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入服务器内容')
        if (id === 0) {
          const { data: res } = await this.$http.post('idc/batchcreateserver', JSON.stringify([this.serverInfo]))
          if (res.status !== 200) return this.$message.error(res.message)
          this.$router.push('/listserver')
          this.$message.success('添加服务器成功')
        } else {
          const { data: res } = await this.$http.post(`idc/batchupdateserver`, JSON.stringify([this.serverInfo]))
          //console.log(JSON.stringify([this.serverInfo]))
          console.log(this.serverInfo)
          if (res.status !== 200) return this.$message.error(res.message)

          this.$router.push('/listserver')
          this.$message.success('更新服务器成功')
        }
      })
    },

    addCancel() {
      this.$refs.serverInfoRef.resetFields()
    },
  },
}
</script>
