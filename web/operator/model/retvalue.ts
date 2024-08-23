//значение, возвращаемое методами API
class Retvalue {
    //флаг успешности операции
    public Success: boolean | undefined;
    //сообщение от API (если есть)
    public Message: string | undefined;
    //данные возвращаемые API
    public Data: any;    
    //тип сообщения, возвращаемого  API
    //0 - ошибка
    //1 - информация
    public Level: number | undefined;
    //инициализация полей при создании класса
    protected static SInit = (() => {
        Retvalue.prototype.Level = 0;
    })();
}

export { Retvalue };