## Дефолтный круд для тасок

Запускается из коробки через docker-compose

### Скрины:

Таска извлекается по id
![img.png](img.png)

Таски извлекаются
![img_1.png](img_1.png)

Таска добавляется
[img_2.png](img_2.png)

Проверяем добавление
![img_3.png](img_3.png)

Обновляем таску по id
![img_4.png](img_4.png)

Проверяем обновление
![img_5.png](img_5.png)

Удаляем таску по id
![img_6.png](img_6.png)

Проверяем удаление
![img_7.png](img_7.png)

### Особенности:
- Есть миграции через golang-migrate (поднимаются в коде при запуске приложения)
- Добавил простенький кеш в redis
- Маршалинг/анмаршалинг делал через генератор easyjson
- Добавил докер для бд/кеша/приложения
- Есть простое логгирование