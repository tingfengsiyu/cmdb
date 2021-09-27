<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
              v-model="queryParam.user"
              placeholder="输入用户查找"
              enter-button
              allowClear
              @search="getOpsRecords"
          />
        </a-col>
      </a-row>
    </a-card>
      <a-table
          rowKey="ID"
          :columns="columns"
          :pagination="pagination"
          :data-source="logs"
          @change="handleTableChange"
      >
        <p slot="expandedRowRender" slot-scope="ops"  style="margin: 0 ;white-space: pre-wrap" >
          {{ ops.log }}

        </p>
      </a-table>

  </div>
</template>

<script>
import moment from 'moment'

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
      logs: [],
      columns : [
        {
          title: 'ID',
          dataIndex: 'ID',
          width: '5%',
          align: 'left',
        },
        {
          title: '用户',
          dataIndex: 'user',
          width: '5%',
          align: 'left',
        },
        {
          title: '终端登录用户',
          dataIndex: 'term_user',
          align: 'left',
        },
        {
          title: '登录协议',
          dataIndex: 'protocol',
          align: 'left',
        },
        {
          title: '源ip',
          dataIndex: 'client_ip',
          align: 'left',
        },
        {
          title: '操作机器',
          dataIndex: 'private_ip_address',
          align: 'left',
        },
        {
          title: '日期',
          dataIndex: 'UpdatedAt',
          align: 'left',
          customRender: (val) => {
            return val ? moment(val).format('YYYY年MM月DD日 HH:mm') : '暂无'
          },
        },
      ],
      queryParam: {
        user: '',
        pagesize: 10,
        pagenum: 1,
      },
    }
  },
  created() {
    this.getOpsRecords()
  },
  methods: {
    // 获取操作记录
    async getOpsRecords() {
      const { data: res } = await this.$http.get('term/consolelog', {
        params: {
          user: this.queryParam.user,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.logs = res.data
      this.pagination.total = res.total
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
      this.getOpsRecords()
    },
  },
}
</script>

<style scoped>
.userSlot {
  display: flex;
  justify-content: left;
}
</style>
