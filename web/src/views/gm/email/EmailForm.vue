<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
           <el-form-item label="物品代码" prop="code">
            <el-input v-model.number="form.code"></el-input>
          </el-form-item>

          <el-form-item label="数量" prop="number">
            <el-input v-model.number="form.number"></el-input>
          </el-form-item>

          <el-form-item label="金币" prop="gold">
            <el-input v-model.number="form.gold">
              <template #append>万</template>
            </el-input>
          </el-form-item>

          <el-form-item label="邮件类型" prop="mail_type">
            <el-radio-group v-model="form.mail_type">
              <el-radio label="normal">普通</el-radio>
              <el-radio label="avata">时装</el-radio>
              <el-radio label="creature">宠物</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" size="small" @click="sendEmail">发送</el-button>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col :span="13" :offset="3">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>邮件类型一定要选择正确!!导致网络中断需要在后台执行“一键恢复”操作！</span>
            </div>
          </template>
          <div>
            <el-row>
              <el-form :inline="true" :model="argQuery" ref="filterFormRef" label-width="100px">
                <div>
                  <el-form-item label="物品名称:" prop="name">
                    <el-input v-model="argQuery.name" clearable></el-input>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="filterGold">查询</el-button>
                    <el-button @click.prevent="reset">重置</el-button>
                  </el-form-item>
                </div>
              </el-form>
            </el-row>
            <el-row>
              <el-table :data="goldState.data" v-loading="goldState.loading" border fit>
                <el-table-column prop="code" label="物品代码" width="150" />
                <el-table-column prop="name" label="物品名称" />
              </el-table>
            </el-row>
            <el-row>
              <el-pagination
                :current-page="pageQuery.page"
                :page-sizes="[10, 20, 30, 40, 50]"
                :page-size="pageQuery.per_page"
                :total="goldState.total"
                background
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              ></el-pagination>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, defineExpose } from 'vue'
import { FormInstance, FormRules } from 'element-plus'
import { Email, FilterGold } from '@/views/gm/email/model'
import { validate } from '@/utils/element/form'
import { sendEmailByRole, getGolds } from '@/api/gm/email'
import { errorMessage, successMessage } from '@/utils/element/message'
import { PageQuery } from '@/models/page'

const formRef = ref<FormInstance>()

const form = reactive<Email & { mail_type?: string }>({
  code: null,
  number: null,
  seperate_upgrade: 0,
  upgrade: 0,
  is_amplify: false,
  amplify_option: 3,
  amplify_value: 0,
  gold: 0,
  seal_flag: false,
  mail_type: 'normal',
  avata_flag: false,
  creature_flag: false
})

const rules = reactive<FormRules>({
  code: [{ type: 'integer', min: 0, message: '请输入物品代码', trigger: 'blur' }],
  number: [
    { required: true, message: '请输入数量', trigger: 'blur' },
    { type: 'number', min: 1, message: '数量必须大于0', trigger: 'blur' }
  ]
})

const characNo = ref<number>(null)

const options = [
  { label: '体力', value: 1 },
  { label: '精神', value: 2 },
  { label: '力量', value: 3 },
  { label: '智力', value: 4 }
]

const goldState = reactive<any>({
  data: [],
  total: 0,
  loading: false
})

const pageQuery = reactive<PageQuery>({
  page: 1,
  per_page: 10
})

let argQuery = reactive<FilterGold>({
  name: ''
})

const filterFormRef = ref<FormInstance>()

const sendEmail = async () => {
  const valid = await validate(formRef)
  if (!valid) return

  if (!characNo.value) {
    errorMessage('未选择角色ID')
    return
  }

  if (typeof form.gold === 'string') {
    form.gold = 0
  }

  if (form.mail_type === 'avata') {
    form.avata_flag = true
    form.creature_flag = false
  } else if (form.mail_type === 'creature') {
    form.avata_flag = false
    form.creature_flag = true
  } else {
    form.avata_flag = false
    form.creature_flag = false
  }

  try {
    await sendEmailByRole(characNo.value, form)
    successMessage('发送成功，请小退一下！')
  } catch (e) {}
}

const setCharacNo = (id: number) => {
  characNo.value = id
}

const getGoldList = async () => {
  goldState.loading = true
  const params = {
    ...argQuery,
    ...pageQuery
  }
  const { items, total } = await getGolds(params)
  goldState.data = items
  goldState.total = total
  goldState.loading = false
}

const handleSizeChange = (params_limit: number) => {
  pageQuery.per_page = params_limit
  getGoldList()
}

const handleCurrentChange = (params_page: number) => {
  pageQuery.page = params_page
  getGoldList()
}

const filterGold = () => {
  pageQuery.page = 1
  getGoldList()
}

const reset = () => {
  filterFormRef.value.resetFields()
  pageQuery.page = 1
  getGoldList()
}

// hook
getGoldList()
defineExpose({ setCharacNo })
</script>

<style scoped></style>
