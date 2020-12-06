import { MutationTree, ActionTree } from 'vuex'
import { Work } from '~/plugins/api'

export interface State {
  works: Work[]
}

export const state: () => State = () => ({
  works: [],
})

export const mutations: MutationTree<State> = {
  addWorks(s, payload: { works: Work | Work[] }) {
    if (Array.isArray(payload.works)) {
      s.works.push(...payload.works)
    } else {
      s.works.push(payload.works)
    }
  },
}

export const actions: ActionTree<State, State> = {
  async fetchWorks() {
    this.commit('addWorks', {
      works: await this.$api.getWorks(),
    })
  },
}
