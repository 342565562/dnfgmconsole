<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="110px" class="mail-form">
           <el-form-item label="物品代码" prop="code">
            <el-input v-model.number="form.code"></el-input>
          </el-form-item>

          <el-form-item label="数量" prop="number">
            <el-input v-model.number="form.number"></el-input>
          </el-form-item>

          <el-form-item label="金币" prop="gold">
            <el-input v-model.number="form.gold">
              <template #append>金币</template>
            </el-input>
          </el-form-item>

          <el-form-item label="邮件类型" prop="mail_type">
            <el-radio-group v-model="form.mail_type" size="large" class="mail-type-group">
              <el-radio label="normal">普通</el-radio>
              <el-radio label="avata">时装</el-radio>
              <el-radio label="creature">宠物</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item>
            <el-button class="send-btn" size="large" @click="sendEmail">发 送</el-button>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col :span="13" :offset="3">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span class="mail-warn"
                ><span class="hl">邮件类型</span>一定要选择正确!!导致<span class="hl">网络中断</span>需要在后台执行"<span class="hl">一键恢复</span>"操作！</span
              >
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
                <el-table-column prop="code" label="物品代码" width="120" />
                <el-table-column prop="name" label="物品名称" />
                <el-table-column prop="rarity" label="稀有度" width="90" />
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
import { ElMessageBox } from 'element-plus'
import { Email, FilterGold } from '@/views/gm/email/model'
import { validate } from '@/utils/element/form'
import { sendEmailByRole, getGolds } from '@/api/gm/email'
import { errorMessage } from '@/utils/element/message'
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
    ElMessageBox.alert(
      '<div style="color:#67c23a;font-weight:600;line-height:2;font-size:15px">发送成功，请小退一下！</div>' +
        '<div style="color:#409eff;line-height:2">如未收到道具，请先在后台操作删除邮件！</div>' +
        '<div style="color:#e6a23c;line-height:2">如网络中断，请执行一件恢复！</div>',
      '发送结果',
      {
        dangerouslyUseHTMLString: true,
        center: true,
        confirmButtonText: '知道了'
      }
    )
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

<style scoped>
/* 警示语：其余文字浅蓝，关键字醒目橙色 */
.mail-warn {
  font-size: 13px;
  color: #409eff;
  font-weight: 500;
}
.mail-warn .hl {
  color: #e6a23c;
  font-weight: 700;
}

/* 表单标签重新设计：加粗醒目 */
.mail-form :deep(.el-form-item__label) {
  font-size: 15px;
  font-weight: 700;
  color: #303133;
}

/* 邮件类型：竖直排列，橙色方框，选中显示对钩 */
.mail-type-group {
  display: flex;
  flex-direction: column;
  gap: 14px;
  align-items: flex-start;
}
.mail-type-group :deep(.el-radio) {
  margin-right: 0;
  height: auto;
}
.mail-type-group :deep(.el-radio__label) {
  font-size: 15px;
  font-weight: 600;
}
/* 橙色方形边框 */
.mail-type-group :deep(.el-radio__inner) {
  width: 18px;
  height: 18px;
  border-radius: 4px;
  border: 2px solid #e6a23c;
}
/* 选中：橙色填充 + 白色对钩 */
.mail-type-group :deep(.el-radio__input.is-checked .el-radio__inner) {
  background: #e6a23c;
  border-color: #e6a23c;
}
.mail-type-group :deep(.el-radio__input.is-checked .el-radio__inner::after) {
  content: '';
  width: 4px;
  height: 8px;
  border: 2px solid #fff;
  border-top: 0;
  border-left: 0;
  border-radius: 0;
  background: transparent;
  transform: rotate(45deg);
  left: 5px;
  top: 1px;
}
.mail-type-group :deep(.el-radio__input.is-checked + .el-radio__label) {
  color: #e6a23c;
}

/* 发送按钮重新设计：醒目 + hover/点击变色 */
.send-btn {
  min-width: 160px;
  height: 44px;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 4px;
  color: #fff;
  border: none;
  border-radius: 8px;
  background: #409eff;
  transition: all 0.2s ease;
}
.send-btn:hover {
  background: #66b1ff;
  box-shadow: 0 4px 14px rgba(64, 158, 255, 0.4);
  transform: translateY(-1px);
}
.send-btn:active {
  background: #2f7fd1;
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(47, 127, 209, 0.5);
}
</style>
