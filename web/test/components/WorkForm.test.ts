/* eslint-disable @typescript-eslint/no-unused-vars */
import Vue from 'vue'
import Vuetify from 'vuetify'
import flushPromises from 'flush-promises'
import { createLocalVue, mount } from '@vue/test-utils'
import WorkForm from '~/components/WorkForm.vue'
import { WorkType } from '~/store/works'

const localVue = createLocalVue()
Vue.use(Vuetify)

describe('WorkForm', () => {
  let vuetify: Vuetify
  beforeEach(() => {
    vuetify = new Vuetify()
  })

  describe('on mount', () => {
    it('not specified prop', () => {
      const wrapper = mount(WorkForm, {
        localVue,
        vuetify,
      })

      expect(wrapper.find('.work-form__type-group').props().value).toBe(
        WorkType.URL
      )
      expect(wrapper.find('.work-form__title').props().value).toBe('')
      expect(wrapper.find('.work-form__content-url').props().value).toBe('')
      expect(wrapper.find('.work-form__description').props().value).toBe('')
    })

    it('specified prop(type is URL)', () => {
      const value = {
        type: WorkType.URL,
        title: 'testTitle',
        contentUrl: 'https://example.com',
        description: 'testDescription',
      }

      const wrapper = mount(WorkForm, {
        localVue,
        vuetify,
        propsData: { value },
      })

      expect(wrapper.find('.work-form__type-group').props().value).toBe(
        value.type
      )
      expect(wrapper.find('.work-form__title').props().value).toBe(value.title)
      expect(wrapper.find('.work-form__content-url').props().value).toBe(
        value.contentUrl
      )
      expect(wrapper.find('.work-form__description').props().value).toBe(
        value.description
      )
    })

    it('specified prop(type is File)', () => {
      const value = {
        type: WorkType.FILE,
        title: 'testTitle',
        thumbnail: {
          type: 'image/png',
        },
        content: {
          type: 'application/zip',
        },
        description: 'testDescription',
      }

      const wrapper = mount(WorkForm, {
        localVue,
        vuetify,
        propsData: { value },
      })

      expect(wrapper.find('.work-form__type-group').props().value).toBe(
        value.type
      )
      expect(wrapper.find('.work-form__title').props().value).toBe(value.title)
      expect(wrapper.find('.work-form__thumbnail').props().value).toBe(
        value.thumbnail
      )
      expect(wrapper.find('.work-form__content').props().value).toBe(
        value.content
      )
      expect(wrapper.find('.work-form__description').props().value).toBe(
        value.description
      )
    })
  })

  describe('on input', () => {
    it('type is URL', async () => {
      const wrapper = mount(WorkForm, {
        localVue,
        vuetify,
        propsData: {
          value: {
            type: WorkType.FILE,
          },
        },
      })

      let times = 0

      wrapper.find('.work-form__type--url').trigger('click')
      await Vue.nextTick()
      expect(wrapper.emitted().input?.[times][0]?.type).toStrictEqual(
        WorkType.URL
      )
      expect(wrapper.find('.work-form__thumbnail').exists()).toBeFalsy()
      expect(wrapper.find('.work-form__content').exists()).toBeFalsy()

      const title = 'testtitle'
      wrapper
        .find('.work-form__title')
        .find('input[type="text"]')
        .setValue(title)
      expect(wrapper.emitted().input?.[++times][0]?.title).toStrictEqual(title)

      const url = 'https://example.com'
      wrapper
        .find('.work-form__content-url')
        .find('input[type="text"]')
        .setValue(url)
      expect(wrapper.emitted().input?.[++times][0]?.contentUrl).toStrictEqual(
        url
      )

      const description = 'testdescription'
      wrapper
        .find('.work-form__description')
        .find('textarea')
        .setValue(description)
      expect(wrapper.emitted().input?.[++times][0]?.description).toStrictEqual(
        description
      )
    })
    it('type is File', async () => {
      const wrapper = mount(WorkForm, {
        localVue,
        vuetify,
        propsData: {
          value: {
            type: WorkType.URL,
          },
        },
      })

      wrapper.find('.work-form__type--file').trigger('click')
      await Vue.nextTick()
      expect(wrapper.emitted().input?.[0][0]?.type).toStrictEqual(WorkType.FILE)
      expect(wrapper.find('.work-form__title').exists()).toBeTruthy()
      expect(wrapper.find('.work-form__content-url').exists()).toBeFalsy()
      expect(wrapper.find('.work-form__thumbnail').exists()).toBeTruthy()
      expect(wrapper.find('.work-form__content').exists()).toBeTruthy()
      expect(wrapper.find('.work-form__description').exists()).toBeTruthy()
    })
  })

  describe('on submit', () => {
    it('type is URL', async () => {
      const el = document.createElement('div')
      el.id = 'root'
      document.body.appendChild(el)

      const wrapper = mount(WorkForm, {
        attachTo: '#root',
        localVue,
        vuetify,
      })

      const value = {
        type: WorkType.URL,
        title: 'testTitle',
        contentUrl: 'https://example.com',
        description: 'testDescription',
      }

      wrapper.find('.work-form__type--url').trigger('click')
      wrapper
        .find('.work-form__title')
        .find('input[type="text"]')
        .setValue(value.title)
      wrapper
        .find('.work-form__content-url')
        .find('input[type="text"]')
        .setValue(value.contentUrl)
      wrapper
        .find('.work-form__description')
        .find('textarea')
        .setValue(value.description)

      wrapper.find('.work-form__submit').trigger('click')
      await flushPromises()
      expect(wrapper.emitted().submit?.[0][0]).toStrictEqual({ value })

      wrapper.destroy()
    })
  })

  describe('on error', () => {
    describe('type is URL', () => {
      it('missing title', () => {
        const wrapper = mount(WorkForm, {
          localVue,
          vuetify,
        })

        const value = {
          type: WorkType.URL,
          title: 'testTitle',
          contentUrl: 'https://example.com',
          description: 'testDescription',
        }

        wrapper.find('.work-form__type--url').trigger('click')
        wrapper
          .find('.work-form__content-url')
          .find('input[type="text"]')
          .setValue(value.contentUrl)
        wrapper
          .find('.work-form__description')
          .find('textarea')
          .setValue(value.description)

        wrapper.find('.work-form__submit').trigger('click')
        expect(wrapper.emitted().submit).toBeFalsy()
      })
    })
  })
})
