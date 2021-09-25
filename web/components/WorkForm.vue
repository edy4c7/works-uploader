<template>
  <validation-observer v-slot="{ invalid, handleSubmit }">
    <v-form class="work-form" @submit.prevent="handleSubmit(submit)">
      <validation-provider vid="type">
        <v-radio-group
          v-model="work.type"
          row
          mandatory
          name="type"
          class="work-form__type-group"
          @change="input"
        >
          <v-radio
            label="URL"
            :value="TYPE_URL"
            class="work-form__type--url"
          ></v-radio>
          <v-radio
            label="ファイル"
            :value="TYPE_FILE"
            class="work-form__type--file"
          ></v-radio>
        </v-radio-group>
      </validation-provider>
      <validation-provider
        v-slot="{ errors }"
        rules="required|max:40"
        name="title"
      >
        <v-text-field
          v-model="work.title"
          class="work-form__title"
          name="title"
          label="タイトル"
          :error-messages="errors[0]"
          counter="40"
          placeholder="タイトル"
          @input="input"
        ></v-text-field>
      </validation-provider>
      <validation-provider
        v-slot="{ errors }"
        rules="required_if:type,1"
        name="content-url"
      >
        <v-text-field
          v-if="work.type === TYPE_URL"
          v-model="work.contentUrl"
          class="work-form__content-url"
          name="content-url"
          label="URL"
          :error-messages="errors[0]"
          placeholder="https://example.com"
          @input="input"
        ></v-text-field>
      </validation-provider>
      <validation-provider
        v-slot="{ errors }"
        rules="required_if:type,2"
        name="thumbnail"
      >
        <v-file-input
          v-if="work.type === TYPE_FILE"
          v-model="work.thumbnail"
          class="work-form__thumbnail"
          clearable
          label="サムネイル"
          :error-messages="errors[0]"
          @change="input"
        ></v-file-input>
      </validation-provider>
      <validation-provider
        v-slot="{ errors }"
        rules="required_if:type,2"
        name="content"
      >
        <v-file-input
          v-if="work.type === TYPE_FILE"
          v-model="work.content"
          class="work-form__content"
          clearable
          label="作品"
          :error-messages="errors[0]"
          @change="input"
        ></v-file-input>
      </validation-provider>
      <validation-provider
        v-slot="{ errors }"
        rules="max:200"
        name="description"
      >
        <v-textarea
          v-model="work.description"
          class="work-form__description"
          name="description"
          label="説明文"
          counter="200"
          :error-messages="errors[0]"
          @input="input"
        ></v-textarea>
      </validation-provider>
      <v-btn
        color="primary"
        class="work-form__submit"
        :disabled="invalid"
        type="submit"
        >投稿</v-btn
      >
    </v-form>
  </validation-observer>
</template>

<script lang="ts">
import Vue, { PropType } from 'vue'
import { extend, ValidationProvider, ValidationObserver } from 'vee-validate'
// eslint-disable-next-line camelcase
import { required, required_if, max } from 'vee-validate/dist/rules'
import { WorkForm } from '~/plugins/api'
import { WorkType } from '~/store'

extend('required', required)
extend('required_if', required_if)
extend('max', max)

export default Vue.extend({
  components: { ValidationProvider, ValidationObserver },
  props: {
    value: {
      type: Object as PropType<WorkForm>,
      default() {
        return {
          type: WorkType.URL,
          title: '',
          contentUrl: '',
          description: '',
        }
      },
    },
  },
}).extend({
  data() {
    return {
      work: { ...this.value },
    }
  },
  computed: {
    TYPE_URL: () => WorkType.URL,
    TYPE_FILE: () => WorkType.FILE,
  },
  methods: {
    input() {
      this.$emit('input', this.work)
    },
    submit() {
      this.$emit('submit', {
        value: this.work,
      })
    },
  },
})
</script>
