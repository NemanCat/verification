// ИС Верификатор
// методы API АРМ Администратора
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

//------------------------------------------------------------------------------
//функции работы со списком пользователей

//получение списка пользователей
func (f *faketype) GetUsersList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.users ORDER BY name) t"
	err := Pool.QueryRow(context.Background(), sql).Scan(&result)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
	} else {
		retvalue, _ = json.Marshal(Retvalue{
			Success: true,
			Message: "",
			Data:    result,
		})
	}
	io.WriteString(w, string(retvalue))
}

//добавление/редактирование аккаунта пользователя
func (f *faketype) EditUser(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var iid, icategory int64
	var err error
	var id, category, name, login, password, password_changed string
	var exists bool

	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id пользователя должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if category = r.FormValue("category"); category == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана категория пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		icategory, err = strconv.ParseInt(category, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Категория пользователя должна быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}

	}
	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указаны ФИО пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if login = r.FormValue("login"); login == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан логин пользователя!",
			Data:    nil,
		})
		return
	}
	if password = r.FormValue("password"); password == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан пароль пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//шифруем пароль пользователя
	hash, _ := HashPassword(password)
	if iid == 0 {
		//новый аккаунт пользователя
		//проверяем наличие в БД аккаунта пользователя с таким логином
		sql := "SELECT EXISTS(SELECT 1 FROM settings.users WHERE login=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, login).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			if exists {
				retvalue, _ = json.Marshal(Retvalue{
					Success: false,
					Message: "Аккаунт пользователя с таким логином уже есть в базе данных!",
					Data:    nil,
				})
			} else {
				var new_id int64
				sql = "INSERT INTO settings.users (name,category,login,password,is_blocked) VALUES ($1,$2,$3,$4,false) RETURNING id"
				err = Pool.QueryRow(context.Background(), sql, name, icategory, login, hash).Scan(&new_id)
				if err != nil {
					retvalue, _ = json.Marshal(Retvalue{
						Success: false,
						Message: err.Error(),
						Data:    nil,
					})
				} else {
					retvalue, _ = json.Marshal(Retvalue{
						Success: true,
						Message: "",
						Data:    new_id,
					})
				}
			}
		}
	} else {
		//изменение данных существующего аккаунта пользователя
		//проверяем наличие в БД аккаунта пользователя с таким логином
		sql := "SELECT EXISTS(SELECT 1 FROM settings.users WHERE login=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, login, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			if exists {
				retvalue, _ = json.Marshal(Retvalue{
					Success: false,
					Message: "Аккаунт пользователя с таким логином уже есть в базе данных!",
					Data:    nil,
				})
			} else {
				//список изменяемых полей зависит от того, изменился ли пароль аккаунта
				if password_changed = r.FormValue("password_changed"); password_changed == "" {
					retvalue, _ = json.Marshal(Retvalue{
						Success: false,
						Message: "Не указан флаг редактирования пароля пользователя!",
						Data:    nil,
					})
				} else {
					if password_changed == "true" {
						//пароль изменился
						sql := "UPDATE settings.users SET name = $1, login = $2, password = $3, category = $4 WHERE id = $5"
						_, err = Pool.Exec(context.Background(), sql, name, login, hash, icategory, iid)
					} else {
						sql := "UPDATE settings.users SET name = $1, login = $2, category = $3 WHERE id = $4"
						_, err = Pool.Exec(context.Background(), sql, name, login, icategory, iid)
					}
					if err != nil {
						retvalue, _ = json.Marshal(Retvalue{
							Success: false,
							Message: err.Error(),
							Data:    nil,
						})
					} else {
						retvalue, _ = json.Marshal(Retvalue{
							Success: true,
							Message: "",
							Data:    iid,
						})
					}
				}
			}
		}

	}
	io.WriteString(w, string(retvalue))
}

//удаление аккаунта пользователя
func (f *faketype) DeleteUser(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		err      error
		retvalue []byte
		iid      int64
		id       string
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id пользователя должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id пользователя должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}

	sql := "DELETE FROM settings.users WHERE id=$1"
	_, err = Pool.Exec(context.Background(), sql, iid)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	} else {
		retvalue, _ = json.Marshal(Retvalue{
			Success: true,
			Message: "",
			Data:    iid,
		})
	}
	io.WriteString(w, string(retvalue))
}

//------------------------------------------------------------------------------
//получение списка открытых пользовательских сессий
func (f *faketype) GetSessionsList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	retvalue, _ = json.Marshal(Retvalue{
		Success: true,
		Message: "",
		Data:    SessionsList.GetList(),
	})
	io.WriteString(w, string(retvalue))
}
