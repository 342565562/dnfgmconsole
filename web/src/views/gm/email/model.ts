export type Email = {
  code: number // 物品代码
  number: number // 物品数量
  seperate_upgrade: number // 武器锻造等级
  upgrade: number // 装备强化等级
  is_amplify: boolean // 具备异界属性
  amplify_option: number // 异界类型
  amplify_value: number // 追加属性
  gold: number // 金币
  seal_flag: boolean // 是否封装
  avata_flag?: boolean // 时装邮件标记
  creature_flag?: boolean // 宠物邮件标记
}

export type FilterGold = {
  name: string
}
