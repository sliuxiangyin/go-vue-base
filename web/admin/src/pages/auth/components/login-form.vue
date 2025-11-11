<script lang="ts" setup>
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import GitHubButton from './github-button.vue'
import GoogleButton from './google-button.vue'
import PrivacyPolicyButton from './privacy-policy-button.vue'
import TermsOfServiceButton from './terms-of-service-button.vue'
import ToForgotPasswordLink from './to-forgot-password-link.vue'

const router = useRouter()
const authStore = useAuthStore()

// 表单验证规则
const formSchema = toTypedSchema(z.object({
  username: z.string().min(3, '用户名至少需要3个字符').max(50, '用户名最多50个字符'),
  password: z.string().min(6, '密码至少需要6个字符').max(50, '密码最多50个字符'),
}))

const { handleSubmit, defineField, errors } = useForm({
  validationSchema: formSchema,
})

const [username, usernameAttrs] = defineField('username')
const [password, passwordAttrs] = defineField('password')

const isLoading = ref(false)

// 提交登录
const onSubmit = handleSubmit(async (values) => {
  isLoading.value = true
  try {
    const result = await authStore.handleLogin(values)
    if (result.success) {
      toast.success(result.message)
      // 登录成功后跳转到首页
      router.push('/dashboard')
    }
    else {
      toast.error(result.message)
    }
  }
  catch (error) {
    toast.error('登录失败，请稍后重试')
  }
  finally {
    isLoading.value = false
  }
})
</script>

<template>
  <UiCard class="w-full max-w-sm">
    <UiCardHeader>
      <UiCardTitle class="text-2xl">
        登录
      </UiCardTitle>
      <UiCardDescription>
        输入您的用户名和密码登录账户。
        还没有账户？
        <UiButton
          variant="link" class="px-0 text-muted-foreground"
          @click="$router.push('/auth/sign-up')"
        >
          注册
        </UiButton>
      </UiCardDescription>
    </UiCardHeader>
    <UiCardContent>
      <form class="grid gap-4" @submit="onSubmit">
        <div class="grid gap-2">
          <UiLabel for="username">
            用户名
          </UiLabel>
          <UiInput 
            id="username" 
            v-model="username"
            v-bind="usernameAttrs"
            type="text" 
            placeholder="请输入用户名" 
            :disabled="isLoading"
          />
          <p v-if="errors.username" class="text-sm text-destructive">
            {{ errors.username }}
          </p>
        </div>
        <div class="grid gap-2">
          <div class="flex items-center justify-between">
            <UiLabel for="password">
              密码
            </UiLabel>
            <ToForgotPasswordLink />
          </div>
          <UiInput 
            id="password" 
            v-model="password"
            v-bind="passwordAttrs"
            type="password" 
            placeholder="请输入密码" 
            :disabled="isLoading"
          />
          <p v-if="errors.password" class="text-sm text-destructive">
            {{ errors.password }}
          </p>
        </div>

        <UiButton class="w-full" type="submit" :disabled="isLoading">
          <span v-if="isLoading">登录中...</span>
          <span v-else>登录</span>
        </UiButton>

        <UiSeparator label="或使用第三方登录" />

        <div class="flex flex-col items-center justify-between gap-4">
          <GitHubButton />
          <GoogleButton />
        </div>

        <UiCardDescription>
          点击登录即表示您同意我们的
          <TermsOfServiceButton />
          和
          <PrivacyPolicyButton />
        </UiCardDescription>
      </form>
    </UiCardContent>
  </UiCard>
</template>

<style scoped>

</style>
