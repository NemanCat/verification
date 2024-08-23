// сервис подсистемы безопасности
// сервер сессий

package main

import (
	"encoding/json"
	"reflect"
	"time"
)

//-------------------------------------------------
//клиентская сессия
type Session struct {
	BaseObject
	//id сессии - является также ключём в списке сессий
	Sessionid string
	//дата и время последней активности пользователя
	Last_activity time.Time
	//категория клиента
	//0 - администратор admin
	//1 - архивариус archivarius
	//2 - оператор БД operator
	//3 - аудитор auditor
	//4 - телефонист telephonist
	//5 - паспортист pasportist
	//6 - суперпользователь superuser
	Category int32
	//id пользователя
	Userid int64
	//ФИО пользователя
	Name string
	//логин пользователя
	Login string
	//ip-адрес клиента
	Ip string
}

//---------------------------------------------
//методы работы с информационным объектом
//конструктор объекта
func NewSession(id uint64, sessionid string, last_activity time.Time, category int32, userid int64, name string, login string, ip string) *Session {
	session := new(Session)
	session.Id = id
	session.Sessionid = sessionid
	session.Last_activity = last_activity
	session.Category = category
	session.Userid = userid
	session.Name = name
	session.Login = login
	session.Ip = ip
	return session
}

//обновление поля last_activity
func (session *Session) SetLastActivity(last_activity time.Time) {
	session.Last_activity = last_activity
}

//обновление поля ip
func (session *Session) SetLastIp(last_ip string) {
	session.Ip = last_ip
}

//декодирование объекта
func (session *Session) UnmarshalJSON(encoded []byte) string {
	err := json.Unmarshal(encoded, session)
	if err != nil {
		return err.Error()
	}
	return ""
}

//---------------------------------------------------------------
//методы сохранения информационного объекта в БД Bolt
//получение списка сеcсий
func (provider *BoltDbProvider) GetSessionsList() (map[string]BaseObjectInterface, error) {
	return provider.List("verification.db", "verification.db", "sessions", reflect.TypeOf(Session{}))
}

//сохранение сессии в БД
func (provider *BoltDbProvider) PersistSession(key string, value BaseObjectInterface) error {
	return provider.Persist("verification.db", "verification.db", "sessions", key, value)
}

//удаление сессии из БД
func (provider *BoltDbProvider) DeleteSession(key string) error {
	return provider.Delete("verification.db", "verification.db", "sessions", key)
}
