<template>
  <TooltipProvider :delay-duration="0">
    <ResizablePanelGroup
        id="resize-panel-group-1"
        direction="horizontal"
        class="h-full items-stretch  rounded-md overflow-hidden p-0"
    >
      <!-- 左侧：数据库表列表 -->
      <ResizablePanel :default-size="20" class="flex flex-col h-full">
        <TableList :tables="tables" @selection-change="handleSelectionChange" />
      </ResizablePanel>
      <ResizableHandle/>


      <!-- 右侧：聊天窗口 -->
      <ResizablePanel class=" h-full">
        <Tabs default-value="unread" class="h-full flex flex-col">
          <div class="flex items-center px-4 py-2">
            <h1 class="text-xl font-bold">
              Inbox
            </h1>
            <TabsList class="ml-auto">
              <TabsTrigger value="all" class="text-zinc-600 dark:text-zinc-200">
               <Bubbles/>  关系一览表
              </TabsTrigger>
              <TabsTrigger value="unread" class="text-zinc-600 dark:text-zinc-200">
                <MessageSquareMore/> 对话
              </TabsTrigger>
            </TabsList>
          </div>
          <Separator />

          <TabsContent value="all" class="flex flex-col  flex-1 m-0">
             <Relationship :selected-tables="selectedTables" :account-id="accountId"/>
          </TabsContent>
          <TabsContent value="unread" class="flex flex-col flex-1 m-0">
            <Chat/>
          </TabsContent>
        </Tabs>

      </ResizablePanel>
    </ResizablePanelGroup>
  </TooltipProvider>
</template>
<script setup lang="ts">
import {ref, nextTick, onUpdated, onMounted} from 'vue'
import {
  TooltipProvider,
} from '@/components/ui/tooltip'

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable'
import { MessageSquareMore,Bubbles } from "lucide-vue-next"
import {Button} from '@/components/ui/button'
import {Textarea} from "@/components/ui/textarea";
import {useRoute} from "vue-router";
import api from "@/lib/api.ts";
import TableList from './TableList.vue'
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs";
import {Separator} from "@/components/ui/separator";
import Relationship from "@/views/database/Relationship.vue";
import Chat from "@/views/database/Chat.vue";

const route = useRoute()

interface Message {
  id: number
  sender: 'user' | 'bot'
  content: string
}

interface TableInfo {
  name: string
  comment: string
  selected?: boolean
}



const tables = ref<TableInfo[]>([])
const selectedTables = ref<TableInfo[]>([])
const accountId = ref<number>(0)



onMounted(() => {
  init()
})
const handleSelectionChange = (selected: TableInfo[]) => {
  selectedTables.value = selected
}
const init = async () => {
  const id = route.params.id
  accountId.value = Number(id)
  const response = await api.get<TableInfo[]>(`/database/init/${id}`);
  tables.value = response.data;

}


</script>
<style scoped>
/* 可选：隐藏滚动条样式更干净 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}
</style>
