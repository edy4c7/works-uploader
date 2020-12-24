<template>
  <v-card>
    <v-img :src="content.thumbnailUrl" :contains="true" />
    <v-card-title>{{ content.title }}</v-card-title>
    <v-card-subtitle>{{ content.author }}</v-card-subtitle>
    <v-card-text>
      <p class="description">
        {{ content.description }}
      </p>
      <p class="link">
        <v-btn :href="content.contentUrl" target="_blank"
          >{{ domainName }}
        </v-btn>
      </p>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { defineComponent, PropType } from '@nuxtjs/composition-api'
import { Work } from '~/plugins/api'

export default defineComponent({
  props: {
    content: {
      type: Object as PropType<Work>,
      required: true,
    },
  },
  setup(props, _) {
    const domainName = props.content.contentUrl.split('/')[2]

    return {
      domainName,
    }
  },
})
</script>
