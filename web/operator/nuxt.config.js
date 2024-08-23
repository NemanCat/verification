export default {
  telemetry: false,  
  target: 'static',
  ssr: false,
  head: {          
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1.0' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' },
      {'http-equiv':'X-UA-Compatible', content:'IE=edge'}
    ],
    script: [     
      { src: 'https://code.jquery.com/jquery-3.5.1.slim.min.js', body: true },
      { src: 'https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js', body: true },   
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      { rel: "stylesheet", href: "https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css"},
		   { href: 'https://use.fontawesome.com/releases/v5.12.0/css/all.css',  rel: 'stylesheet' },
    ]
  }, 
  
  css: [
    {src: 'primevue/resources/primevue.min.css'},
    {src: 'primevue/resources/themes/bootstrap4-light-blue/theme.css'},    
    {src: 'primeicons/primeicons.css'},
    {src: 'primeflex/primeflex.css'},
    {src: '~assets/css/style.css'},
  ],
  
  plugins: [
    {src:'~/plugins/primevue.ts'},
    {src:'~/plugins/http.ts'},
  ],

  components: true,
  
  buildModules: [
    '@nuxt/typescript-build',
  ], 
  
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
  router: {
    base: '/operator/'
  },

  
  build: {
  }
}
