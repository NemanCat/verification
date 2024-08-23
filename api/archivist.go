// ИС Верификатор
// методы API АРМ Архивариуса
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

//------------------------------------------------------------------------------
//функции доступа к справочнику организационно - правовых форм

// получение списка записей
func (f *faketype) GetFormsList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.forms WHERE id<>1 ORDER BY name) t"
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

// получение списка записей включая запись "нет сведений"
func (f *faketype) GetFormsFullList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.forms ORDER BY name) t"
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

// добавление/редактирование записи в справочник организационно - правовых форм
func (f *faketype) EditForm(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, name string
		iid      int64
		err      error
		retvalue []byte
		exists   bool
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id агента должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано название агента",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	if iid == 0 {
		//добавление новой записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.forms WHERE name=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Организационно-правовая форма с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, добавляем новую запись
		var new_id int64
		sql = "INSERT INTO settings.forms (name) VALUES ($1) RETURNING id"
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&new_id)

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
		io.WriteString(w, string(retvalue))
		return
	} else {
		//редактирование существующей записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.forms WHERE name=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Организационно-правовая с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, редактируем запись
		sql = "UPDATE settings.forms SET name = $1 WHERE id = $2"
		_, err = Pool.Exec(context.Background(), sql, name, iid)

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
}

