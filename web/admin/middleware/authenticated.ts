import { Middleware } from '@nuxt/types'

const AuthMiddleware: Middleware = async (context) => {
  //if (process.server) {   
      if(!context.app.$cookies.get("verification-session")) {
        context.redirect('/login');
      }    
  //}
}

export default AuthMiddleware;