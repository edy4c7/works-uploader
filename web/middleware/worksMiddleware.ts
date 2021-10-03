import { Context } from '@nuxt/types'

export default async function ({ store }: Context) {
  await store.dispatch('works/fetchWorks')
}
