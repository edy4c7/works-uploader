import { Plugin } from '@nuxt/types'
import { configure } from 'vee-validate'

const validationPlugin: Plugin = ({ app }, _) => {
  configure({
    defaultMessage(field, values) {
      values._field_ = app.i18n.t(`fields.${field}`)
      return app.i18n
        .t(`validations.messages.${values._rule_}`, values)
        .toString()
    },
  })
}

export default validationPlugin
