
export default { 
  /*
  ** Nuxt rendering mode
  ** See https://nuxtjs.org/api/configuration-mode
  */
  //mode: 'universal',
  /*
  ** Nuxt target
  ** See https://nuxtjs.org/api/configuration-target
  */
  target: 'static',
  ssr: false,
  /*
  ** Headers of the page
  ** See https://nuxtjs.org/api/configuration-head
  */
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' },
      {'http-equiv':'X-UA-Compatible', content:'IE=edge'}
    ],
    script: [     
      { src: 'https://code.jquery.com/jquery-3.3.1.slim.min.js' },
      { src: 'https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/js/bootstrap.min.js' },   
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
       /* Bootstrap v4.5.0 CSS via CDN */
       { rel: "stylesheet", href:"https://maxcdn.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css" },
        /* Font Awesome */
		   { href: 'https://use.fontawesome.com/releases/v5.12.0/css/all.css',  rel: 'stylesheet' },
    ]
  },
  /*
  ** Global CSS
  */
  css: [
    {src: 'primevue/resources/primevue.min.css'},
    {src: 'primevue/resources/themes/bootstrap4-light-blue/theme.css'},    
    {src: 'primeicons/primeicons.css'},
    {src: 'primeflex/primeflex.css'},
    {src: '~assets/css/style.css'},
  ],
  /*
  ** Plugins to load before mounting the App
  ** https://nuxtjs.org/guide/plugins
  */
  plugins: [
    {src:'~/plugins/primevue.ts',mode:'client'},
    {src:'~/plugins/loading.ts',mode:'client'},
    {src:'~/plugins/http.ts',mode:'client'},
  ],
  /*
  ** Auto import components
  ** See https://nuxtjs.org/api/configuration-components
  */
  components: true,
  /*
  ** Nuxt.js dev-modules
  */
  buildModules: [
    '@nuxt/typescript-build',
  ],
  /*
  ** Nuxt.js modules
  */
  modules: [
    '@nuxt/http',
    'cookie-universal-nuxt',
  ],
  http: {    
    proxy: true,
    serverTimeout: 5000,
    clientTimeout: 5000
  },
  proxy: {
    '/api/': {
      target: process.env.API_URL
    }
  },
  /*
  ** Build configuration
  ** See https://nuxtjs.org/api/configuration-build/
  */
  build: {
  },  
  router: {
    base: '/admin/'
  }
}
