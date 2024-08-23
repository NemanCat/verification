// сервис подсистемы безопасности
// список информационных объектов в памяти
package main

import (
	"sync"
)

type MemoryMap struct {
	//мьютекс для многопоточного использования
	mutex *sync.RWMutex
	//объекты
	list map[string]BaseObjectInterface
	//сохранять ли объекты в хранилище
	persistent bool
	//последовательность целочисленных id объектов
	sequence uint64
	//функция загрузки списка объектов из БД
	//nil для списков, не сохраняемых в хранилище
	loadFromPersistence func() (map[string]BaseObjectInterface, error)
	//функция сохранения объекта в БД
	//nil для списков, не сохраняемых в хранилище
	persist func(key string, value BaseObjectInterface) error
	//функция удаления объекта из БД
	//nil для списков, не сохраняемых в хранилище
	remove func(key string) error
}

//-----------------------------------------------
//методы работы со списком
//создание списка
func NewMemoryMap(persistent bool, loadFromPersistence func() (map[string]BaseObjectInterface, error), persist func(key string, value BaseObjectInterface) error, remove func(key string) error) (*MemoryMap, error) {
	mm := new(MemoryMap)
	mm.mutex = new(sync.RWMutex)
	mm.list = make(map[string]BaseObjectInterface)
	mm.persistent = persistent
	mm.loadFromPersistence = loadFromPersistence
	mm.persist = persist
	mm.remove = remove
	if persistent == true {
		pl, err := mm.loadFromPersistence()
		if err != nil {
			return nil, err
		}
		var maxid uint64 = 0
		for key, value := range pl {
			mm.AddObject(key, value)
			if value.GetId() > maxid {
				maxid = value.GetId()
			}
		}
		mm.sequence = maxid
		return mm, nil
	} else {
		mm.sequence = 0
		return mm, nil
	}
}

//получение списка всех элементов
func (memory_map *MemoryMap) GetList() map[string]BaseObjectInterface {
	return memory_map.list
}

//получение элемента по указанному ключу
//возвращает ссылку на элемент
func (memory_map *MemoryMap) FindObject(key string) BaseObjectInterface {
	memory_map.mutex.Lock()
	defer memory_map.mutex.Unlock()
	obj, _ := memory_map.list[key]
	return obj
}

//количество объектов в списке
func (memory_map *MemoryMap) Count() int {
	memory_map.mutex.Lock()
	defer memory_map.mutex.Unlock()
	return len(memory_map.list)
}

//генерация id нового объекта
func (memory_map *MemoryMap) CreateId() uint64 {
	memory_map.mutex.Lock()
	defer memory_map.mutex.Unlock()
	memory_map.sequence++
	return memory_map.sequence
}

//добавление в список элемента с указанным ключом
//для существующего ключа обновляет запись
func (memory_map *MemoryMap) AddObject(key string, value BaseObjectInterface) error {
	memory_map.mutex.Lock()
	defer memory_map.mutex.Unlock()
	memory_map.list[key] = value
	if memory_map.persistent == true {
		err := memory_map.persist(key, value)
		return err
	} else {
		return nil
	}
}

//удаление из списка элемента с указанным ключом
func (memory_map *MemoryMap) DeleteObject(key string) error {
	memory_map.mutex.Lock()
	defer memory_map.mutex.Unlock()
	delete(memory_map.list, key)
	if memory_map.persistent == true {
		err := memory_map.remove(key)
		return err
	} else {
		return nil
	}
}
