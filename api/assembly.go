// ИС Верификатор
// методы API АРМ Телефониста
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
//получение списка анкет, находящихся на обработке у данного телефониста
func (f *faketype) GetAssemblyCustomersList(w http.ResponseWriter, r *http.Request, session *Session) {
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
		"SELECT id, lastname, firstname, middle_initial, phones, '' as phones_string " +
		"FROM customers.customers " +
		"WHERE status = 6 AND telephonistid = $1 " +
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

//получение анкеты клиента для просмотра/редактирования
func (f *faketype) GetAssemblyCustomerData(w http.ResponseWriter, r *http.Request, session *Session) {
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

//------------------------------------------------------------------------------
//получение анкеты клиента на прозвон
func (f *faketype) GetCustomerForCalling(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue          []byte
		userid, result    string
		useriid, recordid int64
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
	sql = "SELECT EXISTS(SELECT 1 FROM customers.customers WHERE status = 4) AS \"exists\""
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
			Message: "Нет ни одной новой анкеты клиента для обработки!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//выбираем первую подходящую для обработки запись и помечаем её как находящуюся в обработке
	sql = "SELECT id FROM customers.customers WHERE status = 4 ORDER BY id LIMIT 1"
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
	sql = "UPDATE customers.customers SET status = 6, telephonistid = $1, telephonist_block_datetime = CURRENT_TIMESTAMP WHERE id = $2"
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
	sql = "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT id,lastname,firstname,middle_initial,phones,'' as phones_string FROM customers.customers WHERE id = $1) t"
	err = tx.QueryRow(context.Background(), sql, recordid).Scan(&result)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
	} else {
		//заканчиваем транзакцию
		tx.Commit(context.Background())
		retvalue, _ = json.Marshal(Retvalue{
			Success: true,
			Message: "",
			Data:    result,
		})
	}
	io.WriteString(w, string(retvalue))
}

//возврат анкеты от телефониста
/*func (f *faketype) GetAnotherCustomerForCalling(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, userid             string
		iid, useriid, recordid int64
		err                    error
		retvalue               []byte
		exists                 bool
		result                 string
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
	//возвращаем ранее загруженного пользователя в общую базу
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
	sql := "UPDATE customers.customers SET status = 4, telephonistid = NULL, telephonist_block_datetime = NULL WHERE id = $1"
	_, err = tx.Exec(context.Background(), sql, iid)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	err = tx.Commit(context.Background())
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//выбираем другого пользователя
	//проверяем имеются ли доступные для обработки записи
	tx, err = Pool.Begin(context.Background())
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
	sql = "SELECT EXISTS(SELECT 1 FROM customers.customers WHERE status = 4 AND phones IS NOT NULL AND id <> $1) AS \"exists\""
	err = tx.QueryRow(context.Background(), sql, iid).Scan(&exists)
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
			Message: "Нет ни одной новой анкеты клиента для обработки!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//выбираем первую подходящую для обработки запись и помечаем её как находящуюся в обработке
	sql = "SELECT id FROM customers.customers WHERE status = 4 AND phones IS NOT NULL AND id <> $1 ORDER BY id LIMIT 1"
	err = tx.QueryRow(context.Background(), sql, iid).Scan(&recordid)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	sql = "UPDATE customers.customers SET status = 6, telephonistid = $1, telephonist_block_datetime = CURRENT_TIMESTAMP WHERE id = $2"
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
	sql = "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT id,lastname,firstname,middle_initial,inn,birthday,birth_place,address,passport_data,phones,emails,passport_scans FROM customers.customers WHERE id = $1) t"
	err = tx.QueryRow(context.Background(), sql, recordid).Scan(&result)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
	} else {
		//заканчиваем транзакцию
		tx.Commit(context.Background())
		retvalue, _ = json.Marshal(Retvalue{
			Success: true,
			Message: "",
			Data:    result,
		})
	}
	io.WriteString(w, string(retvalue))
}*/

//сохранение данных клиента
func (f *faketype) SaveCalledCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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
	sql := "UPDATE customers.customers SET called=CURRENT_TIMESTAMP,lastname=$1,firstname=$2," +
		"middle_initial=$3,birthday=$4,birth_place=$5,address=$6,passport_data=$7,phones=$8,emails=$9,auditor_comment=$10,inn=$11 WHERE id=$12"
	_, err = Pool.Exec(context.Background(), sql, data.Lastname, data.Firstname, data.Middle_initial, ldate,
		data.Birth_place, data.Address, data.Passport, data.Phones, data.Emails, data.Comment, data.Inn, iid)
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
func (f *faketype) DefectAssemblyCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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
	sql := "UPDATE customers.customers SET status=7,called=CURRENT_TIMESTAMP WHERE id=$1"
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

//клиент готов к дальнейшему сотрудничеству
func (f *faketype) ConfirmAssemblyCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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
	sql := "UPDATE customers.customers SET status=9,called=CURRENT_TIMESTAMP WHERE id=$1"
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

//клиент не готов к дальнейшему сотрудничеству
func (f *faketype) RejectAssemblyCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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
	sql := "UPDATE customers.customers SET status=8,called=CURRENT_TIMESTAMP WHERE id=$1"
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
