// ИС Верификатор
// главный модуль
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

//------------------------------------------------------------------
var (
	//максимальный период бездействия пользователя до окончания сессии
	max_session_timeout string
	//частота запуска чистильщика просроченных сессий
	sessions_cleaner_frequency string
	//ощибка
	err error
	//функция преобразования строки к целому числу
	StrToInt = func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	}
)

//количество дочерних сервисов
const services_count int = 1

//------------------------------------------------------------------

func main() {
	//получаем настройки приложения
	godotenv.Load()
	//порт приложения
	port, exists := os.LookupEnv("port")
	if !exists {
		port = "8080"
	}
	//хост сервера БД
	pg_host, exists := os.LookupEnv("pg_host")
	if !exists {
		pg_host = "localhost"
	}
	//порт сервера БД
	pg_port, exists := os.LookupEnv("pg_port")
	if !exists {
		pg_port = "5432"
	}
	//имя пользователя БД
	pg_user, exists := os.LookupEnv("pg_user")
	if !exists {
		pg_user = "postgres"
	}
	//пароль пользователя БД
	pg_password, exists := os.LookupEnv("pg_password")
	if !exists {
		pg_password = "masterkey"
	}
	//название БД
	pg_database, exists := os.LookupEnv("pg_database")
	if !exists {
		pg_database = "verification"
	}
	//директория баз данных Bolt
	bolt_dir, exists := os.LookupEnv("bolt_dir")
	if !exists {
		bolt_dir, _ = os.Getwd()
	}
	//максимальное время бездействия пользователя до отключения сессии, минут
	max_session_timeout, exists = os.LookupEnv("max_session_timeout")
	if !exists {
		max_session_timeout = "30"
	}
	//частота запуска чистильщика просроченных сессий, минут
	sessions_cleaner_frequency, exists = os.LookupEnv("sessions_cleaner_frequency")
	if !exists {
		sessions_cleaner_frequency = "60"
	}
	//директория для записи логов
	logging_dir, exists := os.LookupEnv("logging_dir")
	if !exists {
		path, _ := os.Getwd()
		logging_dir = filepath.Join(path, "logs")
	}
	//проверяем существует ли указанная директория
	//если не существует - создаём её
	if _, err := os.Stat(logging_dir); os.IsNotExist(err) {
		err = os.Mkdir(logging_dir, 0755)
		if err != nil {
			fmt.Println("Could not create logging dir with error message " + err.Error())
			return
		}
	}
	//максимальный размер файла лога, Мб
	max_size, exists := os.LookupEnv("max_size")
	if !exists {
		max_size = "2"
	}
	//количество бэкапов файла лога
	backups, exists := os.LookupEnv("backups")
	if !exists {
		backups = "10"
	}
	//максимальный срок хранения бэкапов, дней
	max_age, exists := os.LookupEnv("max_age")
	if !exists {
		max_age = "30"
	}
	//уровень логгирования
	//0-не вести никаких логов
	//1-логгировать события авторизации и нарушения безопасности
	log_level, exists = os.LookupEnv("log_level")
	if !exists {
		log_level = "0"
	}
	//-------------------------------------------
	//создаём сервис записи сообщений в лог
	logging_service = LoggingService{
		base_logging_dir: logging_dir,
		max_size:         StrToInt(max_size),
		backups:          StrToInt(backups),
		max_age:          StrToInt(max_age),
	}
	//-------------------------------------------------------
	//пытаемся создать подключение к БД
	connectionUri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		pg_user,
		pg_password,
		pg_host,
		pg_port,
		pg_database)
	Pool, err = pgxpool.Connect(context.Background(), connectionUri)
	if err != nil {
		fmt.Println("Could not connect to database with error message: " + err.Error())
		logging_service.WriteLogMessage("main", "Could not connect to database with error message: "+err.Error())
		logging_service.WriteLogMessage("main", "Verification API service stopped")
		return
	}
	//всегда закрываем подключение при остановке приложения
	defer Pool.Close()
	//--------------------------------------------------------
	//проверяем существует ли директория баз данных Bolt
	//если не существует - создаём её
	if _, err := os.Stat(bolt_dir); os.IsNotExist(err) {
		err = os.Mkdir(bolt_dir, 0755)
		if err != nil {
			logging_service.WriteLogMessage("main", "Could not create directory for Bolt databases with error message: "+err.Error())
			logging_service.WriteLogMessage("main", "Verification API service stopped")
			return
		}
	}
	BoltProvider = NewBoltDbProvider(bolt_dir)
	defer BoltProvider.Close()
	//создаём в памяти список сессий и загружаем из БД сохранённые сессии
	SessionsList, err = NewMemoryMap(true, BoltProvider.GetSessionsList, BoltProvider.PersistSession, BoltProvider.DeleteSession)
	if err != nil {
		logging_service.WriteLogMessage("main", "Could not load sessions list from Bolt databases with error message: "+err.Error())
		logging_service.WriteLogMessage("main", "Verification API service stopped")
		return
	}
	//----------------------------------------------------------
	//запускаем http-сервер приложения
	h := &http.Server{Addr: ":" + port, Handler: &SecureProxy{}}
	go func() {
		h.ListenAndServe()
	}()

	//создаём канал для сигнала оповещения об остановке работы дочерних сервисов
	shutdownChannel := make(chan bool)
	//создаём синхронизированную группу дочерних сервисов
	waitGroup := &sync.WaitGroup{}
	//добавляем в синхронизированную группу
	waitGroup.Add(services_count)
	//запускаем сервис очистки просроченных сессий пользователей
	go Cleaner(StrToInt(sessions_cleaner_frequency), StrToInt(max_session_timeout), shutdownChannel, waitGroup)
	go logging_service.WriteLogMessage("main", "Verification API service started at http port "+port)
	fmt.Println("Verification API service started at http port " + port)
	//--------------------------------------------------------
	//создаём канал для сигнала оповещения об остановке работы приложения
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	//ждём получения сигнала остановки
	<-quitChannel
	//посылаем сигнал остановки дочерним сервисам
	for i := 1; i <= services_count; i++ {
		shutdownChannel <- true
	}
	//ждём закрытия всех приложений
	waitGroup.Wait()
	//останавливаем веб-приложение
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	h.Shutdown(ctx)
	fmt.Println("Verification API service stopped at http port " + port)
	logging_service.WriteLogMessage("main", "Verification API service stopped at http port "+port)
}
