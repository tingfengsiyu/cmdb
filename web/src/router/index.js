import Vue from 'vue'
import VueRouter from 'vue-router'
const Login = () => import(/* webpackChunkName: "Login" */ '../views/Login.vue')
const Admin = () => import(/* webpackChunkName: "Admin" */ '../views/Admin.vue')

// 页面路由组件
const Index = () => import(/* webpackChunkName: "Index" */ '../components/admin/Index')
 const AddServer = () => import(/* webpackChunkName: "AddServer" */ '../components/idc/AddServer')
 const ListServer = () => import(/* webpackChunkName: "ListServer" */ '../components/idc/ListServer')
 const EditServer = () => import(/* webpackChunkName: "EditServer" */ '../components/idc/AddServer')
const UserList = () => import(/* webpackChunkName: "UserList" */ '../components/user/UserList')
const OpsRecords = () => import(/* webpackChunkName: "UserList" */ '../components/ops/OpsRecords')
const UpdateCluster = () => import(/* webpackChunkName: "UserList" */ '../components/ops/UpdateCluster')
const OsInit = () => import(/* webpackChunkName: "UserList" */ '../components/ops/OsInit')
const BatchIp = () => import(/* webpackChunkName: "UserList" */ '../components/ops/BatchIp')
const InstallMonitorAgent = () => import(/* webpackChunkName: "UserList" */ '../components/ops/InstallMonitorAgent')
const StorageMount = () => import(/* webpackChunkName: "UserList" */ '../components/ops/StorageMount')

// 路由重复点击捕获错误
const originalPush = VueRouter.prototype.push
VueRouter.prototype.push = function push(location, onResolve, onReject) {
  if (onResolve || onReject) return originalPush.call(this, location, onResolve, onReject)
  return originalPush.call(this, location).catch(err => err)
}

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'login',
    meta: {
      title: '请登录'
    },
    component: Login
  },
  {
    path: '/',
    name: 'admin',
    meta: {
      title: 'CMDB 后台管理页面'
    },
    component: Admin,
    children: [
      {
        path: 'index',
        component: Index,
        meta: {
          title: 'CMDB 后台管理页面'
        }
      },
      {
        path: 'addServer',
        component: AddServer,
        meta: {
          title: '新增服务器'
        }
      },
      {
        path: 'addServer/:id',
        component: EditServer,
        meta: {
          title: '编辑服务器'
        },
        props: true
      },
      {
        path: 'listServer',
        component: ListServer,
        meta: {
          title: '服务器列表'
        }
      },
      {
        path: 'opsRecords',
        component: OpsRecords,
        meta: {
          title: '操作记录列表'
        }
      },
      {
        path: 'updateCluster',
        component: UpdateCluster,
        meta: {
          title: '操作记录列表'
        }
      },
      {
        path: 'osInit',
        component: OsInit,
        meta: {
          title: '系统初始化'
        }
      },
      {
        path: 'batchIp',
        component: BatchIp,
        meta: {
          title: '修改ip'
        }
      },
      {
        path: 'installAgent',
        component: InstallMonitorAgent,
        meta: {
          title: '安装集群监控agent'
        }
      },
      {
        path: 'StorageMount',
        component: StorageMount,
        meta: {
          title: '挂载机器存储'
        }
      },
      {
        path: 'updateCluster',
        component: UpdateCluster,
        meta: {
          title: '更新机器所属集群'
        }
      },

      {
        path: 'userlist',
        component: UserList,
        meta: {
          title: '用户列表'
        }
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

router.beforeEach((to, from, next) => {
  if (to.meta.title) {
    document.title = to.meta.title
  }
  next()
  const token = window.sessionStorage.getItem('token')
  if (to.path === '/login') return next()
  if (!token) {
    next('/login')
  } else {
    next()
  }
})

export default router
