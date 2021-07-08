<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="6">
          <a-input-search
              v-model="queryParam.action"
              placeholder="输入操作分类查找"
              enter-button
              allowClear
              @search="getOpsRecords"
          />
        </a-col>
      </a-row>
      <a-table
          rowKey="ID"
          :columns="columns"
          :pagination="pagination"
          :data-source="recordlist"

          @change="handleTableChange"
      >
<!--        <span slot="state" slot-scope="state">{{ state == 1 ? '成功': '失败' }}</span>-->
        <span slot="state" slot-scope="state">
           <div v-if="state===1">
             成功
    </div>
           <div v-else-if="state===0">
         失败
     </div>
          <div v-else>
         执行中
     </div>
        </span>
      </a-table>
    </a-card>

  </div>
</template>

<script>
import moment from 'moment'
const columns = [
  {
    title: 'ID',
    dataIndex: 'ID',
    width: '3%',
    align: 'center',
  },
  {
    title: '操作分类',
    dataIndex: 'action',
    width: '7%',
    align: 'center',
  },
  {
    title: '操作目标',
    dataIndex: 'object',
    width: '20%',
    align: 'center',
  },
  {
    title: '操作状态',
    dataIndex: 'state',
    width: '5%',
    align: 'center',
    scopedSlots: { customRender: 'state' },
  },
  {
    title: '成功记录',
    dataIndex: 'success',
    width: '20%',
    align: 'center',
  },
  {
    title: '失败记录',
    dataIndex: 'error',
    width: '20%',
    align: 'center',
  },
  {
    title: '创建日期',
    dataIndex: 'CreatedAt',
    width: '7%',
    align: 'center',
    customRender: (val) => {
      return val ? moment(val).format('YYYY年MM月DD日 HH:mm') : '暂无'
    },
  },
  {
    title: '更新日期',
    dataIndex: 'UpdatedAt',
    width: '7%',
    align: 'center',
    customRender: (val) => {
      return val ? moment(val).format('YYYY年MM月DD日 HH:mm') : '暂无'
    },
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
      recordlist: [],
      columns,
      queryParam: {
        action: '',
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
      const { data: res } = await this.$http.get('idc/opsrecords', {
        params: {
          action: this.queryParam.action,
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      })
      if (res.status != 200) return this.$message.error(res.message)
      this.recordlist = res.data
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
.actionSlot {
  display: flex;
  justify-content: center;
}
</style>
