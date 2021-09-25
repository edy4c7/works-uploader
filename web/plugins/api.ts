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

export interface WorkForm {
  type: string
  title: string
  contentUrl?: string
  thumbnail?: File
  content?: File
  description: string
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

const apiClient: ApiClient = {
  getWorks() {
    return Promise.resolve<Work[]>([
      {
        id: '01',
        author: 'taro',
        title: 'hoge',
        description: 'aaaaaaaaaaaaaaaa',
        thumbnailUrl: require('~/assets/image.jpg'),
        contentUrl: 'https://example.com',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '02',
        author: 'hanako',
        title: 'fuga',
        description: 'aaaaaaaaaaaaaaaa',
        thumbnailUrl: require('~/assets/image.jpg'),
        contentUrl: 'https://example.com',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '03',
        author: 'jiro',
        title: 'foo',
        description: 'aaaaaaaaaaaaaaaa',
        thumbnailUrl: require('~/assets/image.jpg'),
        contentUrl: 'https://example.com',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '04',
        author: 'saburo',
        title: 'bar',
        description: 'aaaaaaaaaaaaaaaa',
        thumbnailUrl: require('~/assets/image.jpg'),
        contentUrl: 'https://example.com',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '05',
        author: 'shiro',
        title: 'baz',
        description: 'aaaaaaaaaaaaaaaa',
        thumbnailUrl: require('~/assets/image.jpg'),
        contentUrl: 'https://example.com',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    ])
  },
}

const apiPlugin: Plugin = ({ $axios }, inject) => {
  $axios.setBaseURL('/api/'!)
  // inject('api', new ApiClientImpl($axios))
  inject('api', apiClient)
}

export default apiPlugin
