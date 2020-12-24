<template>
  <v-overlay
    ref="refOl"
    v-model="isVisible"
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
      <div class="prev">
        <slot name="previousButton" />
      </div>
      <slot name="content" />
      <div class="next">
        <slot name="nextButton" />
      </div>
    </div>
  </v-overlay>
</template>

<style scoped>
.modal {
  display: flex;
  justify-content: center;
  align-items: center;
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

.prev {
  position: absolute;
  top: 50%;
  left: 0;
  text-align: center;
}

.next {
  position: absolute;
  top: 50%;
  right: 0;
  text-align: center;
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
    }
  },
})
</script>
