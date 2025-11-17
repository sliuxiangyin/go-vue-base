import { ref } from 'vue'
import { toast } from 'vue-sonner'
import { useAxios } from './use-axios'
import Cookies from 'universal-cookie'
import { TOKEN_KEY } from '@/utils/constants'

const cookies = new Cookies()
const { axiosInstance } = useAxios()
export interface ServerEvent {
  type: 'success' | 'error' | 'warning' | 'info'
  message: string
  data?: Record<string, any>
  timestamp: string
}

// 全局单例状态,确保只有一个 SSE 连接
const connected = ref(false)
const eventSource = ref<EventSource | null>(null)
const lastEvent = ref<ServerEvent | null>(null)

export function useServerEvents() {

  async function connect() {
    // 如果已经有连接,直接返回
    if (eventSource.value) {
      console.log('SSE connection already exists, reusing...')
      return
    }

    try {
      // 获取 token
      const token = cookies.get(TOKEN_KEY)
      if (!token) {
        console.warn('No token found, skipping SSE connection')
        return
      }

      const baseURL = axiosInstance.defaults.baseURL || ''
      // 通过 URL 参数传递 token（SSE 不支持自定义 header）
      const url = `${baseURL}/admin/events/stream?token=${encodeURIComponent(token)}`
      
      const es = new EventSource(url)

      es.onopen = () => {
        connected.value = true
        console.log('SSE connection established')
      }

      es.onmessage = (event) => {
        try {
          const data: ServerEvent = JSON.parse(event.data)
          lastEvent.value = data

          // 根据事件类型显示 toast
          switch (data.type) {
            case 'success':
              toast.success(data.message, {
                description: data.data ? JSON.stringify(data.data) : undefined,
              })
              break
            case 'error':
              toast.error(data.message, {
                description: data.data ? JSON.stringify(data.data) : undefined,
              })
              break
            case 'warning':
              toast.warning(data.message, {
                description: data.data ? JSON.stringify(data.data) : undefined,
              })
              break
            case 'info':
              toast.info(data.message, {
                description: data.data ? JSON.stringify(data.data) : undefined,
              })
              break
          }
        } catch (error) {
          console.error('Failed to parse SSE message:', error)
        }
      }

      es.onerror = (error) => {
        console.error('SSE connection error:', error)
        connected.value = false
        disconnect()
      }

      eventSource.value = es
    } catch (error) {
      console.error('Failed to establish SSE connection:', error)
    }
  }

  function disconnect() {
    if (eventSource.value) {
      eventSource.value.close()
      eventSource.value = null
      connected.value = false
      console.log('SSE connection closed')
    }
  }

  // 自动连接(仅在没有连接时)
  connect()

  return {
    connected,
    lastEvent,
    connect,
    disconnect,
  }
}
