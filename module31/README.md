# skillbox-module31

31.5 Практическая работа

Напишите HTTP-сервис, который принимает входящие соединения с JSON-данными и обрабатывает их.
Запрос должен возвращать ID пользователя и статус 201.

Сделайте обработчик, который делает друзей из двух пользователей.
Запрос должен возвращать статус 200 и сообщение «username_1 и username_2 теперь друзья».

Сделайте обработчик, который удаляет пользователя.
Запрос должен возвращать 200 и имя удалённого пользователя.

Сделайте обработчик, который возвращает всех друзей пользователя.
После /friends/ указывается id пользователя, друзей которого мы хотим увидеть.

Сделайте обработчик, который обновляет возраст пользователя.
Запрос должен возвращать 200 и сообщение «возраст пользователя успешно обновлён».

Отрефакторьте приложение так, чтобы вы могли поднять две реплики данного приложения.
Используйте любую базу данных (например, MongoDB), чтобы сохранять информацию о пользователях, или можете сохранять информацию в файл, предварительно сереализуя в JSON.
Напишите proxy или используйте, например, nginx.
Протестируйте приложение.