// КОРПРЕСУРС Система онлайн-тестирования
// сервис подсистемы безопасности
// провайдер key-value БД для хранения оперативной информации

package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/boltdb/bolt"
)

type BoltDbProvider struct {
	//мьютекс для многопоточного использования
	mutex *sync.RWMutex
	//открытые базы данных
	databases map[string]*bolt.DB
	//директория файлов БД
	dbfolder string
	//имя файла БД настроек системы
	//kernel_name string
}

//--------------------------------------------------
//не экспортируемые функции
//получение ссылки на БД
//возвращает ссылку на базу или nil если база не открывалась
func (provider *BoltDbProvider) getDB(name string) *bolt.DB {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()
	return provider.databases[name]
}

//открытие указанной БД и добавление её в список открытых БД
func (provider *BoltDbProvider) openDB(name string, filename string) (*bolt.DB, error) {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}
	provider.databases[name] = db
	return db, nil
}

//поиск по ключу
func (provider *BoltDbProvider) findByKey(dbname string, filename string, bucketname string, key string, objtype reflect.Type) (BaseObjectInterface, error) {
	var db *bolt.DB
	var err error
	//подключаемся к БД
	db = provider.getDB(dbname)
	if db == nil {
		db, err = provider.openDB(dbname, provider.dbfolder+"/"+filename)
		if err != nil {
			return nil, err
		}
	}
	rec := reflect.New(objtype)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketname))
		if bucket != nil {
			found := bucket.Get([]byte(key))
			if found == nil {
				//запись не найдена
				rec = reflect.Zero(objtype)
				return nil
			} else {
				//запись найдена
				method := rec.MethodByName("UnmarshalJSON")
				inputs := make([]reflect.Value, 1)
				inputs[0] = reflect.ValueOf(found)
				res := method.Call(inputs)
				message := res[0].String()
				if len(message) > 0 {
					rec = reflect.Zero(objtype)
					return errors.New(message)
				} else {
					return nil
				}
			}
		} else {
			rec = reflect.Zero(objtype)
			return nil
		}
	})

	if rec == reflect.Zero(objtype) {
		return nil, nil
	} else {
		return rec.Interface().(BaseObjectInterface), err
	}
}

//количество записей
func (provider *BoltDbProvider) count(dbname string, filename string, bucketname string) (uint32, error) {
	var db *bolt.DB
	var err error
	//подключаемся к БД
	db = provider.getDB(dbname)
	if db == nil {
		db, err = provider.openDB(dbname, provider.dbfolder+"/"+filename)
		if err != nil {
			return 0, err
		}
	}
	var count uint32
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketname))
		if bucket == nil {
			//если букет не существует, создаём его
			bucket, err = tx.CreateBucketIfNotExists([]byte(bucketname))
			if err != nil {
				return err
			}
		}
		count = uint32(bucket.Stats().KeyN)
		return nil
	})
	return count, err
}

//генерация целочисленного id
func (provider *BoltDbProvider) createid(dbname string, filename string, bucketname string) (uint64, error) {
	var db *bolt.DB
	var err error
	//подключаемся к БД
	db = provider.getDB(dbname)
	if db == nil {
		db, err = provider.openDB(dbname, provider.dbfolder+"/"+filename)
		if err != nil {
			return 0, err
		}
	}
	var id uint64
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketname))
		if bucket == nil {
			//если букет не существует, создаём его
			bucket, err = tx.CreateBucketIfNotExists([]byte(bucketname))
			if err != nil {
				return err
			}
		}
		id, err = bucket.NextSequence()
		return err
	})
	return id, err
}

//--------------------------------------------------
//базовые функции провайдера
//конструктор
func NewBoltDbProvider(dbfolder string /*, kernel_name string*/) *BoltDbProvider {
	nbdp := new(BoltDbProvider)
	nbdp.dbfolder = dbfolder
	//	nbdp.kernel_name = kernel_name
	nbdp.databases = make(map[string]*bolt.DB)
	nbdp.mutex = new(sync.RWMutex)
	return nbdp
}

//закрытие всех открытых БД
func (provider *BoltDbProvider) Close() {
	for _, value := range provider.databases {
		value.Close()
	}
}

//получение списка записей
func (provider *BoltDbProvider) List(dbname string, filename string, bucketname string, objtype reflect.Type) (map[string]BaseObjectInterface, error) {
	var db *bolt.DB
	var err error
	//подключаемся к БД
	db = provider.getDB(dbname)
	if db == nil {
		db, err = provider.openDB(dbname, provider.dbfolder+"/"+filename)
		if err != nil {
			return nil, err
		}
	}

	list := make(map[string]BaseObjectInterface)
	err = db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(bucketname))
		if bucket != nil {
			bucket.ForEach(func(k, v []byte) error {
				var tmp = reflect.New(objtype)
				method := tmp.MethodByName("UnmarshalJSON")
				inputs := make([]reflect.Value, 1)
				inputs[0] = reflect.ValueOf(v)
				res := method.Call(inputs)
				message := res[0].String()
				if len(message) != 0 {
					return errors.New(message)
				}
				list[string(k)] = tmp.Interface().(BaseObjectInterface)
				return nil
			})
			return nil
		}
		return nil
	})
	return list, err
}

//удаление записи
func (provider *BoltDbProvider) Delete(dbname string, filename string, bucketname string, key string) error {
	var db *bolt.DB
	var err error
	//подключаемся к БД
	db = provider.getDB(dbname)
	if db == nil {
		db, err = provider.openDB(dbname, provider.dbfolder+"/"+filename)
		if err != nil {
			return err
		}
	}
	err = db.Update(
		func(tx *bolt.Tx) error {
			var bucket *bolt.Bucket
			var err error
			bucket = tx.Bucket([]byte(bucketname))
			if bucket != nil {
				//удаляем
				err = bucket.Delete([]byte(key))
				if err != nil {
					return err
				} else {
					return nil
				}
			} else {
				return nil
			}
		})
	return err
}

//сохранение записи
func (provider *BoltDbProvider) Persist(dbname string, filename string, bucketname string, key string, value BaseObjectInterface) error {
	var db *bolt.DB
	var err error
	//подключаемся к БД
	db = provider.getDB(dbname)
	if db == nil {
		db, err = provider.openDB(dbname, provider.dbfolder+"/"+filename)
		if err != nil {
			return err
		}
	}
	//encoded, err := value.MarshalJSON()
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	//fmt.Println(encoded)
	err = db.Update(
		func(tx *bolt.Tx) error {
			var bucket *bolt.Bucket
			var err error
			bucket = tx.Bucket([]byte(bucketname))
			if bucket == nil {
				//если букет не существует, создаём его
				bucket, err = tx.CreateBucketIfNotExists([]byte(bucketname))
				if err != nil {
					return err
				}
			}
			//сохраняем
			err = bucket.Put([]byte(key), []byte(encoded))
			if err != nil {
				return err
			}
			return nil
		})
	return err
}
