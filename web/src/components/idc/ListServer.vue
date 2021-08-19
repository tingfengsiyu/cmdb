<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
              v-model="queryParam.private_ip_address"
              placeholder="输入私有ip查找"
              enter-button
              allowClear
              @search="getServerList"
          />
        </a-col>
        <a-col :span="3">
          <a-button type="primary" @click="$router.push('/addserver')">新增</a-button>
        </a-col>
        <a-col :span="3">
          <a-button type="primary" @click="batchaddServer = true">批量新增</a-button>
          <!-- 批量新增区域 -->
          <a-modal
              closable
              title="批量新增服务器"
              :visible="batchaddServer"
              width="60%"
              @cancel="addServersCancel"
              destroyOnClose
          >
            <a-upload name="file"  :headers="headers" :action=upUrl accept=".xlsx" method="post" @change="upChange" >
              <a-button> <a-icon type="upload" /> Click to Upload </a-button>
            </a-upload>
          </a-modal>

        </a-col>
        <a-col :span="4">
        <a-button type="primary" :headers="headers" @click="downloadAllServer">导出服务器</a-button>
        </a-col>
        <a-col :span="3">
          <a-select placeholder="请选择集群" style="width:250px" @change="ClusterChange">
            <a-select-option v-for="item in ClusterList" :key="item.id" :value="item.cluster">{{item.cluster}}</a-select-option>
          </a-select>
        </a-col>
        <a-col :span="2">
          <a-button type="info" @click="getAllServerList()">显示全部</a-button>
        </a-col>
      </a-row>

      <a-table
          rowKey="id"
          :columns="columns"
          :pagination="pagination"
          :dataSource="Artlist"
          bordered
          @change="handleTableChange"
      >
        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button
                size="small"
                type="primary"
                icon="edit"
                style="margin-right:15px"
                @click="$router.push(`/addserver/${data.id}`)"
            >编辑</a-button>
            <a-button
                size="small"
                type="danger"
                icon="delete"
                style="margin-right:15px"
                @click="deleteArt(data.id)"
            >删除</a-button>
          </div>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script>
import { Url } from '../../plugin/http'
import FileDownload from "js-file-download"
import axios from 'axios'
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: '5%',
    key: 'id',
    align: 'center',
  },
  {
    title: '主机名',
    dataIndex: 'name',
    width: '10%',
    key: 'name',
    align: 'center',
  },
  {
    title: '私有ip',
    dataIndex: 'private_ip_address',
    width: '10%',
    key: '私有ip',
    align: 'center',
  },
  {
    title: '标签ip',
    dataIndex: 'label_ip_address',
    width: '10%',
    key: '标签ip',
    align: 'center',
  },
  {
    title: '所属集群',
    dataIndex: 'cluster',
    width: '10%',
    //key: 'desc',
    align: 'center',
  },
  {
    title: '所属用户',
    dataIndex: 'user',
    width: '10%',
    //key: 'desc',
    align: 'center',
  },
  {
    title: '机房名',
    dataIndex: 'idc_name',
    width: '5%',
    //key: 'desc',
    align: 'center',
  },
  {
    title: '机柜',
    dataIndex: 'cabinet_number',
    width: '5%',
   // key: 'desc',
    align: 'center',
  },
  {
    title: '上架状态',
    dataIndex: 'state',
    width: '5%',
    //key: 'desc',
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
      Artlist: [],
      upUrl: Url + 'idc/uploadexcel',
      downUrl: Url+ 'idc/exportcsv',
      ClusterList: [],
      columns,
      queryParam: {
        cluster:'',
        private_ip_address: '',
        pagesize: 10,
        pagenum: 1,
      },
      batchaddServer: false,
      headers: {
       Authorization: `Bearer ${localStorage.getItem('token')}`
      },
    }
  },
  created() {
    this.getServerList()
    this.getClusterList()
  },
  methods: {
    // 获取服务器列表
    async getServerList() {
      const { data: res } = await this.$http.get('idc/getnetworktopology', {
        params: {
          private_ip_address: this.queryParam.private_ip_address,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
          cluster: this.queryParam.cluster,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)

      this.Artlist = res.data
      this.pagination.total = res.total
    },
    async getAllServerList() {
      const { data: res } = await this.$http.get('idc/getnetworktopology', {
        params: {
          private_ip_address: this.queryParam.private_ip_address,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)

      this.Artlist = res.data
      this.pagination.total = res.total
    },
    // 获取机器集群
    async getClusterList() {
      const { data: res } = await this.$http.get('idc/getclusters')
      if (res.status !== 200) return this.$message.error(res.message)
      this.ClusterList = res.data
      this.pagination.total = res.total
      // console.log(res.total)
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
      this.getServerList()
    },
    // 删除服务器
    deleteArt(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '确定要删除该服务器吗？一旦删除，无法恢复',
        onOk: async () => {
          const { data: res } = await this.$http.delete(`idc/deleteserver/${id}`)
          if (res.status != 200) return this.$message.error(res.message)
          this.$message.success('删除成功')
          this.getServerList()
        },
        onCancel: () => {
          this.$message.info('已取消删除')
        },
      })
    },
    // 查询集群下的服务器
    ClusterChange(value) {
      this.queryParam.cluster=value
      this.getClusterServer(value)
    },
    async getClusterServer(value) {
      const { data: res } = await this.$http.get(`idc/getnetworktopology`, {
        params: {
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
          cluster: value
        },
      })
      if (res.status !== 200) return this.$message.error(res.message)
      this.Artlist = res.data
      this.pagination.total = res.total
    },
    // 取消批量导入服务器
    addServersCancel() {
      this.batchaddServer = false
      this.$message.info('批量新增服务器已取消')
    },

    downloadAllServer(){
     axios({
        method: 'get',
        url: this.downUrl,
        headers: {
          'Authorization': this.headers
        },
        responseType: 'blob'
      }).then(res => {
        FileDownload(res.data, '服务器.csv');
      })
    },

    upChange(info) {
      if (info.file.response.status === 200 ) {
        this.$message.success(`${info.file.name} ${info.file.response.message}`);
      } else {
        this.$message.error(`${info.file.name} ${info.file.response.message}`);
      }
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
