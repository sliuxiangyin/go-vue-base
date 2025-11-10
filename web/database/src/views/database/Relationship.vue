<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick, onUnmounted } from 'vue'
import { Button } from "@/components/ui/button";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupText
} from "@/components/ui/input-group";
import { BracesIcon, Fullscreen, CornerDownLeftIcon, RefreshCwIcon, PlayIcon, TrashIcon, CopyIcon } from "lucide-vue-next"
import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { createRelationshipFromSelectedTables, getAllRelationships, deleteRelationship, analyzeRelationship } from '@/lib/relationshipApi';
import type { Relationship, TableInfo } from '@/types/relationship';
import { toast } from 'vue-sonner'
import MarkdownViewer from "@/components/ui/MarkdownViewer.vue";
import {Spinner} from "@/components/ui/spinner";

// 添加 fullscreenRef 和 isFullscreen 状态
const fullscreenRef = ref<HTMLElement | null>(null)
const isFullscreen = ref(false)
const fullscreenContent = ref<{ title: string; result: string } | null>(null)

const props = defineProps<{
  selectedTables?: TableInfo[]
  accountId?: number

}>()



const creating = ref(false)
// 存储获取到的关系记录
const relationships = ref<Relationship[]>([])

const selecteds = ref<TableInfo[]>([])

// 定义响应式数据
const titleInput = ref('')
const loadingStates = ref<Record<number, boolean>>({}) // 用于跟踪每个关系记录的加载状态
watch(
    ()=>props.selectedTables?.length,
    async  (newVal, _) => {
      console.log(newVal)
      await nextTick()
      selecteds.value = [...(props.selectedTables || [])] // ✅ 创建新数组引用

    },{  deep: true }
)

// 计算是否有选中的表
const hasSelectedTables = computed(() => {
  return selecteds.value && selecteds.value.length > 0
})


// 获取所有关系记录
const fetchRelationships = async () => {
  try {
    relationships.value = await getAllRelationships()
  } catch (error: any) {
    toast.error(error.message || '获取关系记录失败')
  }
}

// 删除关系记录
const deleteRelationshipRecord = async (id: number) => {
  try {
    await deleteRelationship(id)
    toast.success('关系记录删除成功')
    // 重新获取数据
    await fetchRelationships()
  } catch (error: any) {
    toast.error(error.message || '删除关系记录失败')
  }
}

// 创建关系记录并自动分析
const createRelationshipRecord = async () => {
  console.log(111);
  if (!props.accountId) {
    toast.error('未找到账号ID')
    return
  }
  if ( selecteds.value.length === 0) {
    toast.error('请先选择至少一个数据表')
    return
  }
  if (titleInput.value==""){
    toast.error('标题不能为空')
    return
  }
  try {
    creating.value = true
    const relationship = await createRelationshipFromSelectedTables(
        props.accountId,
        selecteds.value,
        titleInput.value
    )
    toast.success('关系记录创建成功')
    // 自动开始分析
    await analyzeRelationshipRecord(relationship.id)

    titleInput.value = ''
    selecteds.value=[];
    await fetchRelationships()
  } catch (error: any) {
    toast.error(error.message || '创建关系记录失败')
  } finally {
    creating.value = false
  }
}


// 分析关系记录
const analyzeRelationshipRecord = async (id: number) => {
  try {
    // 设置加载状态
    loadingStates.value[id] = true;
    const sse = analyzeRelationship(
        id,
        (msg) => {
          console.log("分析中:", msg);
        },
        async  () => {
          toast.success('关系分析完成');
          // 重新获取数据以更新列表
          await fetchRelationships();

          loadingStates.value[id] = false;
        }
    );

  } catch (error: any) {
    toast.error(error.message || '关系分析失败');
  } finally {
    // 清除加载状态

  }
};

// 全屏展示函数
const showFullscreen = (title: string, result: string) => {
  fullscreenContent.value = { title, result }
  isFullscreen.value = true
  
  // 等待 DOM 更新后进入全屏
  nextTick(() => {
    if (fullscreenRef.value) {
      enterFullscreen(fullscreenRef.value)
    }
  })
}

// 进入全屏模式
const enterFullscreen = (element: HTMLElement) => {
  try {
    if (element.requestFullscreen) {
      element.requestFullscreen()
    } else if ((element as any).mozRequestFullScreen) { /* Firefox */
      (element as any).mozRequestFullScreen()
    } else if ((element as any).webkitRequestFullscreen) { /* Chrome, Safari & Opera */
      (element as any).webkitRequestFullscreen()
    } else if ((element as any).msRequestFullscreen) { /* IE/Edge */
      (element as any).msRequestFullscreen()
    }
  } catch (error) {
    console.error('无法进入全屏模式:', error)
  }
}

// 退出全屏模式
const exitFullscreen = () => {
  try {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    } else if ((document as any).mozCancelFullScreen) { /* Firefox */
      (document as any).mozCancelFullScreen()
    } else if ((document as any).webkitExitFullscreen) { /* Chrome, Safari & Opera */
      (document as any).webkitExitFullscreen()
    } else if ((document as any).msExitFullscreen) { /* IE/Edge */
      (document as any).msExitFullscreen()
    }
  } catch (error) {
    console.error('无法退出全屏模式:', error)
  }
  isFullscreen.value = false
  fullscreenContent.value = null
}

