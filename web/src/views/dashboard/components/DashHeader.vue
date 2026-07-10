<template>
  <div class="dash-header">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="box-card" shadow="hover">
          <div>
            <span class="fs-14">活跃用户</span>
            <p class="stat-value">
              <span :style="fitFontSize(state.user_total)">{{ formatPrice(state.user_total) }}</span>
            </p>
            <div class="right-icon">
              <svg-icon icon-class="person"></svg-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card" shadow="hover">
          <div>
            <span class="fs-14">活跃角色</span>
            <p class="stat-value">
              <span :style="fitFontSize(state.charac_total)">{{ formatPrice(state.charac_total) }}</span>
            </p>
            <div class="right-icon">
              <svg-icon icon-class="nested"></svg-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card" shadow="hover">
          <div>
            <span class="fs-14">充值点券</span>
            <p class="stat-value">
              <span :style="fitFontSize(state.cera_total)">{{ formatPrice(state.cera_total) }}</span>
            </p>
            <div class="right-icon">
              <svg-icon icon-class="dollar" class-name="rela"></svg-icon>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card" shadow="hover">
          <div>
            <span class="fs-14">充值D点</span>
            <p class="stat-value">
              <span :style="fitFontSize(state.cera_point_total)">{{ formatPrice(state.cera_point_total) }}</span>
            </p>
            <div class="right-icon">
              <svg-icon icon-class="yen" class-name="rela"></svg-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { getDashTotal } from '@/api/dash'
import { formatPrice } from '@/utils'
import { reactive } from 'vue'

const state = reactive({
  user_total: 0,
  charac_total: 0,
  cera_point_total: 0,
  cera_total: 0
})

// 根据格式化后（含千分位）的字符串长度自适应字号，
// 数值越大字号越小，避免被右侧图标遮挡或溢出卡片。
const fitFontSize = (val: number) => {
  const len = formatPrice(val).length
  let size = 30
  if (len >= 16) size = 15
  else if (len >= 14) size = 17
  else if (len >= 12) size = 20
  else if (len >= 10) size = 23
  else if (len >= 8) size = 26
  return { fontSize: size + 'px' }
}

const getDashCountStat = async () => {
  const data = await getDashTotal()
  state.cera_total = data.cera_total
  state.cera_point_total = data.cera_point_total
  state.user_total = data.user_total
  state.charac_total = data.charac_total
}

getDashCountStat()
</script>

<style lang="scss" scoped>
.dash-header {
  .el-card {
    position: relative;
  }

  .stat-value {
    // 预留右侧图标区域（图标 right:25px + 宽 60px），大数字不会被遮住
    padding-right: 96px;
    min-height: 40px;
    margin: 8px 0 0;
    white-space: nowrap;

    span {
      font-weight: 600;
      line-height: 40px;
      display: inline-block;
      vertical-align: middle;
      transition: font-size 0.2s ease;
    }
  }

  .right-icon {
    position: absolute;
    top: 53%;
    right: 25px;
    width: 60px;
    height: 60px;
    line-height: 60px;
    color: var(--el-color-primary);
    text-align: center;
    background: var(--el-color-primary-light-9);
    border-radius: 50%;
    transform: translateY(-50%);

    .svg-icon {
      font-size: 28px;

      &.rela {
        position: relative;
        left: 3px;
        top: 4px;
      }
    }
  }
}
</style>
