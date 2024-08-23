// ИС Верификатор
// методы API АРМ Аудитора
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgtype"
)

//------------------------------------------------------------------------------
//получение списка неподтверждённых анкет клиентов
func (f *faketype) GetNewCustomersList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
		"(" +
		"SELECT A.id,A.lastname,A.firstname,A.middle_initial, B.name, B.login, A.inserted, A.importid, " +
		"EXISTS(SELECT 1 FROM customers.customers WHERE importid = A.importid and id <> A.id) " +
		"FROM customers.customers A JOIN settings.users B ON A.operatorid = B.id " +
		"WHERE A.status = 0 " +
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

//получение анкеты клиента для просмотра/редактирования
func (f *faketype) GetCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue   []byte
		id, result string
		iid        int64
		status     int
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
	sql := "SELECT status FROM customers.customers WHERE id=$1"
	err := Pool.QueryRow(context.Background(), sql, iid).Scan(&status)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Database access error: " + err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if status != 0 {
		//другой аудитор верифицировал запись
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Level:   1,
			Message: "Данная анкета уже была верифицирована другим аудитором!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	sql = "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM " +
		"(SELECT id,lastname,firstname,middle_initial,birthday,birth_place,address,passport_data,inn,phones,emails,operatorid,screen_scans,is_moscow,auditor_comment AS comment " +
		"FROM customers.customers WHERE id=$1) t"
	err = Pool.QueryRow(context.Background(), sql, id).Scan(&result)
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

