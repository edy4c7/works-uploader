import { GetterTree, MutationTree, ActionTree } from 'vuex'
import { Work } from '~/plugins/api'

export const WorkType = {
  URL: 1,
  FILE: 2,
} as const
// eslint-disable-next-line no-redeclare
export type WorkType = typeof WorkType[keyof typeof WorkType]

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
  workForm: {
    type: WorkType.URL,
    title: '',
    description: '',
    contentUrl: '',
  },
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
    this.commit('setWorks', {
      works: await this.$api.getWorks(),
    })
  },
  async postWork(_, payload: PostWorkPayload) {
    await this.$axios.post('/works', payload)
  },
}
