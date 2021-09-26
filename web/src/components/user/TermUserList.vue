<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
              v-model="queryParam.username"
              placeholder="输入终端用户名查找"
              enter-button
              allowClear
              @search="getUserList"
          />
        </a-col>
        <a-col :span="4">
          <a-button type="primary" @click="addUserVisible = true">新增</a-button>
        </a-col>
      </a-row>

      <a-table
          rowKey="ID"
          :columns="columns"
          :pagination="pagination"
          :dataSource="userlist"
          bordered
          @change="handleTableChange"
      >
        <span slot="role" slot-scope="data">{{data == 1 ? '管理员':'只读终端用户'}}</span>
        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button
                type="primary"
                icon="edit"
                style="margin-right:15px"
                @click="editUser(data.ID)"
            >编辑</a-button>
            <a-button
                type="danger"
                icon="delete"
                style="margin-right:15px"
                @click="deleteUser(data.ID)"
            >删除</a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 新增终端用户区域 -->
    <a-modal
        closable
        title="新增终端用户"
        :visible="addUserVisible"
        width="60%"
        @ok="addUserOk"
        @cancel="addUserCancel"
        destroyOnClose
    >
      <a-form-model :model="newUser" :rules="addUserRules" ref="addUserRef">
        <a-form-model-item label="用户名" prop="username">
          <a-input v-model="newUser.username"></a-input>
        </a-form-model-item>
        <a-form-model-item label="协议" prop="protocol">
          <a-select  style="width: 120px" @change="handleChange">
            <a-select-option value="ssh">
              ssh
            </a-select-option>
            <a-select-option value="vnc">
              vnc
            </a-select-option>
            <a-select-option value="rdp">
              rdp
            </a-select-option>
          </a-select>

        </a-form-model-item>
        <a-form-model-item label="私钥">
          <a-input v-model="newUser.identity_file"></a-input>
        </a-form-model-item>
        <a-form-model-item label="端口" prop="port">
          <a-input v-model="newUser.port"></a-input>
        </a-form-model-item>

        <a-form-model-item has-feedback label="密码" prop="password">
          <a-input-password v-model="newUser.password"></a-input-password>
        </a-form-model-item>
        <a-form-model-item has-feedback label="确认密码" prop="checkpass">
          <a-input-password v-model="newUser.checkpass"></a-input-password>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <!-- 编辑终端用户区域 -->
    <a-modal
        closable
        destroyOnClose
        title="编辑终端用户"
        :visible="editUserVisible"
        width="60%"
        @ok="editUserOk"
        @cancel="editUserCancel"
    >
      <a-form-model :model="userInfo" :rules="userRules" ref="addUserRef">
        <a-form-model-item label="终端用户名" prop="username">
          <a-input v-model="userInfo.username"></a-input>
        </a-form-model-item>
        <a-form-model-item has-feedback label="密码" prop="password">
          <a-input-password v-model="userInfo.password"></a-input-password>
        </a-form-model-item>
        <a-form-model-item label="协议" prop="protocol">
          <a-input v-model="userInfo.protocol"></a-input>
        </a-form-model-item>
        <a-form-model-item label="私钥">
          <a-input v-model="userInfo.identity_file"></a-input>
        </a-form-model-item>
        <a-form-model-item label="端口" prop="port">
          <a-input v-model="userInfo.port"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
const columns = [
  {
    title: 'ID',
    dataIndex: 'ID',
    width: '10%',
    key: 'id',
    align: 'center',
  },
  {
    title: '用户名',
    dataIndex: 'username',
    width: '20%',
    align: 'center',
  },
  {
    title: '密码',
    dataIndex: 'password',
    width: '20%',
    align: 'center',
  },
  {
    title: '协议',
    dataIndex: 'protocol',
    width: '20%',
  },
  {
    title: '端口',
    dataIndex: 'port',
    width: '20%',
    align: 'center',
  },
  {
    title: '私钥',
    dataIndex: 'identity_file',
    width: '20%',
    align: 'center',
  },
  {
    title: '操作',
    width: '30%',
    key: 'action',
    align: 'center',
    scopedSlots: { customRender: 'action' },
  },
]

