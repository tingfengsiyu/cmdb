<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
              v-model="queryParam.group"
              placeholder="输入组名查找"
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
          rowKey="id"
          :columns="columns"
          :pagination="pagination"
          :dataSource="userlist"
          bordered
          @change="handleTableChange"
      >
        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button
                type="primary"
                icon="edit"
                style="margin-right:15px"
                @click="editPermission(data.ID)"
            >编辑</a-button>
            <a-button
                type="danger"
                icon="delete"
                style="margin-right:15px"
                @click="deleteUserPermissions(data.ID)"
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
        @ok="addPermissionsOk"
        @cancel="addUserCancel"
        destroyOnClose
    >
      <a-form-model :model="newPermission" :role="newPermissionRule" ref="addUserRef" >
        <a-form-model-item label="组名" prop="group">
          <a-input v-model="newPermission.group"></a-input>
        </a-form-model-item>
        <a-form-model-item label="系统用户" v-model="newPermission.UserLists" prop="users">
          <a-select mode="multiple"  placeholder="Please select"  style="width: 200px" @change="UserChange">
            <a-select-option v-for="item in UserList" :key="item.id" :value="item.username" >{{item.username}}</a-select-option>
          </a-select>
        </a-form-model-item>
        <a-form-model-item  label="ip池" v-model="newPermission.Ips"  prop="ips">
          <a-select mode="multiple"  placeholder="Please select"  style="width: 200px" @change="IpChange">
            <a-select-option v-for="item in IpList" :key="item.id" :value="item.private_ip_address" >{{item.private_ip_address}}</a-select-option>
          </a-select>
        </a-form-model-item>
        <a-form-model-item label="终端用户" v-model="newPermission.TermUsers"  prop="usernames">
          <a-select  placeholder="Please select"  style="width: 200px" @change="TermUserChange">
            <a-select-option v-for="item in TermUserList" :key="item.id" :value="item.username" >{{item.username}}</a-select-option>
          </a-select>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <!-- 编辑终端用户区域 -->
    <a-modal
        closable
        destroyOnClose
        title="编辑终端用户权限"
        :visible="editPermissionVisible"
        width="60%"
        @ok="editPermissionOk"
        @cancel="editPermissionCancel"
    >
      <a-form-model :model="userInfo" :rules="addPermissionsRules" ref="addUserRef">
        <a-form-model-item label="组名" prop="username">
          <a-input v-model="userInfo.username"></a-input>
        </a-form-model-item>
        <a-form-model-item label="系统用户名" prop="username">
          <a-input v-model="userInfo.username"></a-input>
        </a-form-model-item>
        <a-form-model-item has-feedback label="授权ip" prop="password">
          <a-input-password v-model="userInfo.password"></a-input-password>
        </a-form-model-item>
        <a-form-model-item label="终端用户" prop="port">
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
    dataIndex: 'id',
    width: '5%',
    key: 'id',
    align: 'center',
  },
  {
    title: '组名',
    dataIndex: 'group',
    align: 'center',
  },
  {
    title: '系统用户',
    dataIndex: 'user',
    align: 'center',
  },
  {
    title: '主机名',
    dataIndex: 'name',
    align: 'center',
  },
  {
    title: 'ip',
    dataIndex: 'private_ip_address',
    width: '10%',
    align: 'center',
  },
  {
    title: '终端用户',
    dataIndex: 'username',
    align: 'center',
  },
  {
    title: '操作',
    width: '15%',
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
      newPermissionRule: {
        usernames: '',
        ips: '',
        users: '',
      },
      columns,
      queryParam: {
        group: '',
        pagesize: 5,
        pagenum: 1,
      },
      IpList: [],
      TermUserList: [],
      UserList: [],
      newPermission: {
        Ips: [],
        TermUsers: [],
        UserLists: [],
        group: '',
      },
      editVisible: false,
      addPermissionsRules: {
        usernames: [{ required: true, message: '选择机器登录用户', trigger: 'blur' }],
        users: [{ required: true, message: '请选择系统用户 ，可多选', trigger: 'blur' }],
        ips: [{ required: true, message: '请选择多个ip，all所有', trigger: 'blur' }],
      },
      editPermissionVisible: false,
      addUserVisible: false,
    }
  },
  created() {
    this.getUserList()
    this.getSrvList()
    this.getTermUserList()
    this.getPermissionsList()
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
    async getPermissionsList() {
      const { data: res } = await this.$http.get('term/permissions', {
        params: {
          group: this.newPermission.group,
          username: this.queryParam.username,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.userlist = res.data
      this.pagination.total = res.total
    },
    async getSrvList() {
      const { data: res } = await this.$http.get('idc/getservers', {
        params: {
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.IpList = [
        { id: 0, private_ip_address: 'all' },
        ...(res.data || [])
      ]
    },
    async getTermUserList() {
      const { data: res } = await this.$http.get('term/getusers', {
        params: {
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.TermUserList = res.data
      this.pagination.total = res.total
    },
    async getUserList() {
      const { data: res } = await this.$http.get('user/getusers', {
        params: {
          group: this.queryParam.group,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.UserList = res.data
    },
    IpChange(value) {
      this.newPermission.Ips =value
    },
    UserChange(value) {
      this.newPermission.UserLists = value
    },
    TermUserChange(value) {
      this.newPermission.TermUsers =value
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
    deleteUserPermissions(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '确定要删除该终端用户吗？一旦删除，无法恢复',
        onOk: async () => {
          const { data: res } = await this.$http.delete(`term/deletePermission/${id}`)
          if (res.status != 200) return this.$message.error(res.message)
          this.$message.success('删除成功')
          this.getUserList()
        },
        onCancel: () => {
          this.$message.info('已取消删除')
        },
      })
    },
    // 新增终端用户权限
    addPermissionsOk() {
      this.$refs.addUserRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        console.log(this.newPermission)
        const { data: res } = await this.$http.post('term/addpermissions', {
          group: this.newPermission.group,
          ips: this.newPermission.Ips,
          term_users:  this.newPermission.TermUsers,
          users: this.newPermission.UserLists,
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
    async editPermission(id) {
      this.editPermissionVisible = true
      const { data: res } = await this.$http.get(`term/getuser/${id}`)
      this.userInfo = res.data
      this.userInfo.id = id
    },
    editPermissionOk() {
      this.$refs.addUserRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const { data: res } = await this.$http.put(`term/editPermission/${this.userInfo.id}`, {
          username: this.userInfo.username,
          password: this.userInfo.password,
          protocol: this.userInfo.protocol,
        })
        if (res.status != 200) return this.$message.error(res.message)
        this.editPermissionVisible = false
        this.$message.success('更新终端用户信息成功')
        this.getUserList()
      })
    },
    editPermissionCancel() {
      this.$refs.addUserRef.resetFields()
      this.editPermissionVisible = false
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
