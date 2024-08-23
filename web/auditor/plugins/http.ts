import { Plugin } from '@nuxt/types'

const httpPlugin: Plugin = function ({ store,$http,app }: any) {

        $http.onRequest((config: any) => {
            config.headers.set('Content-Type', 'application/x-www-form-urlencoded');
            //для запроса авторизации session cookie и заголовок x-request-security-token не нужен
            //для остальных запросов нужен            
            if(config.url.indexOf("authorization") == -1) {
                if(!app.$cookies.get("verification-session")) {
                    return new Response('Auth',{status:403})
                }                                    
                if (window.sessionStorage.getItem('csfr_security_token')) {                    
                    config.headers.set('x-request-security-token', window.sessionStorage.getItem('csfr_security_token'));                    
                } else {                    
                    //если отсутствует csfr security token, отправляем на страницу авторизации
                    return new Response('Auth',{status:403})
                }
            }    
        });

        $http.onError((error: any) => {  
            let message: string = '';
            
            switch (error.statusCode) {
                case 504: 
                    //gateway timeout
                    message = 'Сервер выключен или недоступен, попробуйте повторить операцию позже';                    
                    break;
                case 502: 
                    //bad gateway
                    message = 'Сервер выключен или недоступен, попробуйте повторить операцию позже';                    
                    break;    
                case 403:
                    message = 'Ваш логин и/или пароль не прошли проверку, попробуйте ещё раз. ';
                    break;
                case 404:
                    message = 'Запрошенный ресурс отсутствует на сервере.';
                    break;    
                default:
                    message = 'При получении данных с сервера произошла ошибка ' + error.statusCode;       
            }
            store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: true,
                fetchErrorCode: error.statusCode,
                fetchErrorMessage: message,
                fetchErrorLevel: 0,
            });
        });    
}

export default httpPlugin