<template>
  <v-container>
    <works-form></works-form>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import WorksForm from '~/components/WorkForm.vue'

export default Vue.extend({
  components: {
    WorksForm,
  },
})
</script>
