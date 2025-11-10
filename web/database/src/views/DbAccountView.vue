<template>
  <div class="db-account-view p-6 space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">数据库账号管理</h1>
      <Button @click="openAddDialog">添加账号</Button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-8 text-muted-foreground">
      <Loader2 class="animate-spin inline-block mr-2" /> 加载中...
    </div>

    <!-- 错误提示 -->
    <Alert v-else-if="error" variant="destructive">
      <AlertTitle>加载失败</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <!-- 数据卡片 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card v-for="account in accounts" :key="account.id" class="hover:shadow-lg transition-shadow">
        <CardHeader>
          <div class="flex justify-between items-start">
            <CardTitle>{{ account.name }}</CardTitle>
            <div class="flex space-x-2">
              <Button size="icon" variant="ghost" @click="openEditDialog(account)">
                <Pencil class="h-4 w-4" />
              </Button>
              <Button size="icon" variant="ghost" @click="deleteAccount(account.id)">
                <Trash class="h-4 w-4 text-destructive" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent class="text-sm space-y-2 text-muted-foreground">
          <p><span class="font-medium">主机：</span>{{ account.host }}</p>
          <p><span class="font-medium">端口：</span>{{ account.port }}</p>
          <p><span class="font-medium">创建：</span>{{ formatDate(account.create_time) }}</p>
          <p><span class="font-medium">更新：</span>{{ formatDate(account.update_time) }}</p>
        </CardContent>
      </Card>
    </div>

    <!-- 添加 / 编辑对话框 -->
    <Dialog v-model:open="showDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ isEditing ? '编辑账号' : '添加账号' }}</DialogTitle>
          <DialogDescription>
            请填写数据库账号信息
          </DialogDescription>
        </DialogHeader>

        <form class="space-y-4" @submit.prevent="saveAccount">
          <div class="space-y-2">
            <Label for="name">账号名称</Label>
            <Input id="name" v-model="currentAccount.name" required />
          </div>


          <div class="space-y-2">
            <Label for="host">主机地址</Label>
            <Input id="host" v-model="currentAccount.host" required />
          </div>
          <div class="space-y-2">
            <Label for="host">数据库</Label>
            <Input id="host" v-model="currentAccount.dbname" required />
          </div>
          <div class="space-y-2">
            <Label for="host">账号</Label>
            <Input id="host" v-model="currentAccount.account" required />
          </div>
          <div class="space-y-2">
            <Label for="password">密码</Label>
            <Input id="password" v-model="currentAccount.password" type="password" required />
          </div>

          <div class="space-y-2">
            <Label for="port">端口</Label>
            <Input id="port" v-model.number="currentAccount.port" type="number" required />
          </div>

          <DialogFooter>
            <Button type="button" variant="secondary" @click="closeDialog">取消</Button>
            <Button type="submit">保存</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>

    <!-- 删除确认 -->
    <AlertDialog v-model:open="showDeleteConfirm">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认删除</AlertDialogTitle>
          <AlertDialogDescription>
            确定要删除账号 "{{ accountToDelete?.name }}" 吗？此操作不可撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="cancelDelete">取消</AlertDialogCancel>
          <AlertDialogAction @click="confirmDelete" class="bg-destructive text-white hover:bg-destructive/90">删除</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, type Ref } from "vue"

import { Pencil, Trash, Loader2 } from "lucide-vue-next"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";
import {Label} from "@/components/ui/label";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Alert, AlertDescription, AlertTitle} from "@/components/ui/alert";
import {Input} from "@/components/ui/input";
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from "@/components/ui/alert-dialog";

// 类型定义
interface DBAccount {
  id: number
  name: string
  account: string
  dbname: string
  host: string
  password: string
  port: number
  create_time: string
  update_time: string
}

const accounts: Ref<DBAccount[]> = ref([])
const loading = ref(true)
const error = ref('')
const showDialog = ref(false)
const showDeleteConfirm = ref(false)
const isEditing = ref(false)
const accountToDelete: Ref<DBAccount | null> = ref(null)

const currentAccount: Ref<DBAccount> = ref({
  id: 0,
  name: '',
  account: '',
  dbname: '',
  host: '',
  password: '',
  port: 3306,
  create_time: '',
  update_time: ''
})

const formatDate = (dateString: string) =>
    new Date(dateString).toLocaleString('zh-CN')

const fetchAccounts = async () => {
  try {
    loading.value = true
    const res = await fetch('/api/db-account')
    const result = await res.json()
    if (result.code === 200) accounts.value = result.data
    else error.value = result.msg || '加载失败'
  } catch (e) {
    error.value = '网络错误，请稍后重试'
  } finally {
    loading.value = false
  }
}

const openAddDialog = () => {
  isEditing.value = false
  currentAccount.value = { id: 0, name: '', account: '',dbname: '',host: '', password: '', port: 3306, create_time: '', update_time: '' }
  showDialog.value = true
}

const openEditDialog = (account: DBAccount) => {
  isEditing.value = true
  currentAccount.value = { ...account }
  showDialog.value = true
}

const closeDialog = () => (showDialog.value = false)

const saveAccount = async () => {
  const url = isEditing.value ? `/api/db-account/${currentAccount.value.id}` : '/api/db-account'
  const method = isEditing.value ? 'PUT' : 'POST'
  try {
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(currentAccount.value),
    })
    const result = await res.json()
    if (result.code === 200) {
      showDialog.value = false
      fetchAccounts()
    } else error.value = result.msg || '操作失败'
  } catch (e) {
    error.value = '网络错误'
  }
}

const deleteAccount = (id: number) => {
  const acc = accounts.value.find(a => a.id === id)
  if (acc) {
    accountToDelete.value = acc
    showDeleteConfirm.value = true
  }
}
const cancelDelete = () => (showDeleteConfirm.value = false)
const confirmDelete = async () => {
  if (!accountToDelete.value) return
  try {
    const res = await fetch(`/api/db-account/${accountToDelete.value.id}`, { method: 'DELETE' })
    const result = await res.json()
    if (result.code === 200) fetchAccounts()
  } finally {
    showDeleteConfirm.value = false
  }
}

onMounted(fetchAccounts)
</script>
