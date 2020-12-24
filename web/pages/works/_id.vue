<template>
  <v-container class="p-work">
    <work :content="content" />
  </v-container>
</template>

<style scoped>
.p-work {
  width: 100vw;
  max-width: 600px;
}
</style>

<script lang="ts">
import {
  computed,
  defineComponent,
  getCurrentInstance,
} from '@nuxtjs/composition-api'
import Work from '~/components/Work.vue'

export default defineComponent({
  components: {
    Work,
  },
  setup(_, __) {
    const self = getCurrentInstance()

    const content = computed(() =>
      self?.$store.getters.getWorkById(self?.$route.params.id)
    )

    return {
      content,
    }
  },
})
</script>
