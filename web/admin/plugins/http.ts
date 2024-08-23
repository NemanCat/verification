import { Plugin } from '@nuxt/types'

const httpPlugin: Plugin = function ({ store,$http }: any) {
        $http.setHeader('Content-Type', 'application/x-www-form-urlencoded');
        if (window.sessionStorage.getItem('csfr_security_token')) {
            $http.setHeader('x-request-security-token', window.sessionStorage.getItem('csfr_security_token'));
        }

        $http.onError((error: { statusCode: number; }) => {  
            let message: string = '';
            switch (error.statusCode) {
                case 504: 
                    message = 'Сервер выключен или недоступен, попробуйте повторить операцию позже';
                    break;
                case 403:
                    message = 'Ваш логин и/или пароль не прошли проверку, попробуйте ещё раз';
                    break;
                case 404:
                    message = 'Запрошенный ресурс отсутствует на сервере';
                    break;    
                default:
                    message = 'Ошибка при получении данных с сервера.';        
            }
            store.commit('CHANGE_FETCHERROR', {
                fetchErrorFlag: true,
                fetchErrorCode: error.statusCode,
                fetchErrorMessage: message
            });
        });    
}

export default httpPlugin