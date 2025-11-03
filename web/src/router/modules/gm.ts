import { AppRouteRecordRaw } from '@/router/types'
import { LAYOUT } from '@/router/constant'

const accountsRoutes: Array<AppRouteRecordRaw> = [
  {
    path: '/gm',
    name: 'GM',
    component: LAYOUT,
    meta: {
      title: 'GM中心',
      icon: 'cog'
    },
    redirect: 'accounts',
    children: [

      {
        path: 'accounts',
        name: 'GmAccounts',
        component: () => import('@/views/gm/account/index.vue'),
        meta: {
          title: '点券充值',
          icon: 'list'
        }
      },
      {
        path: 'roles',
        name: 'GmRoles',
        component: () => import('@/views/gm/roles/index.vue'),
        meta: {
          title: '段位胜点',
          icon: 'fork'
        }
      },
      {
        path: 'tasks',
        name: 'GmTasks',
        component: () => import('@/views/gm/tasks/index.vue'),
        meta: {
          title: '删除邮件',
          icon: 'timer'
        }
      },
      {
        path: 'clear-creatures',
        name: 'GmClearCreatures',
        component: () => import('@/views/gm/tasks/ClearCreatures.vue'),
        meta: { title: '删除宠物', icon: 'people' }
      },
      {
        path: 'clear-avatars',
        name: 'GmClearAvatars',
        component: () => import('@/views/gm/tasks/ClearAvatars.vue'),
        meta: { title: '删除时装', icon: 'list-rich' }
      },
      {
        path: 'restore',
        name: 'GmRestore',
        component: () => import('@/views/gm/tasks/RestoreAccount.vue'),
        meta: { title: '一键恢复', icon: 'timer' }
      },
      {
        path: 'email',
        name: 'GmEmail',
        component: () => import('@/views/gm/email/index.vue'),
        meta: {
          title: '发送邮件',
          icon: 'transfer'
        }
      }/*,
      {
        path: 'test',
        name: 'UserTest',
        component: () => import('@/views/test.vue'),
        meta: {
          title: '测试',
          icon: 'list'
        }
      }
      */
    ]
  }
]

export default accountsRoutes
