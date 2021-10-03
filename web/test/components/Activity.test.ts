import Vue from 'vue'
import Vuetify from 'vuetify'
import Vuei18n from 'vue-i18n'
import { createLocalVue, shallowMount } from '@vue/test-utils'
import AcitvityComponent from '~/components/Activity.vue'

Vue.use(Vuetify)
Vue.use(Vuei18n)

const localVue = createLocalVue()
const picture =
  'https://cdn.pixabay.com/photo/2021/05/02/08/33/jellyfish-6222849_960_720.jpg'

describe('Activity component', () => {
  let vuetify: Vuetify
  let i18n: Vuei18n

  beforeEach(() => {
    vuetify = new Vuetify()
    i18n = new Vuei18n({
      locale: 'en',
      messages: {
        en: {
          activities: {
            added: '{user} added {title}',
            updated: '{user} updated {title}',
          },
        },
      },
    })
  })

  it('render users avater', () => {
    const wrapper = shallowMount(AcitvityComponent, {
      localVue,
      vuetify,
      i18n,
      propsData: {
        value: {
          id: 1234,
          type: 1,
          user: {
            id: 'aaaaa',
            name: 'XXX XXXX',
            nickname: 'XXX',
            picture,
          },
          work: {
            id: '000',
            author: 'XXXX',
            title: 'YYY',
            contentUrl: '',
            description: '',
            thumbnailUrl: '',
            createdAt: new Date(),
            updatedAt: new Date(),
          },
          createdAt: new Date(),
        },
      },
    })

    expect(wrapper.find('.activity__avater > img').attributes('src')).toBe(
      picture
    )
  })

  it('render added message', () => {
    const wrapper = shallowMount(AcitvityComponent, {
      localVue,
      vuetify,
      i18n,
      propsData: {
        value: {
          id: 1234,
          type: 1,
          user: {
            id: 'aaaaa',
            name: 'XXX XXXX',
            nickname: 'XXX',
            picture,
          },
          work: {
            id: '000',
            author: 'XXXX',
            title: 'YYY',
            contentUrl: '',
            description: '',
            thumbnailUrl: '',
            createdAt: new Date(),
            updatedAt: new Date(),
          },
          createdAt: new Date(),
        },
      },
    })

    expect(wrapper.find('.activity__message').text()).toContain('XXX added YYY')
  })

  it('render updated message', () => {
    const wrapper = shallowMount(AcitvityComponent, {
      localVue,
      vuetify,
      i18n,
      propsData: {
        value: {
          id: 1234,
          type: 2,
          user: {
            id: 'aaaaa',
            name: 'XXX XXXX',
            nickname: 'XXX',
            picture,
          },
          work: {
            id: '000',
            author: 'XXXX',
            title: 'YYY',
            contentUrl: '',
            description: '',
            thumbnailUrl: '',
            createdAt: new Date(),
            updatedAt: new Date(),
          },
          createdAt: new Date(),
        },
      },
    })

    expect(wrapper.find('.activity__message').text()).toContain(
      'XXX updated YYY'
    )
  })

  it('render timestamp', () => {
    const createdAt = new Date()
    const wrapper = shallowMount(AcitvityComponent, {
      localVue,
      vuetify,
      i18n,
      propsData: {
        value: {
          id: 1234,
          type: 2,
          user: {
            id: 'aaaaa',
            name: 'XXX XXXX',
            nickname: 'XXX',
            picture,
          },
          work: {
            id: '000',
            author: 'XXXX',
            title: 'YYY',
            contentUrl: '',
            description: '',
            thumbnailUrl: '',
            createdAt: new Date(),
            updatedAt: new Date(),
          },
          createdAt,
        },
      },
    })

    expect(wrapper.find('.activity__timestamp').text()).toBe(
      createdAt.toString()
    )
  })
})
