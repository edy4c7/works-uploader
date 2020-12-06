import { Plugin } from '@nuxt/types'
import { AxiosInstance } from 'axios'
import '@nuxtjs/axios'

export interface Work {
  id: string
  author: string
  title: string
  description: string
  thumbnailUrl: string
  contentUrl: string
  createdAt: Date
  updatedAt: Date
}

export interface ApiClient {
  getWorks(): Promise<Work[]>
}

declare module '@nuxt/types' {
  interface Context {
    $api: ApiClient
  }
}

declare module 'vuex/types/index' {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface Store<S> {
    $api: ApiClient
  }
}

export class ApiClientImpl implements ApiClient {
  private axios: AxiosInstance

  constructor(axios: AxiosInstance) {
    this.axios = axios
  }

  async getWorks() {
    return (await this.axios.get<Work[]>('/works')).data
  }
}

const mockedApiClient: ApiClient = {
  getWorks() {
    return Promise.resolve([
      {
        id: '001',
        author: 'taro',
        title: 'work01',
        description: 'hoge',
        thumbnailUrl: 'https://example.com/work01',
        contentUrl: 'https://example.com/work01',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '002',
        author: 'jiro',
        title: 'work02',
        description: 'fuga',
        thumbnailUrl: 'https://example.com/work02',
        contentUrl: 'https://example.com/work02',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    ])
  },
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const apiPlugin: Plugin = ({ $axios }, inject) => {
  // inject('api', new ApiClientImpl($axios))
  inject('api', mockedApiClient)
}

export default apiPlugin
