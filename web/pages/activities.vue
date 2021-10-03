<template>
  <v-menu value="true">
    <v-list>
      <activity
        v-for="act in $store.state.activities.activities"
        :key="act.id"
        :value="act"
      ></activity>
    </v-list>
  </v-menu>
</template>

<script lang="ts">
import Vue from 'vue'
import Activity from '~/components/Activity.vue'
export default Vue.extend({
  components: { Activity },
})
</script>
