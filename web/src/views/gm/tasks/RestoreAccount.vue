<template>
  <div>
    <el-row>
      <div class="console-title clearfix" style="min-height: 50px">
        <div class="pull-left">
          <h4>一键恢复</h4>
          <span style="color: #f56c6c; margin-left: 8px;">解决无法登录游戏和网络中断等问题</span>
        </div>
      </div>
    </el-row>

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

        <panel-step num="3" title="执行恢复" is-last>
          <div>
            <p style="color: #f56c6c; margin-bottom: 12px; font-weight: 500;">
              重新登录生效
            </p>
            <p style="color: #909399; margin-bottom: 12px;">
              解决无法登录游戏和网络中断问题
            </p>
            <el-button type="danger" :disabled="!characNo" @click="doRestore" size="small">
              一键恢复
            </el-button>
          </div>
        </panel-step>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import PanelStep from '@/components/PanelStep'
import SelectAccount from './SelectAccount'
import RoleTable from './RoleTable'
import { reactive, ref } from 'vue'
import { RoleState } from '@/views/gm/roles/model'
import { getRoles } from '@/api/gm/role'
import { restoreAccount } from '@/api/gm/task'
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

const doRestore = async () => {
  if (!characNo.value) return
  await confirmWarning('确认执行一键恢复？将同时删除邮件、宠物和时装（切换角色后生效）')
  await restoreAccount(characNo.value)
  successMessage('执行成功，已恢复！切换角色后生效')
}
</script>

<style scoped>
.console-title {
  padding: 16px 0;
  min-height: 48px;
}

.pull-left {
  float: left;
}

.console-title h4 {
  display: inline-block;
  text-indent: 8px;
  border-left: 2px solid #39f;
  margin-top: 0;
  margin-bottom: 0;
  margin-right: 8px;
  vertical-align: top;
  font-weight: 700;
  color: #555;
  font-size: 18px;
}
</style>