// удаление записи из справочника организационно - правовых форм
func (f *faketype) DeleteForm(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		err      error
		retvalue []byte
		iid      int64
		id       string
		exists   bool
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id организационно - правовой формы!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id организационно - правовой формы должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id организационно - правовой формы должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//---------------------------------------
	//нельзя удалить стандартную запись "нет сведений"
	if iid == 1 {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Нельзя удалить запись по умолчанию из справочника!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//проверяем наличие ссылок на запись в справочнике агентов
	sql := "SELECT EXISTS(SELECT 1 FROM settings.agents WHERE formid=$1) AS \"exists\""
	err = Pool.QueryRow(context.Background(), sql, iid).Scan(&exists)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if exists {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не могу удалить запись - на неё имеется ссылка в справочнике агентов!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//проверяем наличие ссылок на запись в справочнике собственников
	sql = "SELECT EXISTS(SELECT 1 FROM settings.owners WHERE formid=$1) AS \"exists\""
	err = Pool.QueryRow(context.Background(), sql, iid).Scan(&exists)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if exists {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не могу удалить запись - на неё имеется ссылка в справочнике собственников!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//---------------------------------------
	//проверка пройдена, удаляем запись
	sql = "DELETE FROM settings.forms WHERE id=$1"
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
//функции доступа к справочнику агентов

// получение списка агентов
func (f *faketype) GetAgentsList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.agents ORDER BY name) t"
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

// добавление/редактирование записи в справочник агентов
func (f *faketype) EditAgent(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, formid, name, registration_info, address, licension_date, licension_number, status string
		iid, iformid, istatus                                                                  int64
		err                                                                                    error
		retvalue                                                                               []byte
		exists                                                                                 bool
		ldate                                                                                  time.Time
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id агента должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if formid = r.FormValue("formid"); formid == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана организационно-правовая форма агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iformid, err = strconv.ParseInt(formid, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Организационно-правовая форма агента должна быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано название агента",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	licension_number = r.FormValue("licension_number")
	registration_info = r.FormValue("registration_info")
	address = r.FormValue("address")
	licension_number = r.FormValue("licension_number")
	licension_date = r.FormValue("licension_date")
	if len(licension_date) > 0 {
		//проверяем корректность даты выдачи лицензии на туристическую деятельность
		//01.01.1901 означает что дата не указана
		ldate, err = time.Parse("02.01.2006", licension_date)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Дата выдачи лмицензии указана некорректно!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if status = r.FormValue("status"); status == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан статус агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		istatus, err = strconv.ParseInt(status, 10, 64)

		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Статус агента должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if iid == 0 {
		//добавление новой записи
		//проверяем наличие в БД агента с таким названием
		sql := "SELECT EXISTS(SELECT 1 FROM settings.agents WHERE name=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Агент с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, добавляем новую запись
		var new_id int64
		//если указана дата 01.01.1901 - сохраняем в поле даты выдачи лицензии NULL
		if licension_date == "01.01.1901" {
			sql = "INSERT INTO settings.agents (name,formid,registration_info,address,licension_number,status) " +
				"VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
			err = Pool.QueryRow(context.Background(), sql, name, iformid, registration_info, address, licension_number, istatus).Scan(&new_id)
		} else {
			sql = "INSERT INTO settings.agents (name,formid,registration_info,address,licension_number,licension_date,status) " +
				"VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
			err = Pool.QueryRow(context.Background(), sql, name, iformid, registration_info, address, licension_number, ldate, istatus).Scan(&new_id)
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
				Data:    new_id,
			})
		}
		io.WriteString(w, string(retvalue))
		return
	} else {
		//редактирование существующей записи
		//проверяем наличие в БД агента с таким названием
		sql := "SELECT EXISTS(SELECT 1 FROM settings.agents WHERE name=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Агент с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, редактируем запись
		if licension_date == "01.01.1901" {
			sql := "UPDATE settings.agents SET name = $1,formid = $2,registration_info = $3,address = $4,licension_number = $5,status = $6 WHERE id = $7"
			_, err = Pool.Exec(context.Background(), sql, name, iformid, registration_info, address, licension_number, istatus, iid)
		} else {
			sql := "UPDATE settings.agents SET name = $1,formid = $2,registration_info = $3,address = $4,licension_number = $5,licension_date = $6,status = $7 WHERE id = $8"
			_, err = Pool.Exec(context.Background(), sql, name, iformid, registration_info, address, licension_number, ldate, istatus, iid)
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
		io.WriteString(w, string(retvalue))
		return
	}
}

// удаление записи из справочника агентов
func (f *faketype) DeleteAgent(w http.ResponseWriter, r *http.Request, session *Session) {
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
			Message: "Не указан id агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id агента должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id агента должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//---------------------------------------
	//TO DO
	//проверяем наличие ссылок на запись в списке сертификатов

	//---------------------------------------

	sql := "DELETE FROM settings.agents WHERE id=$1"
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
// функции доступа к справочнику управляющих компаний

// получение списка записей
func (f *faketype) GetManagementCompaniesList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.management_companies ORDER BY name) t"
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

// добавление/редактирование записи
func (f *faketype) EditManagementCompany(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, name, address, contacts string
		iid                         int64
		err                         error
		retvalue                    []byte
		exists                      bool
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

	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано название управляющей компании!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	address = r.FormValue("address")
	contacts = r.FormValue("contacts")
	if iid == 0 {
		//добавление новой записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.management_companies WHERE name=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Управляющая компания с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, добавляем новую запись
		var new_id int64
		sql = "INSERT INTO settings.management_companies (name,address,contacts) VALUES ($1,$2,$3) RETURNING id"
		err = Pool.QueryRow(context.Background(), sql, name, address, contacts).Scan(&new_id)

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
		io.WriteString(w, string(retvalue))
		return
	} else {
		//редактирование существующей записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.management_companies WHERE name=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Управляющая компания с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, редактируем запись
		sql = "UPDATE settings.management_companies SET name = $1,address = $2,contacts = $3 WHERE id = $4"
		_, err = Pool.Exec(context.Background(), sql, name, address, contacts, iid)

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
}

// удаление записи из справочника
func (f *faketype) DeleteManagementCompany(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		err      error
		retvalue []byte
		iid      int64
		id       string
		exists   bool
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
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//---------------------------------------
	//проверяем наличие ссылок на запись в справочнике курортов
	sql := "SELECT EXISTS(SELECT 1 FROM settings.resorts WHERE management_companyid=$1) AS \"exists\""
	err = Pool.QueryRow(context.Background(), sql, iid).Scan(&exists)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if exists {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не могу удалить запись - на неё имеется ссылка в справочнике курортов!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	//---------------------------------------
	//проверка пройдена, удаляем запись
	sql = "DELETE FROM settings.management_companies WHERE id=$1"
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
// функции доступа к справочнику уполномоченных маркетинговых агентов

// получение списка записей
func (f *faketype) GetMarketingAgentsList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.marketing_agents ORDER BY name) t"
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

// добавление/редактирование записи
func (f *faketype) EditMarketingAgent(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, name, address string
		iid               int64
		err               error
		retvalue          []byte
		exists            bool
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

	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано название маркетингового агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	address = r.FormValue("address")

	if iid == 0 {
		//добавление новой записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.marketing_agents WHERE name=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Маркетинговый агент с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, добавляем новую запись
		var new_id int64
		sql = "INSERT INTO settings.marketing_agents (name,address) VALUES ($1,$2) RETURNING id"
		err = Pool.QueryRow(context.Background(), sql, name, address).Scan(&new_id)

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
		io.WriteString(w, string(retvalue))
		return
	} else {
		//редактирование существующей записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.marketing_agents WHERE name=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Маркетинговый агент с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, редактируем запись
		sql = "UPDATE settings.marketing_agents SET name = $1,address = $2 WHERE id = $3"
		_, err = Pool.Exec(context.Background(), sql, name, address, iid)

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
}

// удаление записи из справочника
func (f *faketype) DeleteMarketingAgent(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		err      error
		retvalue []byte
		iid      int64
		id       string
		//		exists   bool
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
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//---------------------------------------
	//TO DO
	//проверяем наличие ссылок на запись в списке сертификатов
	//---------------------------------------
	//проверка пройдена, удаляем запись
	sql := "DELETE FROM settings.marketing_agents WHERE id=$1"
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
//функции доступа к справочнику собственников курортов

// получение списка записей
func (f *faketype) GetOwnersList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.owners ORDER BY name) t"
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

// добавление/редактирование записи
func (f *faketype) EditOwner(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, formid, name, address, registration_info, requisites string
		iid, iformid                                             int64
		err                                                      error
		retvalue                                                 []byte
		exists                                                   bool
	)
	//парсим данные запроса с проверкой их значений
	r.ParseForm()
	if id = r.FormValue("id"); id == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iid, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id агента должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if formid = r.FormValue("formid"); formid == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указана организационно-правовая форма агента!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iformid, err = strconv.ParseInt(formid, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Организационно-правовая форма агента должна быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано название собственника курорта!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	address = r.FormValue("address")
	registration_info = r.FormValue("registration_info")
	requisites = r.FormValue("requisites")

	if iid == 0 {
		//добавление новой записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.owners WHERE name=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Собственник курорта с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, добавляем новую запись
		var new_id int64

		sql = "INSERT INTO settings.owners (name,formid,address,registration_info,requisites) VALUES ($1,$2,$3,$4,$5) RETURNING id"
		err = Pool.QueryRow(context.Background(), sql, name, iformid, address, registration_info, requisites).Scan(&new_id)

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
		io.WriteString(w, string(retvalue))
		return
	} else {
		//редактирование существующей записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.owners WHERE name=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Собственник курорта с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, редактируем запись
		sql = "UPDATE settings.owners SET name = $1,formid = $2,address = $3,registration_info = $4,requisites = $5 WHERE id = $6"
		_, err = Pool.Exec(context.Background(), sql, name, iformid, address, registration_info, requisites, iid)

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
}

// удаление записи
func (f *faketype) DeleteOwner(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		err      error
		retvalue []byte
		iid      int64
		id       string
		exists   bool
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
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//---------------------------------------
	//проверяем наличие ссылок на запись в справочнике курортов
	sql := "SELECT EXISTS(SELECT 1 FROM settings.resorts WHERE management_companyid=$1) AS \"exists\""
	err = Pool.QueryRow(context.Background(), sql, iid).Scan(&exists)
	if err != nil {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	if exists {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не могу удалить запись - на неё имеется ссылка в справочнике курортов!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}

	//---------------------------------------
	sql = "DELETE FROM settings.owners WHERE id=$1"
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
// функции доступа к справочнику курортов

// получение списка записей
func (f *faketype) GetResortsList(w http.ResponseWriter, r *http.Request, session *Session) {
	var retvalue []byte
	var result string
	sql := "SELECT array_to_json(coalesce(array_agg (t),'{}')) FROM (SELECT * FROM settings.resorts ORDER BY name) t"
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

// добавление/редактирование записи
func (f *faketype) EditResort(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		id, name, address, ownerid, management_companyid, status string
		iid, iownerid, imanagement_companyid, istatus            int64
		err                                                      error
		retvalue                                                 []byte
		exists                                                   bool
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

	if name = r.FormValue("name"); name == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указано название курорта!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	}
	address = r.FormValue("address")
	if ownerid = r.FormValue("ownerid"); ownerid == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id собственника!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		iownerid, err = strconv.ParseInt(ownerid, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id собственника должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}

	if management_companyid = r.FormValue("management_companyid"); management_companyid == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан id управляющей компании!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		imanagement_companyid, err = strconv.ParseInt(management_companyid, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id управляющей компании должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}

	if status = r.FormValue("status"); status == "" {
		retvalue, _ = json.Marshal(Retvalue{
			Success: false,
			Message: "Не указан статус курорта!",
			Data:    nil,
		})
		io.WriteString(w, string(retvalue))
		return
	} else {
		istatus, err = strconv.ParseInt(status, 10, 64)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Статус курорта должен быть целым числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}

	if iid == 0 {
		//добавление новой записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.resorts WHERE name=$1) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Курорт с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, добавляем новую запись
		var new_id int64

		sql = "INSERT INTO settings.resorts (name,address,ownerid,management_companyid,status) VALUES ($1,$2,$3,$4,$5) RETURNING id"
		err = Pool.QueryRow(context.Background(), sql, name, address, iownerid, imanagement_companyid, istatus).Scan(&new_id)

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
		io.WriteString(w, string(retvalue))
		return
	} else {
		//редактирование существующей записи
		//проверяем уникальность записи
		sql := "SELECT EXISTS(SELECT 1 FROM settings.resorts WHERE name=$1 AND id <> $2) AS \"exists\""
		err = Pool.QueryRow(context.Background(), sql, name, iid).Scan(&exists)
		if err != nil {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		if exists {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "Курорт с таким названием уже есть в базе данных!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
		//проверка пройдена, редактируем запись
		sql = "UPDATE settings.resorts SET name=$1,address=$2,ownerid=$3,management_companyid=$4,status=$5 WHERE id = $6"
		_, err = Pool.Exec(context.Background(), sql, name, address, iownerid, imanagement_companyid, istatus, iid)

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
}

// удаление записи из справочника
func (f *faketype) DeleteResort(w http.ResponseWriter, r *http.Request, session *Session) {
	var (
		err      error
		retvalue []byte
		iid      int64
		id       string
		//		exists   bool
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
		if iid < 1 {
			retvalue, _ = json.Marshal(Retvalue{
				Success: false,
				Message: "id записи должен быть целым положительным числом!",
				Data:    nil,
			})
			io.WriteString(w, string(retvalue))
			return
		}
	}
	//---------------------------------------
	//TO DO
	//проверяем наличие ссылок на запись в списке сертификатов
	//---------------------------------------
	//проверка пройдена, удаляем запись
	sql := "DELETE FROM settings.resorts WHERE id=$1"
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
