<template>
  <div v-if="!isGameAccount">
    <el-form :inline="true" :model="form" ref="formRef" label-width="100px" @submit.native.prevent>
      <div>
        <el-form-item label="玩家账号:" prop="account_name">
          <el-input v-model="form.account_name" clearable placeholder="请输入玩家账号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="filterAccount">查询</el-button>
          <el-button @click.prevent="reset">重置</el-button>
        </el-form-item>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { FormInstance } from 'element-plus'
import { FilterAccountForm } from '@/views/gm/account/model/model'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const isGameAccount = computed(() => userStore.getIsGameAccount)

const formRef = ref<FormInstance>()
const form = reactive<FilterAccountForm>({
  uid: '',
  account_name: ''  // 新增字段
})

// method
const filterAccount = () => {
  emit('setParams', form)
}

const reset = () => {
  formRef.value.resetFields()
  emit('setParams', {})
}

// hook
const emit = defineEmits(['setParams'])
</script>

<style scoped></style>
