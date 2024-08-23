// КОРПРЕСУРС Система онлайн-тестирования
// сервис единой точки доступа к API
// чистильщик просроченных сессий пользователей
package main

import (
	"context"
	"sync"
	"time"
)

// очистка списка сессий от устаревших записей
func clearSessions(session_lifetime int) {
	for key, value := range SessionsList.GetList() {
		//проверяем дату и время последней активности пользователя
		diff := time.Now().Sub(value.(*Session).Last_activity)
		if diff.Minutes() > float64(session_lifetime) {
			if log_level == "1" {
				message := "Пользователь " + value.(*Session).Name +
					" (" + value.(*Session).Login + "), id сессии " + value.(*Session).Sessionid + ", ip-адрес " + value.(*Session).Ip +
					" сессия удалена чистильщиком сессий, время последней активности " +
					value.(*Session).Last_activity.String()
				go logging_service.WriteLogMessage("authorization", message)
			}
			//превышен максимальный интервал ожидания
			//удаляем сессию
			SessionsList.DeleteObject(key)
			return
		}
	}
}

// возвращение в обработку подвисших записей в таблице импорта
func restoreImportRecords() {
	sql := "UPDATE customers.import SET status = 0, operatorid = NULL, operator_block_datetime = NULL " +
		"WHERE status = 1 AND (DATE_PART('day', CURRENT_TIMESTAMP::timestamp - operator_block_datetime::timestamp) * 24 + " +
		"DATE_PART('hour', CURRENT_TIMESTAMP::timestamp - operator_block_datetime::timestamp)) * 60 + " +
		"DATE_PART('minute', CURRENT_TIMESTAMP::timestamp - operator_block_datetime::timestamp) > 120"
	_, err := Pool.Exec(context.Background(), sql)
	if err != nil {
		logging_service.WriteLogMessage("cleaner", "Could not restore import table records with error message: "+err.Error())
	}
}

//interval - интервал запуска сервиса в минутах
//session_lifetime - продолжительность неактивности сессии пользователя, минут
//shutdownChannel - канал для сигнала остановки сервиса
//waitGroup - синхронизированная группа дочерних серверов
func Cleaner(interval int, session_lifetime int, shutdownChannel chan bool, waitGroup *sync.WaitGroup) {
	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	defer ticker.Stop()
	defer waitGroup.Done()
	for {
		select {
		case _ = <-shutdownChannel:
			//прекращение работы
			return
		case <-ticker.C:
			go clearSessions(session_lifetime)
			go restoreImportRecords()
		default:

		}
	}
}
