import { storageLocal } from '@/utils/storage'
import { defineStore } from 'pinia'
import store from '@/store'

interface UserState {
  token?: string
  username?: string
  role?: string
  gameUid?: number
  isGameAccount?: boolean
}

export const useUserStore = defineStore({
  id: 'app-user',
  state: (): UserState => ({
    token: undefined,
    username: null,
    role: null,
    gameUid: undefined,
    isGameAccount: false
  }),
  getters: {
    getToken(): string {
      return this.token || storageLocal.getItem('token')
    },
    getUsername(): string {
      return this.username
    },
    getUserRole(): string {
      return this.role
    },
    getGameUid(): number | undefined {
      return this.gameUid || (storageLocal.getItem('gameUid') ? Number(storageLocal.getItem('gameUid')) : undefined)
    },
    getIsGameAccount(): boolean {
      return this.isGameAccount || storageLocal.getItem('isGameAccount') === 'true'
    }
  },
  actions: {
    setToken(token: string) {
      this.token = token
      storageLocal.setItem('token', token)
    },
    setUserInfo(username: string, role: string, gameUid?: number, isGameAccount?: boolean) {
      this.username = username
      this.role = role
      this.gameUid = gameUid
      this.isGameAccount = isGameAccount || false
      if (gameUid !== undefined) {
        storageLocal.setItem('gameUid', String(gameUid))
      }
      storageLocal.setItem('isGameAccount', String(this.isGameAccount))
    },
    removeUserStore() {
      this.token = ''
      this.username = ''
      this.role = ''
      this.gameUid = undefined
      this.isGameAccount = false
      storageLocal.clear()
    }
  }
})

// Need to be used outside the setup
export function useUserStoreWithOut() {
  return useUserStore(store)
}
