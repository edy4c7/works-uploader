<template>
  <div ref="refEl">
    <slot />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from '@vue/composition-api'
import { disableBodyScroll, enableBodyScroll } from 'body-scroll-lock'

export default defineComponent({
  setup() {
    const refEl = ref<HTMLElement>()
    onMounted(() => {
      disableBodyScroll(refEl.value as HTMLElement)
    })

    return {
      refEl,
    }
  },
})
</script>
