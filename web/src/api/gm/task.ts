import { http } from '@/utils/http'
import { Task } from '@/views/gm/tasks/model'

export const getTasks = (characNo: number): Promise<Task[]> => {
  return http.request({
    url: `/gm/roles/${characNo}/tasks`,
    method: 'get'
  })
}

export const updateTasks = (characNo: number, data: any): Promise<any> => {
  return http.request({
    url: `/gm/roles/${characNo}/tasks`,
    method: 'put',
    data
  })
}

// 清空宠物栏（并清空邮件）
export const clearCreatures = (characNo: number): Promise<void> => {
  return http.request({
    url: `/gm/roles/${characNo}/clear_creatures`,
    method: 'post'
  })
}

// 清空时装栏（并清空邮件）
export const clearAvatars = (characNo: number): Promise<void> => {
  return http.request({
    url: `/gm/roles/${characNo}/clear_avatars`,
    method: 'post'
  })
}

// 一键恢复功能：同时执行删除邮件、删除宠物、删除时装
export const restoreAccount = (characNo: number): Promise<void> => {
  return http.request({
    url: `/gm/roles/${characNo}/restore`,
    method: 'post'
  })
}
