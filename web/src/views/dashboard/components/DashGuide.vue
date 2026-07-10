<template>
  <div class="dash-guide">
    <el-card class="box-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <svg-icon icon-class="list-rich"></svg-icon>
          <span class="card-title">操作介绍</span>
        </div>
      </template>

      <!-- 顶部图文链接(整行可点，超链样式) -->
      <a
        class="guide-banner"
        href="https://docs.qq.com/doc/DTHJzVUZVbGZ4Z1R3"
        target="_blank"
        rel="noopener noreferrer"
      >
        <span class="banner-icon">📘</span>
        <span class="banner-text">图文操作介绍，请点<span class="banner-link">此处</span>（新手必看）</span>
      </a>

      <!-- 操作说明卡片网格 -->
      <div class="guide-grid">
        <div v-for="item in guides" :key="item.no" class="guide-item">
          <div class="guide-no">{{ item.no }}</div>
          <div class="guide-detail">
            <h4 class="guide-title">{{ item.title }}</h4>
            <p class="guide-desc" v-html="item.desc"></p>
          </div>
        </div>
      </div>

      <!-- 举例说明 -->
      <div class="guide-callout example">
        <div class="callout-label">举例说明</div>
        <p>
          <strong>天空套：</strong>搜索关键字“稀有”，出现天空套盒子代码。搜索关键字“强化券”，出现各种强化券代码，然后通过邮件发送。
        </p>
      </div>

      <!-- 网络中断处理 -->
      <div class="guide-callout warn">
        <div class="callout-label warn-label">网络中断游戏进不去问题处理（执行后需切换角色生效）</div>
        <ul class="warn-list">
          <li>点击角色网络中断、刷普通道具：后台执行 <strong>删除邮件</strong></li>
          <li>刷了宠物之后网络中断：后台执行 <strong>删除宠物</strong></li>
          <li>刷了散件时装有问题：后台执行 <strong>删除时装</strong></li>
          <li>以上均未解决：后台执行 <strong>一键恢复</strong></li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
const guides = [
  {
    no: 1,
    title: '后台刷一切',
    desc: '可充值点券、胜点，刷道具/装备/武器/材料/宠物/任务物品。<b>充值胜点后需切换角色或重新登录才生效！</b>'
  },
  {
    no: 2,
    title: '金币和道具发送',
    desc: '金币与道具发送绑定在一起，发送金币需带道具代码一起发送，退到角色选择界面重新进入即可领取。'
  },
  {
    no: 3,
    title: '转职和觉醒',
    desc: '无需做任务可直接发送相关道具，在后台物品代码中搜索关键字 <b>“转职”</b> 和 <b>“觉醒”</b>。'
  },
  {
    no: 4,
    title: '邮件自动清空',
    desc: '全服角色邮件每隔 30 分钟自动清空一次，<b>请避开半点和整点刷道具。</b>'
  }
]
</script>

<style lang="scss" scoped>
.dash-guide {
  margin-top: 18px;

  .card-header {
    display: flex;
    align-items: center;
    font-size: 15px;
    font-weight: 600;

    .card-title {
      margin-left: 10px;
    }
  }

  // 链接横幅(整行可点，超链样式)
  .guide-banner {
    display: flex;
    align-items: center;
    padding: 12px 18px;
    margin-bottom: 20px;
    border-radius: 8px;
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
    font-size: 15px;
    font-weight: 500;
    text-decoration: none;
    transition: all 0.25s ease;

    .banner-icon {
      font-size: 20px;
      margin-right: 10px;
    }

    .banner-text {
      flex: 1;
    }

    // “此处”呈现超链接样式
    .banner-link {
      color: var(--el-color-primary);
      text-decoration: underline;
      font-weight: 700;
      margin: 0 2px;
    }

    &:hover {
      background: var(--el-color-primary-light-8);

      .banner-text {
        text-decoration: underline;
      }
    }
  }

  // 网格
  .guide-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    margin-bottom: 20px;
  }

  @media (max-width: 900px) {
    .guide-grid {
      grid-template-columns: 1fr;
    }
  }

  .guide-item {
    display: flex;
    padding: 16px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    background: #fafcff;
    transition: all 0.25s ease;

    &:hover {
      border-color: var(--el-color-primary-light-5);
      box-shadow: 0 4px 14px rgba(0, 0, 0, 0.06);
      transform: translateY(-2px);
    }

    .guide-no {
      flex-shrink: 0;
      width: 30px;
      height: 30px;
      line-height: 30px;
      text-align: center;
      border-radius: 50%;
      background: var(--el-color-primary);
      color: #fff;
      font-weight: 700;
      font-size: 15px;
      margin-right: 14px;
    }

    .guide-detail {
      flex: 1;
    }

    .guide-title {
      margin: 2px 0 8px;
      font-size: 15px;
      font-weight: 600;
      color: #303133;
    }

    .guide-desc {
      margin: 0;
      font-size: 13px;
      line-height: 1.7;
      color: #606266;

      :deep(b) {
        color: var(--el-color-danger);
        font-weight: 600;
      }
    }
  }

  // 通用提示块
  .guide-callout {
    padding: 14px 16px;
    border-radius: 8px;
    margin-bottom: 16px;

    .callout-label {
      display: inline-block;
      font-size: 13px;
      font-weight: 600;
      margin-bottom: 8px;
      padding: 2px 10px;
      border-radius: 4px;
    }

    p {
      margin: 0;
      font-size: 13px;
      line-height: 1.7;
      color: #606266;

      strong {
        color: #303133;
      }
    }

    &.example {
      background: #f5f7fa;
      border-left: 4px solid var(--el-color-primary);

      .callout-label {
        background: var(--el-color-primary-light-8);
        color: var(--el-color-primary);
      }
    }

    &.warn {
      background: #fef6ec;
      border-left: 4px solid var(--el-color-warning);

      .warn-label {
        background: var(--el-color-danger-light-9);
        color: var(--el-color-danger);
        font-size: 15px;
        font-weight: 700;
      }

      .warn-list {
        margin: 4px 0 0;
        padding-left: 20px;

        li {
          font-size: 13px;
          line-height: 1.9;
          color: #606266;

          strong {
            color: var(--el-color-warning);
            font-weight: 600;
          }
        }
      }
    }
  }
}
</style>
