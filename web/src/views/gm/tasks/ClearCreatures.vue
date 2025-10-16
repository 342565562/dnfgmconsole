<template>
  <div>
    <panel-title title="删除宠物"></panel-title>

    <div class="tc-step-border">
      <div class="l">
        <panel-step num="1" title="账号选择">
          <select-account @setUid="setUid"></select-account>
        </panel-step>

        <panel-step num="2" title="角色选择">
          <role-table
            ref="roleTableRef"
            :loading="state.loading"
            :data="state.data"
            @setCharacNo="setCharacNo"
          ></role-table>
        </panel-step>

        <panel-step num="3" title="执行删除" is-last>
          <el-button type="danger" :disabled="!characNo" @click="doClear" size="small">清空宠物背包（切换角色后生效）</el-button>
        </panel-step>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import PanelTitle from '@/components/PanelTitle'
import PanelStep from '@/components/PanelStep'
import SelectAccount from './SelectAccount'
import RoleTable from './RoleTable'
import { reactive, ref } from 'vue'
import { RoleState } from '@/views/gm/roles/model'
import { getRoles } from '@/api/gm/role'
import { clearCreatures } from '@/api/gm/task'
import { confirmWarning } from '@/utils/element/messageBox'
import { successMessage } from '@/utils/element/message'

const state = reactive<RoleState>({ data: [], loading: false })
const roleTableRef = ref<any>(null)
const characNo = ref<number | null>(null)

const setUid = async (uid: number) => {
  state.data = []
  roleTableRef.value.resetCharacNo()
  try {
    state.loading = true
    state.data = await getRoles(uid)
    state.loading = false
  } catch (e) {
    state.loading = false
  }
}

const setCharacNo = (id: number) => {
  characNo.value = id
}

const doClear = async () => {
  if (!characNo.value) return
  await confirmWarning('确认清空该角色宠物背包？')
  await clearCreatures(characNo.value)
  successMessage('执行成功，已清空！切换角色后生效')
}
</script>

<style scoped></style>


