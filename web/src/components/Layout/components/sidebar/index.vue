<template>
  <el-scrollbar wrap-class="scrollbar-wrapper">
    <div class="sidebar-header">
      <span class="custom-logo" v-html="customLogoSvg" :class="[!isCollapse ? 'show_log' : 'hide_log']" />
      <span v-show="!isCollapse">{{ settings.title }}</span>
    </div>

    <el-menu
      :default-active="activeMenu"
      :collapse="isCollapse"
      unique-opened
      mode="vertical"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
      :collapse-transition="false"
    >
      <sidebar-item v-for="route in routes" :key="route.path" :item="route" :base-path="route.path"></sidebar-item>
    </el-menu>
  </el-scrollbar>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SidebarItem from './SidebarItem.vue'
import settings from '@/settings'
import { usePermissionStore } from '@/store/modules/permission'
import { useAppStore } from '@/store/modules/app'
import { storeToRefs } from 'pinia'

export default defineComponent({
  name: 'SidebarMenu',
  components: {
    SidebarItem
  },
  setup() {
    // const router = useRouter().options.routes
    const permissionStore = usePermissionStore()
    const appStore = useAppStore()

    const { sidebar } = storeToRefs(appStore)
    const { routers } = storeToRefs(permissionStore)

    const route = useRoute()

    const activeMenu = computed(() => {
      const { meta, path } = route
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    })

    const customLogoSvg = `<svg t="1752830120226" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1854" width="57" height="57"><path d="M308.235 298.34S146.999 197.517 58.661 206.385c-8.907 0.434-13.079 1.315-21.819 2.856-9.389 1.655-18.362 5.317-26.744 9.957 0 0 184.802 43.354 255.349 102.339 0 0-16.563 39.571-22.084 42.3 0 0 10.352-30.702 9.663-38.207-0.69-7.505-73.327 157.814-161.492 240.153 0 0-38.648 19.104 2.762 30.021l15.182-4.094 19.323 6.824h16.563l16.563-9.553s109.732-8.87 118.705-16.375c0 0 84.886-28.653 88.335-55.945l-2.76 24.563 26.224-34.115c-0.001 0.002 45.549-124.851-84.196-208.769zM186.369 555.982l91.127-209.377c140.11 160.564-91.127 209.377-91.127 209.377z m799.578-76.16l-309.181 13.644S714.49 392.492 981.57 377.483l32.332-2.008-17.148-14.366s-349.663-14.328-411.086 136.452l-16.563 2.729 33.127-68.225-13.802 20.466 6.9-38.206-28.986 35.478 6.904-23.196s-59.354 78.459-74.537 80.504c0 0-28.295-115.301-84.197-60.038l12.423-1.366s-32.436 60.723-34.508 61.404c0 0-35.196 59.357-24.843 73.685 0 0 38.647-66.862 44.168-65.497l-12.423 28.655 38.648-45.031s-8.282 83.236 45.549 69.591l35.887-16.374-19.324 15.009s47.596-27.546 55.482-44.74l1.108-4.38c0 1.356-0.396 2.827-1.108 4.38l-5.795 22.908 6.903 8.188-109.038 226.987 38.648-65.496-42.86 103.17 31.745-37.909 52.523-96.646-26.227 76.413s91.099-192.192 125.604-211.978L799.61 534.4l190.479-28.654-81.437-5.457 77.295-20.467z" fill="#ffffff" p-id="1855"></path></svg>`;
    return {
      routes: routers,
      activeMenu,
      isCollapse: computed(() => !sidebar.value.opened),
      settings,
      customLogoSvg
    }
  }
})
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import '../../../../styles/mixin.scss';

.sidebar-container {
  .el-menu {
    height: calc(100% - 50px) !important;
  }
}

.show_log {
  width: 38px;
  height: 38px;
  color: #ffffff;
  margin-left: 10px;
}

.hide_log {
  width: 33px;
  height: 33px;
  color: #ffffff;
}

.sidebar-header {
  @include backGround(#324157);
  border-bottom: 1px solid #333;
  width: 100%;
  height: 50px;
  padding: 5px 10px;

  span {
    color: #ffffff;
    font-size: 20px;
    position: absolute;
    left: 75px;
    top: 13px;
  }
}

.custom-logo {
  display: inline-block;
  vertical-align: middle;
  width: 57px;
  height: 57px;
  margin-top: -15px;
  margin-left: -60px;
}
</style>

<style rel="stylesheet/scss" lang="scss">
body {
  .router-link-active {
    .el-menu-item {
      color: rgb(64, 158, 255) !important;
    }
  }
}
</style> 