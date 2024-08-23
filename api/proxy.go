// ИС Верификатор
// сервис единой точки доступа к API
// проксирование запросов в микросервисы
package main

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"time"
)

//---------------------------------------------------
type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(r, false)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	return http.DefaultTransport.RoundTrip(r)
}

//---------------------------------------------------
type SecureProxy struct {
}

//---------------------------------------------------
//проверка права доступа запроса к запрошенному ресурсу
//возвращаемые значения
//*Session сессия
//bool пройдена ли проверка
//error возникшая процессе проверки ошибка
func checkAccess(route Route, r *http.Request) (*Session, bool, error) {
	//проверяем наличие идентификационного cookie
	cookie, err := r.Cookie(session_cookie_name)
	if err != nil {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " без cookie сессии, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, err
	}
	if cookie == nil {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " с cookie сессии = null, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, nil
	}
	//проверяем наличие cookie подписи
	signature_cookie, err := r.Cookie(session_signature_cookie_name)
	if err != nil {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " без cookie подписи, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, err
	}
	if cookie == nil {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " с cookie подписи = null, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, nil
	}
	//id сессии и подпись
	credentials := struct {
		Sessionid string `json:"sessionid"`
		Signature string `json:"signature"`
	}{}
	std, _ := b64.StdEncoding.DecodeString(cookie.Value)
	credentials.Sessionid = string(std)
	std, _ = b64.StdEncoding.DecodeString(signature_cookie.Value)
	credentials.Signature = string(std)
	//проверяем правильность подписи cookie
	hash := md5.Sum([]byte(cookie_secret + credentials.Sessionid))
	signature := hex.EncodeToString(hash[:])
	if signature != credentials.Signature {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " с невалидной подписью cookie, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, nil
	}
	//проверяем наличие заголовка с csfr security token
	csfr_token_header := r.Header.Get("x-request-security-token")
	if len(csfr_token_header) == 0 {
		//заголовок csfr security token отсутствует в запросе
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " без заголовка с csfr security token, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, nil
	}
	//проверяем соответствие session id и csfr security token
	hash = md5.Sum([]byte(csfr_secret + credentials.Sessionid))
	signature = hex.EncodeToString(hash[:])
	if signature != csfr_token_header {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " с несоответствующими session id и csfr security token, ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, nil
	}
	//ищем сессию пользователя в списке сессий
	session := SessionsList.FindObject(credentials.Sessionid)
	if session == nil {
		if log_level == "1" {
			message := "Попытка доступа к маршруту " + r.URL.EscapedPath() + " несуществующей сессией " + credentials.Sessionid + ", ip-адрес " + GetIP(r)
			go logging_service.WriteLogMessage("security", message)
		}
		return nil, false, nil
	}
	//проверяем категорию пользователя
	//маршруты, доступные всем пользователям, имеют category = -1
	//администраторы имеют доступ ко всем АРМам
	required_category := route.Category
	if required_category == -1 {
		return session.(*Session), true, nil
	} else {
		user_category := session.(*Session).Category
		if (required_category != user_category) && (user_category != 0) {
			if log_level == "1" {
				message := "Попытка несанкционированного доступа к маршруту " + r.URL.EscapedPath() + " пользователя " + session.(*Session).Name +
					" (" + session.(*Session).Login + "), ip-адрес " + GetIP(r)
				go logging_service.WriteLogMessage("security", message)
			}
			return nil, false, nil
		}
	}
	return session.(*Session), true, nil
}

//---------------------------------------------------------------
//функция безопасной обработки запросов к API
func (p *SecureProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var f faketype
	var rt Route
	var ok bool
	//принимаем только POST-запросы
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//ищем запрошенный клиентом endpoint в списке маршрутов
	//если запрошенный endpoint не найден - возвращаем 404
	if rt, ok = ApiRoutes[r.URL.EscapedPath()]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var session *Session = nil
	var has_access bool
	//запрос требует авторизации?
	if rt.Is_protected {
		//проверяем права доступа
		session, has_access, _ = checkAccess(rt, r)
		if !has_access {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//проверяем дату и время последней активности пользователя
		diff := time.Now().Sub(session.Last_activity)
		if diff.Minutes() > float64(StrToInt(max_session_timeout)) {
			//превышен максимальный интервал ожидания, разлогиниваем пользователя
			if log_level == "1" {
				message := "Пользователь " + session.Name +
					" (" + session.Login + "), id сессии " + session.Sessionid + ", ip-адрес " + GetIP(r) +
					" принудительно выведен из системы при попытке доступа к маршруту " +
					r.URL.EscapedPath() + " вследствие длительной неактивности, время последней активности " +
					session.Last_activity.String()
				go logging_service.WriteLogMessage("authorization", message)
			}
			//удаляем сессию
			SessionsList.DeleteObject(session.Sessionid)
			//удаляем cookie
			expire := time.Now().Add(-7 * 24 * time.Hour)
			http.SetCookie(w, &http.Cookie{
				Name:    session_cookie_name,
				Expires: expire,
			})
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//устанавливаем дату и время последней активности пользователя
		go session.SetLastActivity(time.Now())
		//если изменился ip-адрес подключения, обновляем ip-адрес
		if r.RemoteAddr != session.Ip {
			go session.SetLastIp(GetIP(r))
		}
	}
	//обрабатываем запрос
	fn := reflect.ValueOf(&f).MethodByName(rt.Handler)
	if !fn.IsValid() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		inputs := make([]reflect.Value, 3)
		inputs[0] = reflect.ValueOf(w)
		inputs[1] = reflect.ValueOf(r)
		inputs[2] = reflect.ValueOf(session)
		fn.Call(inputs)
	}
}