export default {
  data() {
    return {
      pagination: {
        pageSizeOptions: ['5', '10', '20'],
        pageSize: 10,
        total: 0,
        showSizeChanger: true,
        showTotal: (total) => `共${total}条`,
      },
      userlist: [],
      userInfo: {
        username: '',
        password: '',
        checkPass: '',
        identity_file:'',
        protocol: '',
        port: '',
      },
      newUser: {
        username: '',
        password: '',
        identity_file:'',
        protocol: '',
        port: '',
        checkPass: '',

      },
      columns,
      queryParam: {
        username: '',
        pagesize: 5,
        pagenum: 1,
      },
      editVisible: false,
      userRules: {
        username: [
          {
            validator: (rule, value, callback) => {
              if (this.userInfo.username == '') {
                callback(new Error('请输入终端用户名'))
              }
              if ([...this.userInfo.username].length < 3 || [...this.userInfo.username].length > 12) {
                callback(new Error('终端用户名应当在3到12个字符之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        password: [
          {
            validator: (rule, value, callback) => {
              if (this.userInfo.password == '') {
                callback(new Error('请输入密码'))
              }
              if ([...this.userInfo.password].length < 6 || [...this.userInfo.password].length > 40) {
                callback(new Error('密码应当在6到20位之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        checkpass: [
          {
            validator: (rule, value, callback) => {
              if (this.userInfo.checkpass == '') {
                callback(new Error('请输入密码'))
              }
              if (this.userInfo.password !== this.userInfo.checkpass) {
                callback(new Error('密码不一致，请重新输入'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      addUserRules: {
        username: [
          {
            validator: (rule, value, callback) => {
              if (this.newUser.username == '') {
                callback(new Error('请输入终端用户名'))
              }
              if ([...this.newUser.username].length < 3 || [...this.newUser.username].length > 12) {
                callback(new Error('终端用户名应当在3到12个字符之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        password: [
          {
            validator: (rule, value, callback) => {
              if (this.newUser.password == '') {
                callback(new Error('请输入密码'))
              }
              if ([...this.newUser.password].length < 6 || [...this.newUser.password].length > 40) {
                callback(new Error('密码应当在6到40位之间'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
        checkpass: [
          {
            validator: (rule, value, callback) => {
              if (this.newUser.checkpass == '') {
                callback(new Error('请输入密码'))
              }
              if (this.newUser.password !== this.newUser.checkpass) {
                callback(new Error('密码不一致，请重新输入'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      editUserVisible: false,
      addUserVisible: false,
    }
  },
  created() {
    this.getUserList()
  },
  computed: {
    IsAdmin: function () {
      if (this.userInfo.role === 1) {
        return true
      } else {
        return false
      }
    },
  },
  methods: {
    // 获取终端用户列表
    async getUserList() {
      const { data: res } = await this.$http.get('term/getusers', {
        params: {
          username: this.queryParam.username,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.userlist = res.data
      this.pagination.total = res.total
    },
    handleChange(value) {
      console.log(value)
      this.newUser.protocol = value
    },
    // 更改分页
    handleTableChange(pagination) {
      var pager = { ...this.pagination }
      pager.current = pagination.current
      pager.pageSize = pagination.pageSize
      this.queryParam.pagesize = pagination.pageSize
      this.queryParam.pagenum = pagination.current

      if (pagination.pageSize !== this.pagination.pageSize) {
        this.queryParam.pagenum = 1
        pager.current = 1
      }
      this.pagination = pager
      this.getUserList()
    },
    // 删除终端用户
    deleteUser(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '确定要删除该终端用户吗？一旦删除，无法恢复',
        onOk: async () => {
          const { data: res } = await this.$http.delete(`term/deleteuser/${id}`)
          if (res.status != 200) return this.$message.error(res.message)
          this.$message.success('删除成功')
          this.getUserList()
        },
        onCancel: () => {
          this.$message.info('已取消删除')
        },
      })
    },
    // 新增终端用户
    addUserOk() {
      this.$refs.addUserRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const { data: res } = await this.$http.post('term/adduser', {
          username: this.newUser.username,
          password: this.newUser.password,
          protocol: this.newUser.protocol,
          identityFile: this.newUser.identity_file,
          port: Number(this.newUser.port),
        })
        if (res.status != 200) return this.$message.error(res.message)
        this.$refs.addUserRef.resetFields()
        this.addUserVisible = false
        this.$message.success('添加终端用户成功')
        this.getUserList()
      })
    },
    addUserCancel() {
      this.$refs.addUserRef.resetFields()
      this.addUserVisible = false
      this.$message.info('新增终端用户已取消')
    },
    // 编辑终端用户
    async editUser(id) {
      this.editUserVisible = true
      const { data: res } = await this.$http.get(`term/getuser/${id}`)
      this.userInfo = res.data
      this.userInfo.id = id
    },
    editUserOk() {
      this.$refs.addUserRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const { data: res } = await this.$http.put(`term/edituser/${this.userInfo.id}`, {
          username: this.userInfo.username,
          password: this.userInfo.password,
          protocol: this.userInfo.protocol,
          identityFile: this.userInfo.identity_file,
          port: Number(this.userInfo.port),
        })
        if (res.status != 200) return this.$message.error(res.message)
        this.editUserVisible = false
        this.$message.success('更新终端用户信息成功')
        this.getUserList()
      })
    },
    editUserCancel() {
      this.$refs.addUserRef.resetFields()
      this.editUserVisible = false
      this.$message.info('编辑已取消')
    },
  },
}
</script>

<style scoped>
.actionSlot {
  display: flex;
  justify-content: center;
}
</style>
