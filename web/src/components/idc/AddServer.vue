<template>
  <div>
    <a-card>
      <h3>{{id? '编辑服务器':'新增服务器'}}</h3>

      <a-form-model
          :model="serverInfo"
          ref="serverInfoRef"
          :rules="serverInfoRules"
          :hideRequiredMark="true"
      >
        <a-row gutter="24">
          <a-col :span="10">
            <a-form-model-item label="服务器名称" prop="title">
              <a-input style="width:300px" v-model="serverInfo.name"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器型号" prop="model">
              <a-input style="width:300px" v-model="serverInfo.models"></a-input>
            </a-form-model-item>
            <a-form-model-item label="机架位置" prop="location">
              <a-input style="width:300px" v-model="serverInfo.location"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器ip" prop="private_ip_address">
              <a-input style="width:300px" v-model="serverInfo.private_ip_address"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器公网ip" >
              <a-input style="width:300px" v-model="serverInfo.public_ip_address"></a-input>
            </a-form-model-item>
            <a-form-model-item label="服务器标签ip" prop="label_ip_address">
              <a-input style="width:300px" v-model="serverInfo.label_ip_address"></a-input>
            </a-form-model-item>
            <a-form-model-item label="cpu" prop="title">
              <a-input style="width:300px" v-model="serverInfo.cpu"></a-input>
            </a-form-model-item>
            <a-form-model-item label="内存" prop="title">
              <a-input style="width:300px" v-model="serverInfo.memory"></a-input>
            </a-form-model-item>
          </a-col>
          <a-col :span="4">
            <a-form-model-item label="磁盘" prop="title">
              <a-input style="width:300px" v-model="serverInfo.disk"></a-input>
            </a-form-model-item>
            <a-form-model-item label="所属用户" prop="title">
              <a-input style="width:300px" v-model="serverInfo.user"></a-input>
            </a-form-model-item>
            <a-form-model-item label="所属集群" prop="title">
              <a-input style="width:300px" v-model="serverInfo.cluster"></a-input>
            </a-form-model-item>
            <a-form-model-item label="上架状态" prop="title">
              <a-input style="width:300px" v-model="serverInfo.state"></a-input>
            </a-form-model-item>
            <a-form-model-item label="城市" prop="title">
              <a-input style="width:300px" v-model="serverInfo.city"></a-input>
            </a-form-model-item>
            <a-form-model-item label="机房名" prop="title">
              <a-input style="width:300px" v-model="serverInfo.idc_name"></a-input>
            </a-form-model-item>
            <a-form-model-item label="机柜号" prop="title">
              <a-input style="width:300px" v-model="serverInfo.cabinet_number"></a-input>
            </a-form-model-item>

          </a-col>
        </a-row>
<!--        <a-form-model-item label="服务器内容" prop="content">-->
<!--          <Editor v-model="serverInfo.content"></Editor>-->
<!--        </a-form-model-item>-->

        <a-form-model-item>
          <a-button
              type="danger"
              style="margin-right:15px"
              @click.once="artOk(serverInfo.id)"
          >{{serverInfo.id?'更新':"提交"}}</a-button>
          <a-button type="primary" @click.once="addCancel">取消</a-button>
        </a-form-model-item>
      </a-form-model>
    </a-card>
  </div>
</template>

<script>
import { Url } from '../../plugin/http'
//import Editor from '../editor/index'
export default {
  //components: { Editor },
  props: ['id'],
  data() {
    return {
      serverInfo: [{
        id: 0,
       name:'',
        models:'',
        location:'',
        private_ip_address:'',
        public_ip_address:'',
        label:'',
        cluster:'',
        label_ip_address:'',
        cpu:'',
        memory:'',
        disk:'',
        user:'',
        state:'',
        city:'',
        idc_name:'',
        cabinet_number:''
      }],
      Catelist: [],
      upUrl: Url + 'upload',
      headers: {},
      fileList: [],
      serverInfoRules: {
        title: [{ required: true, message: '请输入服务器标题', trigger: 'blur' }],
        cid: [{ required: true, message: '请选择服务器分类', trigger: 'change' }],
        desc: [
          { required: true, message: '请输入服务器描述', trigger: 'blur' },
          { max: 120, message: '描述最多可写120个字符', trigger: 'change' },
        ],
        img: [{ required: true, message: '请选择服务器缩略图', trigger: 'blur' }],
        content: [{ required: true, message: '请输入服务器内容', trigger: 'blur' }],
      },
    }
  },
  created() {
    this.getCateList()
    this.headers = { Authorization: `Bearer ${window.sessionStorage.getItem('token')}` }
    if (this.id) {
      this.getserverInfo(this.id)
    }
  },
  methods: {
    // 查询服务器信息
    async getserverInfo(id) {
      const { data: res } = await this.$http.get(`article/info/${id}`)
      if (res.status !== 200) return this.$message.error(res.message)
      this.serverInfo = res.data
      this.serverInfo.id = res.data.ID
    },
    // 获取分类列表
    async getCateList() {
      const { data: res } = await this.$http.get('idc/getclusters')
      if (res.status !== 200) return this.$message.error(res.message)
      this.Catelist = res.data
    },
    // 选择分类
    cateChange(value) {
      this.serverInfo.cid = value
    },
    artOk(id) {
      this.$refs.serverInfoRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数验证未通过，请按要求录入服务器内容')
        if (id === 0) {
          const { data: res } = await this.$http.post('article/add', this.serverInfo)
          if (res.status !== 200) return this.$message.error(res.message)
          this.$router.push('/artlist')
          this.$message.success('添加服务器成功')
        } else {
          const { data: res } = await this.$http.put(`article/${id}`, this.serverInfo)
          if (res.status !== 200) return this.$message.error(res.message)

          this.$router.push('/artlist')
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
