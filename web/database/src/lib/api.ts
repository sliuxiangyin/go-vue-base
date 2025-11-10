import axios, {type AxiosInstance, type AxiosResponse} from 'axios';

// 创建axios实例
const apiClient: AxiosInstance = axios.create({
  baseURL: '/api', // 默认API基础URL
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config: any)  => {
    // 在发送请求之前做些什么
    // 例如：添加认证token
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    // 对请求错误做些什么
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // 对响应数据做点什么
    return response.data;
  },
  (error) => {
    // 对响应错误做点什么
    if (error.response?.status === 401) {
      // 处理未授权错误
      // 例如：重定向到登录页面
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;

// 封装常用的HTTP方法