// 监听全屏状态变化
const handleFullscreenChange = () => {
  if (!document.fullscreenElement) {
    isFullscreen.value = false
    fullscreenContent.value = null
  }
}



// 组件挂载时获取数据
onMounted(() => {
  fetchRelationships()
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

// 组件卸载时移除事件监听器
onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
})
</script>

<template>
  <div class="h-full flex flex-col p-4">

    <ScrollArea class="flex-1 m-h-0 h-full">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 w-full h-full">
        <!-- 显示已有关系记录的InputGroup列表 -->
        <div v-for="relationship in relationships" :key="relationship.id" class="min-w-0">
          <InputGroup>
            <ScrollArea class="h-45 w-full">
              <div class="text-sm p-2">
                <div class="mt-2">
                  <h5 class="font-medium text-xs">表信息:</h5>
                  <ul class="space-y-1 mt-1">
                    <li v-for="(table, index) in JSON.parse(relationship.tables || '[]')" :key="index" class="flex items-start">
                      <span class="mr-2">•</span>
                      <span class="text-sm ">{{ table.name }} - {{ table.comment }}</span>
                    </li>
                  </ul>
                </div>
              </div>
            </ScrollArea>
            <InputGroupAddon align="block-end" class="border-t">
              <div class="text-xs text-gray-500 truncate">
                {{ new Date(relationship.created_time).toLocaleString() }}
              </div>
              <InputGroupButton size="sm" class="ml-auto" variant="ghost" @click="deleteRelationshipRecord(relationship.id)">
                <TrashIcon class="w-4 h-4" />
              </InputGroupButton>
            </InputGroupAddon>
            <InputGroupAddon align="block-start" class="border-b">
              <InputGroupText class="font-mono font-medium truncate" >
                <Badge :variant="relationship.result ? 'success' : 'secondary'">
                  {{ relationship.result ? "success" : "wait" }}
                </Badge>
                {{ relationship.title }}
              </InputGroupText>
              <InputGroupButton 
                class="ml-auto" 
                size="icon-xs" 
                :disabled="loadingStates[relationship.id]"
                @click="analyzeRelationshipRecord(relationship.id)"
              >

                <RefreshCwIcon v-if="!loadingStates[relationship.id]" class="w-4 h-4" />
                <div v-else class="w-4 h-4 rounded-full border-2 border-blue-500 border-t-transparent animate-spin"></div>
              </InputGroupButton>
              <InputGroupButton variant="ghost" size="icon-xs" @click="showFullscreen(relationship.title, relationship.result)">
                <Fullscreen />
              </InputGroupButton>
            </InputGroupAddon>
          </InputGroup>
        </div>
        <!-- 现有的用于创建新关系的InputGroup -->
        <div v-if="hasSelectedTables" class="min-w-0">
          <InputGroup>
            <ScrollArea class="h-45 w-full">
              <div class="text-sm p-2">
                <div class="mt-2">
                  <h5 class="font-medium text-xs">表信息:</h5>
                  <ul class="space-y-1 mt-1">
                    <li v-for="table in selecteds" :key="table.name" class="flex items-center">
                      <span class="mr-2">•</span>
                      <span class="text-sm truncate">{{ table.name }} - {{ table.comment }}</span>
                    </li>
                  </ul>
                </div>
              </div>
            </ScrollArea>
            <InputGroupAddon align="block-end" class="border-t">
              <Input v-model="titleInput" placeholder="标题" />
              <InputGroupButton
                  :disabled="creating"
                  @click="createRelationshipRecord"
                  size="sm"
                  class="ml-auto"
                  variant="default"
              >
                <template v-if="creating">
                  <Spinner /> Running...
                </template>
                <template v-else>
                  Run <CornerDownLeftIcon />
                </template>
              </InputGroupButton>
            </InputGroupAddon>
            <InputGroupAddon align="block-start" class="border-b">
              <InputGroupText class="font-mono font-medium">
                <BracesIcon />
              </InputGroupText>
              <InputGroupButton class="ml-auto" size="icon-xs">
                <RefreshCwIcon />
              </InputGroupButton>
              <InputGroupButton variant="ghost" size="icon-xs" @click="showFullscreen('选中的表', JSON.stringify(selecteds, null, 2))">
                <Fullscreen />
              </InputGroupButton>
            </InputGroupAddon>
          </InputGroup>
        </div>
      </div>
    </ScrollArea>

    <!-- 全屏展示模态框 -->
    <div 
      v-if="isFullscreen" 
      ref="fullscreenRef"
      class="fixed inset-0 z-50 bg-background flex flex-col"
    >
      <div class="flex items-center justify-between p-4 border-b">
        <h2 class="text-xl font-bold">{{ fullscreenContent?.title }}</h2>
        <Button variant="ghost" size="icon" @click="exitFullscreen">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </Button>
      </div>
      <ScrollArea class="flex-1 p-4 min-h-0">
        <MarkdownViewer :source="fullscreenContent?.result" />
      </ScrollArea>
    </div>
  </div>
</template>

<style scoped>
/* 可以添加一些自定义样式 */
</style>