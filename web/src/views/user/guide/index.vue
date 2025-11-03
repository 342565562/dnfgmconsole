<template>
  <div>
    <panel-title title="操作介绍"></panel-title>

    <el-row>
      <el-col :span="18">
        <div class="guide-content">
          <div class="guide-section">
            <h3>后台刷一切</h3>
            <p>可以刷道具、装备、武器、材料、宠物、任务物品</p>
          </div>

          <div class="guide-section">
            <h3>邮件自动清空</h3>
            <p>全服角色邮件每隔30分钟会自动清空一次，避开半点和整点刷道具</p>
          </div>

          <div class="guide-section">
            <h3>网络中断问题处理</h3>
            <ul class="guide-list">
              <li>如果遇到点击角色网络中断，刷普通道具要在后台执行<strong>删除邮件</strong></li>
              <li>刷了宠物之后网络中断就去后台执行<strong>删除宠物</strong></li>
              <li>如果是刷了散件时装有问题，就去后台执行<strong>删除时装</strong></li>
            </ul>
          </div>

          <div class="guide-section">
            <h3>举例说明</h3>
            <div class="guide-example">
              <p><strong>天空套：</strong>搜索关键字"稀有"，发送天空套盒子</p>
            </div>
          </div>

          <div class="guide-section">
            <h3>转职和觉醒</h3>
            <p>转职和觉醒无需做任务可直接发送相关道具，需要在后台物品代码中搜索关键字"转职"和"觉醒"</p>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import PanelTitle from '@/components/PanelTitle'
</script>

<style scoped>
.guide-content {
  padding: 20px;
  line-height: 1.8;
}

.guide-section {
  margin-bottom: 30px;
}

.guide-section h3 {
  color: #303133;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.guide-section p {
  color: #606266;
  font-size: 14px;
  margin: 8px 0;
}

.guide-list {
  color: #606266;
  font-size: 14px;
  padding-left: 20px;
  margin: 12px 0;
}

.guide-list li {
  margin: 10px 0;
  line-height: 1.6;
}

.guide-list strong {
  color: #409eff;
  font-weight: 600;
}

.guide-example {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  margin-top: 10px;
}

.guide-example p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.guide-example strong {
  color: #303133;
  font-weight: 600;
}
</style>

