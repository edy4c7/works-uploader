import colors from 'vuetify/es5/util/colors'

export default {
  // Target (https://go.nuxtjs.dev/config-target)
  target: 'static',
  // Global page headers (https://go.nuxtjs.dev/config-head)
  head: {
    titleTemplate: '%s - Works uploader',
    title: 'Works uploader',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: '' },
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
  },

  env: {
    apiUrl:
      process.env.NODE_ENV === 'production'
        ? ''
        : 'http://localhost:8000/api/v1',
  },

  // Global CSS (https://go.nuxtjs.dev/config-css)
  css: ['@/assets/main.css'],

  // Plugins to run before rendering page (https://go.nuxtjs.dev/config-plugins)
  plugins: ['@/plugins/validation'],

  // Auto import components (https://go.nuxtjs.dev/config-components)
  components: true,

  // Modules for dev and build (recommended) (https://go.nuxtjs.dev/config-modules)
  buildModules: [
    // https://go.nuxtjs.dev/typescript
    '@nuxt/typescript-build',
    // https://go.nuxtjs.dev/stylelint
    '@nuxtjs/stylelint-module',
    // https://go.nuxtjs.dev/vuetify
    '@nuxtjs/vuetify',
    '@nuxtjs/composition-api',
  ],

  // Modules (https://go.nuxtjs.dev/config-modules)
  modules: [
    // https://go.nuxtjs.dev/axios
    '@nuxtjs/axios',
    '@nuxtjs/i18n',
  ],

  // Axios module configuration (https://go.nuxtjs.dev/config-axios)
  axios: {
    proxy: true,
  },

  i18n: {
    locales: [
      { code: 'ja', iso: 'ja_JP', file: 'ja.yaml' },
      { code: 'en', iso: 'en-US', file: 'en.yaml' },
    ],
    defaultLocale: 'ja',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      fallbackLocale: 'en',
      redirectOn: 'root',
    },
    lazy: true,
    langDir: '@/i18n/',
  },

  proxy: {
    '/api/':
      process.env.NODE_ENV === 'production'
        ? ''
        : 'http://localhost:8000/api/v1',
  },

  // Vuetify module configuration (https://go.nuxtjs.dev/config-vuetify)
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    theme: {
      dark: true,
      themes: {
        dark: {
          primary: colors.blue.darken2,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
      },
    },
  },

  // Build Configuration (https://go.nuxtjs.dev/config-build)
  build: {
    transpile: ['vee-validate/dist/rules'],
    extend(config, _) {
      config.module.rules.push({
        test: /\.ya?ml$/,
        type: 'json', // Required by Webpack v4
        use: 'yaml-loader',
      })
    },
  },

  srcDir: 'web/',

  generate: {
    dir: 'public',
    subFolders: false,
    routes: ['/works/01', '/works/02', '/works/03', '/works/04', '/works/05'],
  },
  router: {
    middleware: ['worksMiddleware'],
  },
}
