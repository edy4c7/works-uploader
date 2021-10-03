import { GetterTree, MutationTree, ActionTree } from 'vuex'

export const WorkType = {
  URL: 1,
  FILE: 2,
} as const
// eslint-disable-next-line no-redeclare
export type WorkType = typeof WorkType[keyof typeof WorkType]

export interface Work {
  id: string
  author: string
  title: string
  thumbnailUrl: string
  contentUrl: string
  description: string
  createdAt: Date
  updatedAt: Date
}

export interface WorkForm {
  type: WorkType
  title: string
  contentUrl?: string
  thumbnail?: File
  content?: File
  description: string
}

export interface State {
  works: Work[]
}

export const state: () => State = () => ({
  works: [],
})

export const getters: GetterTree<State, State> = {
  getWorkById(s) {
    return (id: string) => s.works.find((i) => i.id === id)
  },
}

export const mutations: MutationTree<State> = {
  addWorks(s, payload: { works: Work | Work[] }) {
    if (Array.isArray(payload.works)) {
      s.works.push(...payload.works)
    } else {
      s.works.push(payload.works)
    }
  },
  setWorks(s, payload: { works: Work[] }) {
    s.works = payload.works
  },
}

export interface PostWorkPayload {
  value: WorkForm
}

export const actions: ActionTree<State, State> = {
  async fetchWorks() {
    this.commit('works/setWorks', {
      works: await Promise.resolve([
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
      ]),
    })
  },
  async postWork(_, payload: PostWorkPayload) {
    await this.$axios.post('/works', payload)
  },
}
