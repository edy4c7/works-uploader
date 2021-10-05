<template>
  <v-app-bar app>
    <v-app-bar-nav-icon @click.stop="$emit('click:nav-icon')" />
    <v-spacer></v-spacer>
    <v-btn class="alert" icon>
      <v-icon>mdi-bell</v-icon>
    </v-btn>
    <v-menu activator=".alert" offset-y>
      <v-list>
        <activity
          v-for="act in $store.state.activities.activities"
          :key="act.id"
          :value="act"
        ></activity>
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script lang="ts">
import Vue from 'vue'
import Activity from '~/components/Activity.vue'
export default Vue.extend({
  components: { Activity },
})
</script>
