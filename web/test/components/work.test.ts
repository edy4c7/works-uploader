import Vue from 'vue'
import Vuetify from 'vuetify'
import { createLocalVue, shallowMount } from '@vue/test-utils'
import Work from '~/components/Work.vue'
import { Work as IWork } from '~/store'

Vue.use(Vuetify)

const localVue = createLocalVue()

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

describe('Work', () => {
  let vuetify: Vuetify
  beforeEach(() => {
    vuetify = new Vuetify()
  })

  it('show content title', () => {
    const wrapper = shallowMount(Work, {
      localVue,
      vuetify,
      propsData: {
        content,
      },
    })

    expect(wrapper.html()).toContain(content.title)
  })

  it('show author of content ', () => {
    const wrapper = shallowMount(Work, {
      localVue,
      vuetify,
      propsData: {
        content,
      },
    })

    expect(wrapper.html()).toContain(content.author)
  })

  it('show description of content ', () => {
    const wrapper = shallowMount(Work, {
      localVue,
      vuetify,
      propsData: {
        content,
      },
    })

    expect(wrapper.html()).toContain(content.description)
  })

  it('link to content ', () => {
    const wrapper = shallowMount(Work, {
      localVue,
      vuetify,
      propsData: {
        content,
      },
    })

    expect(wrapper.html()).toContain(`href="${content.contentUrl}"`)
    expect(wrapper.text()).toContain('example.com')
  })
})
