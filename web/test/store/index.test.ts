import Vuex from 'vuex'
import { createLocalVue } from '@vue/test-utils'
import { State, getters, Work } from '~/store/works'

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

const localVue = createLocalVue()
localVue.use(Vuex)

describe('getters', () => {
  it('get work by id', () => {
    const store = new Vuex.Store<State>({
      state: () => ({
        works: data,
      }),
      getters,
    })

    const result = store.getters.getWorkById('001') as Work

    expect(result.id).toEqual('001')
  })
})
