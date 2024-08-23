// ИС Верификатор
// методы API АРМ Паспортиста
package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//------------------------------------------------------------------------------
//получение списка анкет, находящихся на обработке у данного паспортиста
func (f *faketype) GetPasportistCustomersList(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		userid, result string
		useriid        int64
		retvalue       []byte
		err            error
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if userid = r.FormValue("userid"); userid == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		useriid, err = strconv.ParseInt(userid, 10, 64)
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

	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
		"(" +
		"SELECT id, lastname, firstname, middle_initial " +
		"FROM customers.customers " +
		"WHERE (status = 3) AND (pasportistid = $1) AND (importid not in " +
		"(SELECT importid FROM customers.customers GROUP BY importid HAVING COUNT(importid) > 1))" +
		"ORDER BY lastname, firstname, middle_initial" +
		") t"
	err = Pool.QueryRow(context.Background(), sql, useriid).Scan(&result)
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

//получение новой верифицированной анкеты клиента на обработку
func (f *faketype) GetVeirifedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		userid, result    string
		useriid, recordid int64
		retvalue          []byte
		err               error
		exists            bool
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if userid = r.FormValue("userid"); userid == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id пользователя!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		useriid, err = strconv.ParseInt(userid, 10, 64)
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
	//стартуем транзакцию
	tx, err := Pool.Begin(context.Background())
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	defer tx.Rollback(context.Background())
	//блокируем таблицу
	sql := "LOCK TABLE customers.customers IN ACCESS EXCLUSIVE MODE"
	_, err = tx.Exec(context.Background(), sql)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	//проверяем имеются ли доступные для обработки записи
	if useriid == 48 {
		sql = "SELECT EXISTS(SELECT 1 FROM customers.customers WHERE (status = 1) AND (is_kazakh = true) AND (importid not in " +
			"(SELECT importid FROM customers.customers GROUP BY importid HAVING COUNT(importid) > 1))) AS \"exists\""
	} else {
		sql = "SELECT EXISTS(SELECT 1 FROM customers.customers WHERE (status = 1) AND (is_kazakh = false OR is_kazakh is NULL) AND (importid not in " +
			"(SELECT importid FROM customers.customers GROUP BY importid HAVING COUNT(importid) > 1))) AS \"exists\""
	}
	err = tx.QueryRow(context.Background(), sql).Scan(&exists)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if !exists {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Level:   1,
			Message: "Нет ни одной анкеты клиента для обработки!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	//выбираем первую подходящую для обработки запись и помечаем её как находящуюся в обработке
	if useriid == 48 {
		sql = "SELECT id FROM customers.customers WHERE (status = 1) AND (is_kazakh = true) ORDER BY id LIMIT 1"
	} else {
		sql = "SELECT id FROM customers.customers WHERE (status = 1) AND (is_kazakh = false OR is_kazakh is NULL) ORDER BY id LIMIT 1"
	}

	err = tx.QueryRow(context.Background(), sql).Scan(&recordid)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	sql = "UPDATE customers.customers SET status = 3, pasportistid = $1, pasportist_block_datetime = CURRENT_TIMESTAMP WHERE id = $2"
	_, err = tx.Exec(context.Background(), sql, useriid, recordid)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//выбираем все данные выбранной записи
	sql = "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT id, lastname, firstname, middle_initial FROM customers.customers WHERE id = $1) t"
	err = tx.QueryRow(context.Background(), sql, recordid).Scan(&result)
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
		//заканчиваем транзакцию
		tx.Commit(context.Background())
	}
	io.WriteString(w, string(retvalue))
}

//получение анкеты клиента для просмотра/редактирования
func (f *faketype) GetPassportCustomerData(w http.ResponseWriter, r *http.Request, session *Session) {
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
		"(SELECT id,lastname,firstname,middle_initial,birthday,birth_place,address,passport_data,inn,phones,emails,operatorid,screen_scans,passport_scans,auditor_comment as comment " +
		"FROM customers.customers WHERE id=$1) t"
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

//возврат анкеты в общую базу
func (f *faketype) ReturnPassportCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id       string
		iid      int64
		err      error
		retvalue []byte
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

	sql := "UPDATE customers.customers SET status = 1, pasportistid = NULL, pasportist_block_datetime = NULL WHERE id = $1"
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
	return
}

//сохранение данных клиента
func (f *faketype) SavePassportCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue []byte
		iid      int64
	)

	//получаем данные запроса клиентского приложения и проверяем их
	data := new(AssemblyCustomerData)
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err := json.Unmarshal(b, data)

	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if len(data.Id) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id записи!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(data.Id, 10, 64)
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

	if len(data.Lastname) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана фамилия клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	if len(data.Firstname) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано имя клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if len(data.Birthday) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана дата рождения клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if len(data.Birth_place) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано место рождения клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if len(data.Passport) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указаны паспортные данные клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	//проверяем корректность даты рождения
	ldate, err := time.Parse("02.01.2006", data.Birthday)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата рождения указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	//сохраняем данные клиента в БД
	sql := "UPDATE customers.customers SET lastname=$1,firstname=$2," +
		"middle_initial=$3,birthday=$4,birth_place=$5,address=$6,passport_data=$7,phones=$8,emails=$9,inn=$10,auditor_comment=$11,passport_scans=$12 WHERE id=$13"
	_, err = Pool.Exec(context.Background(), sql, data.Lastname, data.Firstname, data.Middle_initial, ldate,
		data.Birth_place, data.Address, data.Passport, data.Phones, data.Emails, data.Inn, data.Comment, data.Scans, iid)
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
			Data:    data.Id,
		})
	}
	io.WriteString(w, string(retvalue))
	return
}

//браковка анкеты клиента
func (f *faketype) DefectPassportCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id       string
		iid      int64
		retvalue []byte
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

	//бракуем анкету
	sql := "UPDATE customers.customers SET status=5,pasported=CURRENT_TIMESTAMP WHERE id=$1"
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
	return
}

//одобрение анкеты клиента
func (f *faketype) ConfirmPassportCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id       string
		iid      int64
		retvalue []byte
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

	//одобряем анкету
	sql := "UPDATE customers.customers SET status=4,pasported=CURRENT_TIMESTAMP WHERE id=$1"
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
	return
}
