// ИС Верификатор
// методы API АРМ Суперпользователя
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

//------------------------------------------------------------------------------
//получение списка анкет клиентов
func (f *faketype) GetCustomersList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
		"(" +
		"SELECT A.id,A.lastname,A.firstname,A.middle_initial,A.status " +
		"FROM customers.customers A " +
		"ORDER BY A.lastname, A.firstname, A.middle_initial" +
		") t"
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

//получение списка забракованных анкет из таблиц импорта
func (f *faketype) GetDefectedCustomersList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
		"(" +
		"SELECT * FROM " +
		"(SELECT A.id,A.lastname,A.firstname,A.middle_initial,A.status,B.name as operator, '' AS auditor, '' AS pasportist, 0 AS source " +
		"FROM customers.import A JOIN settings.users B ON A.operatorid = B.id WHERE A.status = 3 " +
		"UNION " +
		"SELECT A.id,A.lastname,A.firstname,A.middle_initial,A.status,B.name AS operator,C.name AS auditor,D.name AS pasportist, 1 AS source " +
		"FROM customers.customers A JOIN settings.users B ON A.operatorid = B.id LEFT JOIN settings.users C ON A.auditorid = C.id " +
		"LEFT JOIN settings.users D ON A.pasportistid = D.id " +
		"WHERE A.status = 2 OR A.status = 5) AS p " +
		"ORDER BY lastname, firstname, middle_initial " +
		") t"
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

//получение анкеты клиента для просмотра/редактирования
func (f *faketype) GetSuperuserCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue   []byte
		id, result string
		iid        int64
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//выбираем данные записи
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
		"(SELECT A.*, B.name AS operator, C.name AS auditor, D.name AS pasportist, E.name AS telephonist " +
		" FROM customers.customers A JOIN settings.users B on A.operatorid = B.id " +
		" LEFT JOIN settings.users C ON A.auditorid = C.id " +
		" LEFT JOIN settings.users D ON A.pasportistid = D.id " +
		" LEFT JOIN settings.users E ON A.telephonistid = E.id " +
		" WHERE A.id=$1) t"
	err = Pool.QueryRow(context.Background(), sql, iid).Scan(&result)
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

//получение забракованной анкеты клиента для просмотра/редактирования
func (f *faketype) GetSuperuserDefectedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue           []byte
		id, source, result string
		iid                int64
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if source = r.FormValue("source"); source == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан источник записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//выбираем данные записи
	//если source = 0, выбираем данные из таблицы импорта
	//если source - 1, выбираем данные из рабочей таблицы
	if source == "0" {
		sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
			"(SELECT lastname,firstname,middle_initial,address,comment  " +
			" FROM customers.import WHERE id=$1) t"
		err = Pool.QueryRow(context.Background(), sql, iid).Scan(&result)
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
	} else {
		sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
			"(SELECT lastname, firstname, middle_initial,address,auditor_comment as comment " +
			" FROM customers.customers WHERE id=$1) t"
		err = Pool.QueryRow(context.Background(), sql, iid).Scan(&result)
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
	}

	io.WriteString(w, string(retvalue))
}

//сохранение комментария в анкете клиента
func (f *faketype) SaveComment(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue    []byte
		id, comment string
		iid         int64
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}

	if comment = r.FormValue("comment"); comment == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан текст комментария!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	sql := "UPDATE customers.customers SET auditor_comment = $2 WHERE id = $1"
	_, err = Pool.Exec(context.Background(), sql, iid, comment)
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
			Data:    id,
		})
	}
	io.WriteString(w, string(retvalue))
}

//возвращение забракованной анкеты в обработку
func (f *faketype) ReturnDefectedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue                []byte
		id, source, status, sql string
		iid                     int64
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if source = r.FormValue("source"); source == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан источник записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if status = r.FormValue("status"); status == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан статус!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if source == "0" {
		//анкета была забракована оператором
		sql = "UPDATE customers.import SET status = 0, operatorid = NULL, operator_block_datetime = NULL WHERE id = $1"
		_, err = Pool.Exec(context.Background(), sql, iid)
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
				Data:    id,
			})
		}
		io.WriteString(w, string(retvalue))
	} else {
		if status == "2" {
			//анкета была забракована аудитором
			sql = "UPDATE customers.import SET status = 0, operatorid = NULL, operator_block_datetime = NULL " +
				"WHERE id IN (SELECT importid FROM customers.customers WHERE id=$1)"
			_, err = Pool.Exec(context.Background(), sql, iid)
			if err != nil {
				retvalue, _ = json.Marshal(Retvalue{
					Success: false,
					Message: "Database access error: " + err.Error(),
					Data:    nil,
				})
				io.WriteString(w, string(retvalue))
				return
			}
			sql = "DELETE FROM customers.customers WHERE id=$1"
			_, err = Pool.Exec(context.Background(), sql, iid)
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
					Data:    id,
				})
			}
			io.WriteString(w, string(retvalue))
		} else {
			//анкета была забракована паспортистом
			sql = "UPDATE customers.customers SET status = 1, pasportistid = NULL, pasportist_block_datetime = NULL WHERE id = $1"
			_, err = Pool.Exec(context.Background(), sql, iid)
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
					Data:    id,
				})
			}
			io.WriteString(w, string(retvalue))
		}
	}
}
