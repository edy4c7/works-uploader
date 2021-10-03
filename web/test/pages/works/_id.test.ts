import Vue from 'vue'
import Vuex from 'vuex'
import Vuetify from 'vuetify'
import { createLocalVue, shallowMount } from '@vue/test-utils'
import Works from '~/pages/works/_id.vue'
import { Work as IWork } from '~/store'

Vue.use(Vuetify)

const localVue = createLocalVue()
localVue.use(Vuex)

const content: IWork = {
  id: '01',
  author: 'taro',
  title: 'hoge',
  description: 'aaaaaaaaaaaaaaaa',
  thumbnailUrl: 'https://example.com/thumbnail',
  contentUrl: 'https://example.com',
  createdAt: new Date(),
  updatedAt: new Date(),
}

describe('Works/_id', () => {
  let vuetify: Vuetify
  const getWorkById = jest.fn((_) => content)

  beforeEach(() => {
    vuetify = new Vuetify()
  })

  it('show work by id', () => {
    const id = '001'
    shallowMount(Works, {
      localVue,
      vuetify,
      store: new Vuex.Store({
        getters: {
          getWorkById: (_) => getWorkById,
        },
      }),
      mocks: {
        $route: {
          params: {
            id,
          },
        },
      },
    })

    expect(getWorkById).toBeCalledWith(id)
  })
})
