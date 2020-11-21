<template>
  <v-overlay
    ref="refOl"
    v-model="isVisible"
    v-resized="onWindowResized"
    :style="{ cursor: canOverlayClose ? 'pointer' : 'default' }"
    @click.native="overlayClose"
  >
    <div
      ref="refModal"
      class="modal"
      :style="{ maxWidth: maxWidthOfCloseButton() }"
    >
      <v-btn icon large dark class="close" @click="close">
        <v-icon> mdi-close </v-icon>
      </v-btn>
      <slot />
    </div>
  </v-overlay>
</template>

<style scoped>
.modal {
  width: 100vw;
  max-height: 100vh;
  padding: 44px;
}

.close {
  position: absolute;
  top: 0;
  right: 0;
  text-align: center;
  z-index: 10;
  cursor: pointer;
}
</style>

<script lang="ts">
import Vue from 'vue'
import { defineComponent, ref, watch } from '@vue/composition-api'
import { disableBodyScroll, enableBodyScroll } from 'body-scroll-lock'

export default defineComponent({
  props: {
    isVisible: {
      type: Boolean,
      default: false,
    },
    canOverlayClose: {
      type: Boolean,
      default: false,
    },
    maxWidth: {
      type: [Number, String],
      default: 600,
    },
  },

  setup(props, ctx) {
    const refOl = ref<Vue>()
    const lockedClassName = '--locked'

    function close() {
      ctx.emit('update:isVisible', false)
    }

    function overlayClose() {
      if (props.canOverlayClose) {
        ctx.emit('update:isVisible', false)
      }
    }

    function maxWidthOfCloseButton() {
      if (typeof props.maxWidth === 'number') {
        return `${props.maxWidth}px`
      }
      return props.maxWidth
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    function onWindowResized(evt: UIEvent, el: HTMLElement) {
      ctx.emit('resized')
    }

    watch(
      () => props.isVisible,
      () => {
        const el = refOl.value?.$el as Element
        const html = document.getElementsByTagName('html')[0]
        if (props.isVisible) {
          html.classList.add(lockedClassName)
          disableBodyScroll(el, {
            reserveScrollBarGap: true,
          })
          return
        }
        enableBodyScroll(el)
        html.classList.remove(lockedClassName)
      }
    )

    return {
      refOl,
      close,
      overlayClose,
      maxWidthOfCloseButton,
      onWindowResized,
    }
  },
})
</script>
