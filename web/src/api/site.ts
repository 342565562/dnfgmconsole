import { http } from '@/utils/http'

export interface SiteConfig {
  title: string
  login_name: string
}

// 站点信息(标题/登录名)，免鉴权公开接口
export const getSiteConfig = (): Promise<SiteConfig> => {
  return http.request({
    url: '/site-config',
    method: 'get'
  })
}
