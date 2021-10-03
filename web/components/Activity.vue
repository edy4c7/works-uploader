<template>
  <v-list-item class="activity">
    <v-list-item-avatar class="activity__avater">
      <img :src="value.user.picture" alt="avatar" />
    </v-list-item-avatar>
    <v-list-item-content>
      <v-list-item-title class="activity__message">
        {{ message }}
      </v-list-item-title>
      <v-list-item-subtitle class="activity__timestamp">
        {{ value.createdAt }}
      </v-list-item-subtitle>
    </v-list-item-content>
  </v-list-item>
</template>

<script lang="ts">
import Vue, { PropType } from 'vue'
import { Activity, ActivityType } from '~/store/activities'

export default Vue.extend({
  props: {
    value: {
      type: Object as PropType<Activity>,
      required: true,
    },
  },
}).extend({
  computed: {
    message() {
      let key = ''
      switch (this.value.type) {
        case ActivityType.NEW:
          key = 'activities.added'
          break
        case ActivityType.UPDATE:
          key = 'activities.updated'
          break
        default:
          throw new Error('invalid activity type')
      }
      return this.$t(key, {
        user: this.value.user.nickname,
        title: this.value.work.title,
      })
    },
  },
})
</script>
