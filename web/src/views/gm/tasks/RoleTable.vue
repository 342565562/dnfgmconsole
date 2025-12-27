<template>
  <div>
    <el-table v-loading="loading" :data="data" ref="tableRef" border>
      <el-table-column width="80" label="请勾选" align="center" class-name="select-column">
        <template #default="scope">
          <el-radio v-model="characNo" :label="scope.row.charac_no" class="radio" @change="selectCharacNo"></el-radio>
        </template>
      </el-table-column>
      <el-table-column prop="charac_no" label="角色ID" width="180" />
      <el-table-column prop="charac_name" label="角色名称" width="180"></el-table-column>
      <el-table-column prop="create_time" label="创建时间" width="180">
        <template #default="scope">
          <span>{{ dateFormat(scope.row.create_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="lev" label="等级"></el-table-column>
      <el-table-column prop="m_id" label="uid" width="150" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { defineExpose, reactive, ref } from 'vue'
import { RoleState, Role } from '@/views/gm/roles/model'
import { dateFormat } from '@/utils'

const characNo = ref<number>(null)

defineProps({
  data: {
    type: Array as PropType<Role[]>,
    required: true,
    default: () => []
  },
  loading: {
    type: Boolean,
    required: true,
    default: false
  }
})

const selectCharacNo = () => {
  emit('setCharacNo', characNo.value)
}

const resetCharacNo = () => {
  characNo.value = null
}

// hook
const emit = defineEmits(['setCharacNo'])

defineExpose({
  resetCharacNo
})
</script>

<style scoped>
:deep(.select-column .cell) {
  white-space: nowrap;
  min-width: 80px;
  width: 80px;
}

/* 请勾选表头文字设置为红色 */
:deep(.select-column.is-leaf .cell) {
  color: #f56c6c;
}

/* 单选按钮圆圈线条颜色改为绿色 */
:deep(.select-column .el-radio__inner) {
  border-color: #67c23a;
}

:deep(.select-column .el-radio__input:hover .el-radio__inner) {
  border-color: #67c23a;
}

:deep(.select-column .el-radio__input.is-checked .el-radio__inner) {
  border-color: #67c23a;
  background-color: #67c23a;
}

:deep(.select-column .el-radio__input.is-checked .el-radio__inner::after) {
  background-color: #fff;
}
</style>
