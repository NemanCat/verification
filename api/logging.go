// ИС Верификатор
// сервис логирования
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
)

//------------------------------------------------------------------
type LoggingService struct {
	base_logging_dir string
	max_size         int
	backups          int
	max_age          int
}

//функция записи сообщения в лог
func (service *LoggingService) WriteLogMessage(logname string, message string) error {
	//директория вывода лога
	destination_dir := filepath.Join(service.base_logging_dir, logname)
	//проверяем существует ли директория вывода лога
	//если не существует - создаём её
	if _, err := os.Stat(destination_dir); os.IsNotExist(err) {
		err = os.Mkdir(destination_dir, 0755)
		if err != nil {
			fmt.Println("Could not create logging dir with error message " + err.Error())
			return err
		}
	}
	//имя файла лога
	logfile_name := filepath.Join(destination_dir, logname+".log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   logfile_name,
		MaxSize:    service.max_size,
		MaxBackups: service.backups,
		MaxAge:     service.max_age,
		Compress:   true,
		LocalTime:  true,
	})
	log.Println(message)
	return nil
}
