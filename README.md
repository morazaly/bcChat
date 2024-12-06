для поднятие базы
$env:DB_DSN="kursUser:kursPswd@tcp(127.0.0.1:3306)/TEST?parseTime=true"
goose -dir migrations mysql "$env:DB_DSN" up


//для подключения к вэбсокет
ws://localhost:8080/chat/connect
//тип сообщения
   {
        "room_id": "1",
        "content": "hello"
    }