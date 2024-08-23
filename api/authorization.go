// ИС Верификатор
// функция авторизации пользователя

package main

import (
	"context"
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//------------------------------------------------------------
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//------------------------------------------------------------------------------
//функция авторизации пользователя
func (f *faketype) Authorization(w http.ResponseWriter, r *http.Request, session *Session) {
	var login, password, category, name, encrypted_password string
	var id int64
	var icategory int32
	//проверяем наличие в теле клиентского запроса параметров login, password и category
	r.ParseForm()
	if login = r.FormValue("login"); login == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if password = r.FormValue("password"); password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if category = r.FormValue("category"); category == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//ищем в БД пользователя с указанным логином
	sql := "SELECT id,name,password,category FROM settings.users WHERE login=$1 AND is_blocked = false"
	err := Pool.QueryRow(context.Background(), sql, login).Scan(&id, &name, &encrypted_password, &icategory)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	authorized := CheckPasswordHash(password, encrypted_password)
	if authorized == false {
		message := "Попытка доступа к системе с невалидными реквизитами доступа, логин  " + login + " пароль " + password + " ip-адрес " + GetIP(r)
		go logging_service.WriteLogMessage("security", message)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//проверяем категорию пользователя кроме администраторов системы
	//администраторы имеют доступ ко всем АРМ
	if icategory > 0 {
		if int32(StrToInt(category)) != icategory {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	//пользователь авторизовался, создаём для него сессию
	//ID записи
	new_id := SessionsList.CreateId()
	//уникальный ID сессии
	guid := uuid.New()
	//новая сессия
	new_session := NewSession(new_id, guid.String(), time.Now(), icategory, id, name, login, GetIP(r))
	//подпись cookie
	hash := md5.Sum([]byte(cookie_secret + guid.String()))
	signature := hex.EncodeToString(hash[:])
	//CSFR security token
	hash = md5.Sum([]byte(csfr_secret + guid.String()))
	csfr_token := hex.EncodeToString(hash[:])
	//добавляем новую сессию в список сессий
	err = SessionsList.AddObject(guid.String(), new_session)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		if log_level == "1" {
			message := "Пользователь " + new_session.Name + " (" + login + ") вошёл в систему с ip-адреса " + new_session.Ip +
				",id сессии " + new_session.Sessionid
			go logging_service.WriteLogMessage("authorization", message)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     session_signature_cookie_name,
			Value:    string(b64.StdEncoding.EncodeToString([]byte(signature))),
			HttpOnly: true,
			//TO-DO добавить Domain и Path
		})
		w.WriteHeader(http.StatusOK)
		retvalue, _ := json.Marshal(struct {
			Id         int64  `json:"id"`
			Name       string `json:"name"`
			Csfr_token string `json:"csfr_token"`
			Sessionid  string `json:"sessionid"`
		}{
			id,
			name,
			csfr_token,
			string(b64.StdEncoding.EncodeToString([]byte(guid.String()))),
		})
		io.WriteString(w, string(retvalue))
	}
}

//функция штатного выхода из системы
func (f *faketype) Logout(w http.ResponseWriter, r *http.Request, session *Session) {
	if log_level == "1" {
		message := "Пользователь " + session.Name + " (" + session.Login + "), id сессии " + session.Sessionid + " вышел из системы"
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
	http.SetCookie(w, &http.Cookie{
		Name:    session_signature_cookie_name,
		Expires: expire,
	})
	w.WriteHeader(http.StatusOK)
}

//функция принудительного выхода из системы по таймауту на клиенте
func (f *faketype) SystemLogout(w http.ResponseWriter, r *http.Request, session *Session) {
	if log_level == "1" {
		message := "Пользователь " + session.Name + " (" + session.Login + "), id сессии " + session.Sessionid + " выведен из системы по тайм-ауту на клиенте"
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
	http.SetCookie(w, &http.Cookie{
		Name:    session_signature_cookie_name,
		Expires: expire,
	})
	w.WriteHeader(http.StatusOK)
}
