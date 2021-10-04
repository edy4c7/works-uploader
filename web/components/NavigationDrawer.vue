<template>
  <v-navigation-drawer
    :value="value"
    :mini-variant="miniVariant"
    fixed
    :clipped="clipped"
    temporary
    app
    @input="$emit('input', $event)"
  >
    <v-list nav>
      <v-list-item
        v-for="(item, i) in items"
        :key="i"
        :to="item.to"
        router
        exact
      >
        <v-list-item-action>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title v-text="item.title" />
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import Vue from 'vue'
export default Vue.extend({
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      clipped: false,
      items: [
        {
          icon: 'mdi-apps',
          title: 'Welcome',
          to: '/',
        },
        {
          icon: 'mdi-chart-bubble',
          title: 'New',
          to: '/works/new',
        },
      ],
      miniVariant: false,
    }
  },
})
</script>
