<script setup lang="ts">

import {Textarea} from "@/components/ui/textarea";
import {Button} from "@/components/ui/button";
import {nextTick, onUpdated, ref} from "vue";

interface Message {
  id: number
  sender: 'user' | 'bot'
  content: string
}

// 假数据
const messages = ref<Message[]>([
  {id: 1, sender: 'bot', content: '你好！我是你的智能助手。'},
  {id: 2, sender: 'user', content: '你好，请帮我列出数据库表。'},
  {id: 3, sender: 'bot', content: '好的，请在左侧选择你要查看的表。'},
])
const inputMessage = ref('')
const chatContainer = ref<HTMLElement | null>(null)

// 模拟发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim()) return

  messages.value.push({
    id: Date.now(),
    sender: 'user',
    content: inputMessage.value,
  })

  const userMsg = inputMessage.value
  inputMessage.value = ''

  await nextTick()
  scrollToBottom()

  // 模拟机器人回复
  setTimeout(async () => {
    messages.value.push({
      id: Date.now() + 1,
      sender: 'bot',
      content: `你说的是「${userMsg}」吗？`,
    })
    await nextTick()
    scrollToBottom()
  }, 600)
}



const scrollToBottom = () => {
  const el = chatContainer.value
  if (el) {
    el.scrollTo({
      top: el.scrollHeight,
      behavior: 'smooth',
    })
  }
}

// 页面每次更新（例如加载假数据）也保持滚动到底部
onUpdated(() => {
  scrollToBottom()
})
</script>

<template>
    <div ref="chatContainer" class="flex-1 overflow-y-auto p-4 space-y-3">
      <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex"
          :class="msg.sender === 'user' ? 'justify-end' : 'justify-start'"
      >
        <div
            class="max-w-xs px-3 py-2 rounded-lg "
            :class="msg.sender === 'user'
                ? 'flex w-max max-w-[75%] flex-col gap-2 rounded-lg px-3 py-2 text-sm ml-auto bg-primary text-primary-foreground'
                : 'flex w-max max-w-[75%] flex-col gap-2 rounded-lg px-3 py-2 text-sm bg-muted'"
        >
          {{ msg.content }}
        </div>
      </div>
    </div>
    <!-- 输入框固定在底部 -->
    <div class="border-t p-3 flex items-center space-x-2">
          <Textarea
              v-model="inputMessage"
              placeholder="输入消息..."
              class="flex-1"
              @keyup.enter="sendMessage"
          />
      <Button @click="sendMessage">发送</Button>
    </div>
</template>

<style scoped>

</style>