//браковка анкеты клиента
func (f *faketype) DefectCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, userid, comment string
		iid, useriid        int64
		retvalue            []byte
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
	if comment = r.FormValue("comment"); comment == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан комментарий аудитора!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//бракуем анкету
	sql := "UPDATE customers.customers SET status=2,auditorid=$1,audited=CURRENT_TIMESTAMP,auditor_comment=$2 WHERE id=$3"
	_, err = Pool.Exec(context.Background(), sql, useriid, comment, iid)
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
func (f *faketype) ConfirmCustomer(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, userid, lastname, firstname, middle_initial, birthday, birth_place, address, passport, comment string
		iid, useriid                                                                                       int64
		retvalue                                                                                           []byte
		is_moscow                                                                                          pgtype.Bool
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
	if lastname = r.FormValue("lastname"); lastname == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана фамилия клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if firstname = r.FormValue("firstname"); firstname == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано имя клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if birthday = r.FormValue("birthday"); birthday == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана дата рождения клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if birth_place = r.FormValue("birth_place"); birth_place == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано место рождения клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if passport = r.FormValue("passport"); passport == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указаны паспортные данные клиента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//проверяем корректность даты рождения
	ldate, err := time.Parse("02.01.2006", birthday)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата рождения указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	middle_initial = r.FormValue("middle_initial")
	address = r.FormValue("address")
	comment = r.FormValue("comment")
	if r.FormValue("is_moscow") == "null" {
		is_moscow.Set(nil)
	} else {
		flag, _ := strconv.ParseBool(r.FormValue("is_moscow"))
		is_moscow.Set(flag)
	}
	//сохраняем данные и верифицируем анкету
	sql := "UPDATE customers.customers SET status=1,auditorid=$1,audited=CURRENT_TIMESTAMP,lastname=$2,firstname=$3," +
		"middle_initial=$4,birthday=$5,birth_place=$6,address=$7,passport_data=$8,auditor_comment=$9, is_moscow=$10 WHERE id=$11"
	_, err = Pool.Exec(context.Background(), sql, useriid, lastname, firstname, middle_initial, ldate, birth_place, address, passport, comment, is_moscow, iid)
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

//------------------------------------------------------------------------------
//отчёт по операторам за указанный период времени
func (f *faketype) AuditorReport(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue                   []byte
		result, date_from, date_to string
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
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
	ldate_from, err := time.Parse("02.01.2006 15:04", date_from)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата начала интервала указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	ldate_from = ldate_from.Add(-3 * time.Hour)

	ldate_to, err := time.Parse("02.01.2006 15:04", date_to)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата окончания интервала указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	ldate_to = ldate_to.Add(-3 * time.Hour)
	//получаем данные отчёта

	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (" +
		"SELECT b.name, COUNT(a.importid) AS total_forms," +
		"COUNT(distinct importid) AS total_customers," +
		"COUNT(A.importid) FILTER (WHERE is_moscow = true) AS total_moscow_forms," +
		"COUNT(distinct importid) FILTER (WHERE is_moscow = true) AS total_moscow_customers," +
		"COUNT(A.importid) FILTER (WHERE is_moscow = false) AS total_not_moscow_forms," +
		"COUNT(distinct importid) FILTER (WHERE is_moscow = false) AS total_not_moscow_customers," +
		"COUNT(A.importid) FILTER (WHERE is_moscow = NULL) AS total_unknown_forms," +
		"COUNT(distinct importid) FILTER (WHERE is_moscow = NULL) AS total_unknown_customers," +

		"COUNT(A.importid) FILTER (WHERE status = 1 OR status > 2) AS verified_forms," +
		"COUNT(DISTINCT A.importid) FILTER (WHERE status = 1 OR status > 2) AS verified_customers," +
		"COUNT(A.importid) FILTER (WHERE (status = 1 OR status > 2) AND is_moscow = true)  AS verified_moscow_forms," +
		"COUNT(DISTINCT A.importid) FILTER (WHERE (status = 1 OR status > 2) AND is_moscow = true) AS verified_moscow_customers," +
		"COUNT(A.importid) FILTER (WHERE (status = 1 OR status > 2) AND is_moscow = false)  AS verified_not_moscow_forms," +
		"COUNT(DISTINCT A.importid) FILTER (WHERE (status = 1 OR status > 2) AND is_moscow = false) AS verified_not_moscow_customers," +
		"COUNT(A.importid) FILTER (WHERE (status = 1 OR status > 2) AND is_moscow = NULL)  AS verified_unknown_forms," +
		"COUNT(DISTINCT A.importid) FILTER (WHERE (status = 1 OR status > 2) AND is_moscow = false) AS verified_unknown_customers," +

		"COUNT(A.operatorid) FILTER (WHERE status = 2) AS defected_forms," +
		"COUNT(DISTINCT A.operatorid) FILTER (WHERE status = 2) AS defected_customers, " +
		"COUNT(A.operatorid) FILTER (WHERE status = 2 AND is_moscow = true) AS defected_moscow_forms," +
		"COUNT(DISTINCT A.operatorid) FILTER (WHERE status = 2 AND is_moscow = true) AS defected_moscow_customers, " +
		"COUNT(A.operatorid) FILTER (WHERE status = 2 AND is_moscow = false) AS defected_not_moscow_forms," +
		"COUNT(DISTINCT A.operatorid) FILTER (WHERE status = 2 AND is_moscow = false) AS defected_not_moscow_customers, " +
		"COUNT(A.operatorid) FILTER (WHERE status = 2 AND is_moscow = NULL) AS defected_unknown_forms," +
		"COUNT(DISTINCT A.operatorid) FILTER (WHERE status = 2 AND is_moscow = NULL) AS defected_unknown_customers " +

		"FROM customers.customers A JOIN settings.users B ON A.operatorid = B.id  " +
		"WHERE A.inserted >= $1 AND  A.inserted <= $2 GROUP BY B.name" +
		") t	"
	err = Pool.QueryRow(context.Background(), sql, ldate_from, ldate_to).Scan(&result)
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
//отчёт по паспортистам за указанный период времени
func (f *faketype) PassportReport(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		retvalue                   []byte
		result, date_from, date_to string
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
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
	ldate_from, err := time.Parse("02.01.2006 15:04", date_from)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата начала интервала указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	ldate_from = ldate_from.Add(-3 * time.Hour)

	ldate_to, err := time.Parse("02.01.2006 15:04", date_to)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Дата окончания интервала указана некорректно!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	ldate_to = ldate_to.Add(-3 * time.Hour)
	//получаем данные отчёта

	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (" +
		"SELECT b.name," +
		"COUNT(A.pasportistid) FILTER (WHERE status = 4 OR status > 5) AS verified_forms," +
		"COUNT(A.pasportistid) FILTER (WHERE status = 5) AS defected_forms " +
		"FROM customers.customers A JOIN settings.users B ON A.pasportistid = B.id  " +
		"WHERE A.pasported >= $1 AND A.pasported <= $2 GROUP BY B.name" +
		") t"
	err = Pool.QueryRow(context.Background(), sql, ldate_from, ldate_to).Scan(&result)
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
