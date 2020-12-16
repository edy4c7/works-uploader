import Vuex from 'vuex'
import { createLocalVue } from '@vue/test-utils'
import { State, mutations, actions } from '~/store/'
import { Work, ApiClient } from '~/plugins/api'

const data: Work[] = [
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
]

const ApiMock = jest.fn<ApiClient, []>().mockImplementation(() => ({
  getWorks: jest.fn().mockResolvedValue(data),
}))

const localVue = createLocalVue()
localVue.use(Vuex)

describe('actions', () => {
  it('fetch works', async () => {
    const store = new Vuex.Store<State>({
      state: () => ({
        works: [],
      }),
      mutations,
      actions,
    })
    store.$api = new ApiMock()

    await store.dispatch('fetchWorks')
    expect(store.state.works).toEqual(data)

    await store.dispatch('fetchWorks')
    expect(store.state.works).toEqual(data)
  })
})
