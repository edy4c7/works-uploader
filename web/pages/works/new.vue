<template>
  <v-container>
    <works-form></works-form>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import WorksForm from '~/components/WorkForm.vue'

export default defineComponent({
  components: {
    WorksForm,
  },
})
</script>
