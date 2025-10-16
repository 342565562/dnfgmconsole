<template>
  <div>
    <panel-title title="点券充值"></panel-title>

    <filter-account @setParams="toSetParams"></filter-account>
    <account-table ref="accountTableRef"></account-table>
  </div>
</template>

<script setup lang="ts">
import PanelTitle from '@/components/PanelTitle'
import AccountTable from './table'
import FilterAccount from './filter'
import { FilterAccountForm } from '@/views/gm/account/model/model'
import { ref } from 'vue'

const accountTableRef = ref(null)

const toSetParams = (data: FilterAccountForm) => {
  accountTableRef.value.setParams(data)
}
</script>

<style scoped></style>
