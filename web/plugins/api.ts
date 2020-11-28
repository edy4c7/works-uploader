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

export class ApiClientImpl implements ApiClient {
  private axios: AxiosInstance

  constructor(axios: AxiosInstance) {
    this.axios = axios
  }

  async getWorks() {
    return (await this.axios.get<Work[]>('/works')).data
  }
}

const apiPlugin: Plugin = ({ $axios }, inject) => {
  inject('api', new ApiClientImpl($axios))
}

export default apiPlugin
