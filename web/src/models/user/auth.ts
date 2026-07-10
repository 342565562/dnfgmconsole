export type UserInfo = {
  id: number
  username: string
  time: string
  last_login_time: string | null
  email: string
  role: string
  desc: string
  game_uid?: number | null
  is_game_account?: boolean
}

export type EditUserInfo = {
  desc: string
}

export type Login = {
  username: string
  password: string
  activationCode?: string
}

export type LoginResult = {
  token: string
  username: string
  time: string
  last_login_time: string | null
  email: string
  role: string
  desc: string
  game_uid?: number | null
  is_game_account?: boolean
}
