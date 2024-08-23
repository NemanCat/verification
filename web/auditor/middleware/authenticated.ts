import { Middleware } from '@nuxt/types'

const AuthMiddleware: Middleware = async (context) => {  
  //проверяем наличие cookie сессии
  if(!context.app.$cookies.get("verification-session")) {
        context.redirect('/login');
  }    
  //проверяем наличие csfr security token
  if((window.sessionStorage.getItem("csfr_security_token") == null) || (window.sessionStorage.getItem("csfr_security_token")?.length == 0)) {    
    context.redirect('/login');
  }
}

export default AuthMiddleware;