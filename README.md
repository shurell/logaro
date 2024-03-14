# logaro
simple log service
1. Using sqlite for log storage
2. Web interface for logviews and to configure log sources, etc

простой сервис логов
1. Для хранения логов используется sqlite (github.com/mattn/go-sqlite3)
2. Вебинтерфейс для удобного просмотра, редактирования источников и типов сообщений

для регистрации сообщения от приложения нужно отправить сервису POST-запрос, содержащий JSON следующего вида:

{
    "log_owner_token": "0000-1111-AAAA-FFFF",
    "message_type": 6,
    "source_utime": 342342345,
    "message": "Это сообщение информационное"
}

где log_owner_token - токен предварительно зарегистрированного приложения,
    message_type - соответственно Id  предварительно зарегистрированного типа сообщения (например 6 - Инфо),
    source_utime - unix-время события, типа Int64
    message - текст сообщения в кодировке utf-8
