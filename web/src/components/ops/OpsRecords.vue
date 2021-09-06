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
    </a-card>
      <a-table
          rowKey="ID"
          :columns="columns"
          :pagination="pagination"
          :data-source="recordlist"
          @change="handleTableChange"
      >
        <span slot="state" slot-scope="state">
           <div v-if="state===1">成功</div>
           <div v-else-if="state===0">失败</div>
           <div v-else>执行中</div>
        </span>

        <p slot="expandedRowRender" slot-scope="ops"  style="margin: 0 ;white-space: pre-wrap" >
          {{ ops.success }}

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
      recordlist: [],
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
          title: '操作分类',
          dataIndex: 'action',
          width: '10%',
          align: 'left',
        },
        {
          title: '操作目标',
          dataIndex: 'object',
          width: '30%',
          align: 'left',
          customRender: (object) => {
            const textArr = object.split('\n')
            return (<div>
              {
                textArr.map(t => {
                  return (<li>{t}</li>)
                })
              }
            </div>)
          }
        },
        {
          title: '操作状态',
          dataIndex: 'state',
          width: '15%',
          align: 'left',
          scopedSlots: { customRender: 'state' },
        },
        // {
        //   title:'成功记录',
        //   align:'left',
        //   dataIndex: 'success',
        //   customRender: (success) => {
        //     const textArr = success.split('\n')
        //     return (<div>
        //       {
        //         textArr.map(t => {
        //           return (<li>{t}</li>)
        //         })
        //       }
        //     </div>)
        //   }
        // },
        {
          title: '失败记录',
          dataIndex: 'error',
          width: '15%',
          align: 'left',
          customRender: (error) => {
            const textArr = error.split('\n')
            return (<div>
              {
                textArr.map(t => {
                  return (<li>{t}</li>)
                })
              }
            </div>)
          }
        },
        {
          title: '日期',
          dataIndex: 'CreatedAt',
          width: '15%',
          align: 'left',
          customRender: (val) => {
            return val ? moment(val).format('YYYY年MM月DD日 HH:mm') : '暂无'
          },
        },
      ],
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
  justify-content: left;
}
</style>
