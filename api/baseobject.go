// КОРПРЕСУРС Система онлайн-тестирования
// сервис подсистемы безопасности
// базовый информационный объект - включается во все информационные объекты
package main

//разделитель полей при кодировании объектов
var Delimiter []byte = []byte("{}")

type BaseObject struct {
	//целочисленный идентификатор объекта
	Id uint64
}

func (obj BaseObject) GetId() uint64 {
	return obj.Id
}

func (obj BaseObject) SetId(id uint64) {
	obj.Id = id
}

//методы, которые должны быть реализованы во всех объектах
type BaseObjectInterface interface {
	//получение целочисленного id объекта
	GetId() uint64
	//декодировка объекта из JSON
	UnmarshalJSON([]byte) string
}
