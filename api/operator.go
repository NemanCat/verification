// ИС Верификатор
// методы API АРМ Оператора БД
package main

import (
	"context"
	"encoding/json"

	//	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//-----------------------------------------------------------------------------
//структура данных клиента при сохранении анкеты клиента
type CustomerData struct {
	//id исходной записи в таблице импорта
	Importid string `json:"importid"`
	//фамилия
	Lastname string `json:"lastname"`
	//имя
	Firstname string `json:"firstname"`
	//отчество
	Middle_initial string `json:"middle_initial"`
	//дата рождения
	Birthday string `json:"birthday"`
	//ИНН
	Inn string `json:"inn"`
	//место рождения
	Birth_place string `json:"birth_place"`
	//адрес регистрации
	Address string `json:"address"`
	//паспортные данные
	Passport string `json:"passport"`
	//номер телефона 1
	Phone1 string `json:"phone1"`
	//номер телефона 2
	Phone2 string `json:"phone2"`
	//адрес электронной почты
	Email string `json:"email"`
	//список сканов экрана
	Scans []Scan `json:"scans"`
	//Id оператора
	Operatorid string `json:"operatorid"`
	//комментарий
	Comment string `json:"comment"`
	//флаг Москва / регионы
	Is_moscow string `json:"is_moscow"`
}

//-----------------------------------------------------------------------------
// методы работы с импортированными данными
// получение на обработку новой записи из таблицы импорта
func (f *faketype) GetImportedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var userid, result string
	var useriid, recordid int64
	var exists bool

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
	sql := "LOCK TABLE customers.import IN ACCESS EXCLUSIVE MODE"
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
	sql = "SELECT EXISTS(SELECT 1 FROM customers.import WHERE status = 0) AS \"exists\""
	//sql = "SELECT EXISTS(SELECT 1 FROM customers.import WHERE status = 0 AND is_kazakh = true) AS \"exists\""
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
	sql = "SELECT id FROM customers.import WHERE status = 0 ORDER BY id LIMIT 1"
	//sql = "SELECT id FROM customers.import WHERE status = 0 AND is_kazakh = true ORDER BY id LIMIT 1"
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

	sql = "UPDATE customers.import SET status = 1, operatorid = $1, operator_block_datetime = CURRENT_TIMESTAMP WHERE id = $2"
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
	sql = "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM customers.import WHERE id = $1) t"
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

// подтверждение завершения обработки
func (f *faketype) ConfirmImportedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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
	sql := "UPDATE customers.import SET status = 2 WHERE id = $1"
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

// браковка анкеты в таблице импорта
func (f *faketype) DefectImportedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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
	comment := r.FormValue("comment")
	sql := "UPDATE customers.import SET status = 3, comment=$2, operator_block_datetime = NULL WHERE id = $1"
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
	return
}

// возврат анкеты от оператора
func (f *faketype) ReturnImportedCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
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

	sql := "UPDATE customers.import SET status = 0, operatorid = NULL, operator_block_datetime = NULL WHERE id = $1"
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

//------------------------------------------------------------------------------
// сохранение анкеты клиента
func (f *faketype) SaveCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue []byte
	)
	//получаем данные запроса клиентского приложения и проверяем их
	data := new(CustomerData)
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

	if len(data.Scans) == 0 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Нет ни одного скана паспорта клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	var scans_string = "["
	for _, value := range data.Scans {
		scans_string += "{\"img\":\"" + value.Img + "\"},"
	}
	scans_string += "]"

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

	//сохраняем новую анкету клиента в БД
	var new_id int64
	var phones []Phone
	if len(data.Phone1) > 0 {
		phones = append(phones, Phone{Phone: data.Phone1})
	}
	if len(data.Phone2) > 0 {
		phones = append(phones, Phone{Phone: data.Phone2})
	}
	var emails []Email
	if len(data.Email) > 0 {
		emails = append(emails, Email{Email: data.Email})
	}

	sql := "INSERT INTO customers.customers (lastname,firstname,middle_initial,birthday,birth_place," +
		"address,passport_data,inn,status,operatorid,screen_scans,importid,phones,emails,auditor_comment,is_moscow) " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id"
	err = Pool.QueryRow(context.Background(), sql, data.Lastname, data.Firstname, data.Middle_initial, ldate, data.Birth_place,
		data.Address, data.Passport, data.Inn, 0, data.Operatorid, data.Scans, data.Importid, phones, emails, data.Comment, data.Is_moscow).Scan(&new_id)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//помечаем запись в таблице импорта как обработанную
	iimportid, _ := strconv.ParseInt(data.Importid, 10, 64)
	sql = "UPDATE customers.import SET status = 2, operator_block_datetime = NULL WHERE id = $1"
	_, err = Pool.Exec(context.Background(), sql, iimportid)
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
			Data:    new_id,
		})
	}
	io.WriteString(w, string(retvalue))
}

//------------------------------------------------------------------------------
// отчёт по анкетам за указанный период времени
func (f *faketype) OperatorReport(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue                           []byte
		result, userid, date_from, date_to string
		useriid                            int64
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
	if date_from = r.FormValue("date_from"); date_from == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана дата начала интервала!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if date_to = r.FormValue("date_to"); date_to == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана дата окончания интервала!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//проверяем корректность дат
	ldate_from, err := time.Parse("02.01.2006", date_from)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата начала интервала указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	ldate_to, err := time.Parse("02.01.2006", date_to)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата окончания интервала указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	ldate_to = ldate_to.AddDate(0, 0, 1)
	//получаем данные отчёта
	sql := "SELECT array_to_json(array_agg(t)) FROM " +
		"(SELECT COUNT(*) as total_customers," +
		"(SELECT COUNT(*) FROM customers.customers WHERE operatorid = $1 AND (status = 1 OR status > 2) AND inserted BETWEEN $2 AND $3) as verified_customers," +
		"(SELECT COUNT(*) FROM customers.customers WHERE operatorid = $1 AND status = 2 AND inserted BETWEEN $2 AND $3) as defected_customers " +
		"FROM customers.customers WHERE operatorid = $1 AND inserted BETWEEN $2 AND $3) t"
	err = Pool.QueryRow(context.Background(), sql, useriid, ldate_from, ldate_to).Scan(&result)
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
