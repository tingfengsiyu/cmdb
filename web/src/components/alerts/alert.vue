<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
              v-model="queryParam.username"
              placeholder="输入告警项查找"
              enter-button
              allowClear
              @search="getalertlist"
          />
        </a-col>
      </a-row>
      <a-table
          rowKey="id"
          :columns="columns"
          :pagination="pagination"
          :dataSource="alertlist"
          bordered
          @change="handleTableChange"
      >
      </a-table>
    </a-card>
  </div>
</template>

<script>
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: '4%',
    key: 'id',
    align: 'left',
  },
  {
    title: '告警名',
    dataIndex: 'alertname',
    align: 'left',
  },
  {
    title: '集群',
    dataIndex: 'cluster',
    width: '10%',
    align: 'left',
  },
  {
    title: '告警机器',
    dataIndex: 'instance',
    align: 'left',
  },
  {
    title: '告警值',
    dataIndex: 'value',
    align: 'left',
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
      alertlist: [],
      columns,
      queryParam: {
        username: '',
        pagesize: 5,
        pagenum: 1,
      },
    }
  },
  created() {
    this.getalertlist()
  },
  computed: {
  },
  methods: {
    // 获取终端用户列表
    async getalertlist() {
      const { data: res } = await this.$http.get('idc/prometheusalerts', {
        params: {
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.alertlist = res.data
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
      this.getalertlist()
    },
  },
}
</script>

<style scoped>
th.column-money,
td.column-money {
  text-align: right !important;
}
</style>
