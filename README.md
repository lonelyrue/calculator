# Calculator Service
 
Это веб-сервис для вычисления арифметических выражений. Пользователь может отправить арифметическое выражение по HTTP и получить его результат.
 
**Установка:**
 

*   Для запуска проекта вам потребуется Go. Убедитесь, что у вас установлен Go и добавлен в PATH.
     
*   Запуск проекта
     
    Вы можете запустить проект с помощью следующей команды bash:
    
        go run ./cmd/calc_service/...
    
    **Примеры использования:**
     
*   Успешный запрос bash:
    
        curl --location 'localhost:8080/api/v1/calculate' \
    
        --header 'Content-Type: application/json' \
    
        --data '{
        "expression": "2 + 2 * 2" }'
    
    **Ответ:**
     
        json {
         "result": "6"
        }
    
*   Ошибка 422 (Неверное выражение):

        curl --location 'localhost:8080/api/v1/calculate' \
        --header 'Content-Type: application/json' \
        --data '{
          "expression": "2 + a"
        }'
    
    **Ответ:**
     
        json
        {
          "error": "Expression is not valid"
        }
    
    *   Ошибка 500 (Внутренняя ошибка сервера) :
         
            curl --location 'localhost:8080/api/v1/calculate' \
            --header 'Content-Type: application/json' \
            --data '{
              "expression": ""
            }'
        
        **Ответ:**
        
            json
            {
              "error": "Internal server error"
            }
