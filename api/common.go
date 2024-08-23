// ИС Верификатор
// глобальные определения системы
package main

import (
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

//-------------------------------------------------------
//формальный тип для использования рефлексии
type faketype int

//--------------------------------------------------------
//структура стандартного ответа API
type Retvalue struct {
	//флаг успешности операции
	Success bool `json:"Success"`
	//сообщение об ошибке в случае наличия ошибок
	Message string `json:"Message"`
	//возвращаемые данные
	Data interface{} `json:"Data"`
	//уровень серьёзности ошибки
	//0 - ошибка
	//1 - информация
	Level int `default:"0"`
}

//информационный объект Маршрут
type Route struct {
	//требуется ли авторизация
	Is_protected bool
	//функция - обработчик
	Handler string
	//категория пользователя
	Category int32
}

//структура контактный телефон
type Phone struct {
	//номер телефона
	Phone string `json:"phone"`
	//дата и время прозвона
	Call_datetime string `json:"call_datetime"`
	//результат прозвона
	//0-не проверялся
	//1-дозвонились успешно
	//2-номер продиктован клиентом
	//3-не дозвонились
	//4-номер не существует или не обслуживается
	//5- номер принадлежит другому человеку
	Result int16 `json:"result"`
}

//структура адрес электронной почты
type Email struct {
	Email string `json:"email"`
}

//структура скан экрана
type Scan struct {
	//имя файла скана
	Filename string `json:"filename"`
	//графическое изображение закодированное в Base64
	Img string `json:"img"`
}

//структура данных клиента при сохранении анкеты клиента
type AssemblyCustomerData struct {
	//id записи
	Id string `json:"id"`
	//фамилия
	Lastname string `json:"lastname"`
	//имя
	Firstname string `json:"firstname"`
	//отчество
	Middle_initial string `json:"middle_initial"`
	//дата рождения
	Birthday string `json:"birthday"`
	//место рождения
	Birth_place string `json:"birth_place"`
	//адрес регистрации
	Address string `json:"address"`
	//паспортные данные
	Passport string `json:"passport"`
	//ИНН
	Inn string `json:"inn"`
	//комментарий
	Comment string `json:"comment"`
	//адреса электронной почты
	Emails []Email `json:"emails"`
	//список номеров телефонов
	Phones []Phone `json:"phones"`
	//Id телефониста
	Userid string `json:"userid"`
	//статус клиента
	//4-все телефоны бракованные
	//5-отказался от сотрудничества
	//6-согласен на сотрудничество
	//Status int16 `json:"status"`
	//список сканов экрана
	Scans []Scan `json:"scans"`
}

//---------------------------------------------------------
const (
	//секретное слово для подписи cookie
	cookie_secret = "gtrVajNmJLZ4ZbNkFZnRjcFT9SpGHSQB"
	//секретное слово для формирования CSFR security token
	csfr_secret = "HkVxPQyBYZhT92UaaQvbrwQdegUHJUxJ"
	//имя сессионного cookie
	session_cookie_name = "verification-session"
	//имя cookie с подписью сессии
	session_signature_cookie_name = "verification-session-signature"
)

//глобальные переменные
var (
	//пул подключений к БД Postgres
	Pool *pgxpool.Pool
	//список сессий в памяти
	SessionsList *MemoryMap
	//провайдер доступа к БД приложения
	BoltProvider    *BoltDbProvider
	logging_service LoggingService
	//уровень логгирования
	log_level string
)

//глобальные функции
// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

//------------------------------------------------------------------
//список маршрутов и обрабатывающих функций
var ApiRoutes = map[string]Route{
	//авторизация пользователя в системе
	"/authorization": {Handler: "Authorization", Is_protected: false, Category: -1},
	//штатный выход пользователя из системы
	"/logout": {Handler: "Logout", Is_protected: true, Category: -1},
	//выход пользователя из системы по тайм-ауту на клиенте
	"/system-logout": {Handler: "SystemLogout", Is_protected: true, Category: -1},
	//--------------------------------------------------------------------------
	//API АРМ Администратора
	//получение списка пользователей
	"/admin/users": {Handler: "GetUsersList", Is_protected: true, Category: 0},
	//добавление/редактирование аккаунта пользователя
	"/admin/users/edit": {Handler: "EditUser", Is_protected: true, Category: 0},
	//удаление аккаунта пользователя
	"/admin/users/remove": {Handler: "DeleteUser", Is_protected: true, Category: 0},
	//получение списка пользовательских сессий
	"/admin/sessions": {Handler: "GetSessionsList", Is_protected: false, Category: 0},
	//--------------------------------------------------------------------------
	//API АРМ Архивариуса
	//получение списка организационно - правовых форм
	"/archivist/forms": {Handler: "GetFormsList", Is_protected: true, Category: 1},
	//получение списка организационно - правовых форм включая запись "по умолчанию"
	"/archivist/forms/full-list": {Handler: "GetFormsFullList", Is_protected: true, Category: 1},
	//добавление/редактирование записи в справочнике организационно - правовых форм
	"/archivist/forms/edit": {Handler: "EditForm", Is_protected: true, Category: 1},
	//удаление записи из справочника организационно - правовых форм
	"/archivist/forms/remove": {Handler: "DeleteForm", Is_protected: true, Category: 1},
	//------------------------
	//получение списка агентов (продавцов)
	"/archivist/agents": {Handler: "GetAgentsList", Is_protected: true, Category: 1},
	//добавление/редактирование записи в справочнике агентов
	"/archivist/agents/edit": {Handler: "EditAgent", Is_protected: true, Category: 1},
	//удаление записи из справочника агентов
	"/archivist/agents/remove": {Handler: "DeleteAgent", Is_protected: true, Category: 1},
	//------------------------
	//получение списка управляющих компаний
	"/archivist/management-companies": {Handler: "GetManagementCompaniesList", Is_protected: true, Category: 1},
	//добавление/редактирование записи в справочнике управляющих компаний
	"/archivist/management-companies/edit": {Handler: "EditManagementCompany", Is_protected: true, Category: 1},
	//удаление записи из справочника управляющих компаний
	"/archivist/management-companies/remove": {Handler: "DeleteManagementCompany", Is_protected: true, Category: 1},
	//------------------------
	//получение списка маркетинговых агентов
	"/archivist/marketing-agents": {Handler: "GetMarketingAgentsList", Is_protected: true, Category: 1},
	//добавление/редактирование записи в справочнике маркетинговых агентов
	"/archivist/marketing-agents/edit": {Handler: "EditMarketingAgent", Is_protected: true, Category: 1},
	//удаление записи из справочника маркетинговых агентов
	"/archivist/marketing-agents/remove": {Handler: "DeleteMarketingAgent", Is_protected: true, Category: 1},
	//------------------------
	//получение списка собственников курортов
	"/archivist/owners": {Handler: "GetOwnersList", Is_protected: true, Category: 1},
	//добавление/редактирование записи в справочнике собственников курортов
	"/archivist/owners/edit": {Handler: "EditOwner", Is_protected: true, Category: 1},
	//удаление записи из справочника маркетинговых собственников курортов
	"/archivist/owners/remove": {Handler: "DeleteOwner", Is_protected: true, Category: 1},
	//------------------------
	//получение списка курортов
	"/archivist/resorts": {Handler: "GetResortsList", Is_protected: true, Category: 1},
	//добавление/редактирование записи в справочнике курортов
	"/archivist/resorts/edit": {Handler: "EditResort", Is_protected: true, Category: 1},
	//удаление записи из справочника курортов
	"/archivist/resorts/remove": {Handler: "DeleteResort", Is_protected: true, Category: 1},
	//--------------------------------------------------------------------------
	//API АРМ Оператора БД
	//получение импортированной анкеты на обработку
	"/operator/import/get-customer": {Handler: "GetImportedCustomer", Is_protected: true, Category: 2},
	//подтверждение обработки анкеты оператором
	"/operator/import/confirm-customer": {Handler: "ConfirmImportedCustomer", Is_protected: true, Category: 2},
	//браковка анкеты в таблице импорта
	"/operator/import/defect-customer": {Handler: "DefectImportedCustomer", Is_protected: true, Category: 2},
	//возврат анкеты от оператора
	"/operator/import/return-customer": {Handler: "ReturnImportedCustomer", Is_protected: true, Category: 2},
	//сохранение новой анкеты клиента
	"/operator/save-customer": {Handler: "SaveCustomer", Is_protected: true, Category: 2},
	//отчёт по анкетам за указанный период
	"/operator/report": {Handler: "OperatorReport", Is_protected: true, Category: 2},
	//--------------------------------------------------------------------------
	//API АРМ Аудитора
	//получение списка неподтверждённых анкет клиентов
	"/auditor/customers": {Handler: "GetNewCustomersList", Is_protected: true, Category: 3},
	//получение анкеты клиента для просмотра/редактирования
	"/auditor/get-customer": {Handler: "GetCustomer", Is_protected: true, Category: 3},
	//браковка анкеты клиента аудитором
	"/auditor/defect-customer": {Handler: "DefectCustomer", Is_protected: true, Category: 3},
	//одобрение анкеты клиента аудитором
	"/auditor/confirm-customer": {Handler: "ConfirmCustomer", Is_protected: true, Category: 3},
	//отчёт по операторам за указанный период
	"/auditor/report": {Handler: "AuditorReport", Is_protected: true, Category: 3},
	//отчёт по паспортистам за указанный период
	"/auditor/passport-report": {Handler: "PassportReport", Is_protected: true, Category: 3},
	//--------------------------------------------------------------------------
	//API АРМ Паспортиста
	//получение списка анкет, находящихся на обработке у данного паспортиста
	"/passport/get-customers-list": {Handler: "GetPasportistCustomersList", Is_protected: false, Category: 5},
	//получение новой верифицированной анкеты клиента на обработку
	"/passport/get-customer": {Handler: "GetVeirifedCustomer", Is_protected: false, Category: 5},
	//получение данных анкеты клиента для просмотра и редактирования
	"/passport/get-customer-data": {Handler: "GetPassportCustomerData", Is_protected: false, Category: 5},
	//возврат анкеты в общую базу
	"/passport/return-customer": {Handler: "ReturnPassportCustomer", Is_protected: false, Category: 5},
	//сохранение данных анкеты клиента
	"/passport/save-customer": {Handler: "SavePassportCustomer", Is_protected: true, Category: 5},
	//браковка анкеты клиента паспортистом
	"/passport/defect-customer": {Handler: "DefectPassportCustomer", Is_protected: true, Category: 5},
	//одобрение анкеты клиента паспортистом
	"/passport/confirm-customer": {Handler: "ConfirmPassportCustomer", Is_protected: true, Category: 5},
	//--------------------------------------------------------------------------
	//API АРМ Телефониста
	//получение списка анкет, находящихся на обработке у данного телефониста
	"/assembly/get-customers-list": {Handler: "GetAssemblyCustomersList", Is_protected: false, Category: 4},
	//получение данных анкеты клиента для просмотра и редактирования
	"/assembly/get-customer-data": {Handler: "GetAssemblyCustomerData", Is_protected: false, Category: 4},
	//получение анкеты на прозвон
	"/assembly/get-customer": {Handler: "GetCustomerForCalling", Is_protected: false, Category: 4},
	//браковка анкеты клиента телефонистом
	"/assembly/defect-customer": {Handler: "DefectAssemblyCustomer", Is_protected: true, Category: 4},
	//клиен готов к сотрудничеству
	"/assembly/confirm-customer": {Handler: "ConfirmAssemblyCustomer", Is_protected: true, Category: 4},
	//клиен не готов к сотрудничеству
	"/assembly/reject-customer": {Handler: "RejectAssemblyCustomer", Is_protected: true, Category: 4},
	//сохранение данных анкеты клиента
	"/assembly/save-customer": {Handler: "SaveCalledCustomer", Is_protected: true, Category: 4},
	//--------------------------------------------------------------------------
	//API АРМ Суперпользователя
	//получение списка анкет клиентов
	"/superuser/get-customers-list": {Handler: "GetCustomersList", Is_protected: false, Category: 6},
	//получение списка забракованных анкет из таблиц импорта
	"/superuser/get-defected-customers-list": {Handler: "GetDefectedCustomersList", Is_protected: false, Category: 6},
	//получение анкеты клиента для просмотра/редактирования
	"/superuser/get-customer": {Handler: "GetSuperuserCustomer", Is_protected: true, Category: 6},
	//получение забракованной анкеты клиента для просмотра/редактирования
	"/superuser/get-defected-customer": {Handler: "GetSuperuserDefectedCustomer", Is_protected: true, Category: 6},
	//возвращение забракованной анкеты в обработку
	"/superuser/return-defected-customer": {Handler: "ReturnDefectedCustomer", Is_protected: true, Category: 6},
	//сохранение комментария
	"/superuser/save-comment": {Handler: "SaveComment", Is_protected: true, Category: 6},
}
