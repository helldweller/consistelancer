
## Setup
    kubectl -n test apply -f k8s-env.yaml
    skaffold dev

## to do

+ Ограничить sa отдельной ролью
* Конфиг
* Модель данных редис
* Синкер конфига(апстримы) из редиса
* Лок в редисе для воркера apiWatcher, ttl <5s
+ Прокси в отдельной горутине
+ Логика выбор апстрима
+ Хелфчеки
+ Метрики
+ redis: 2022/01/08 09:08:17 pubsub.go:159: redis: discarding bad PubSub connection: read tcp 172.17.0.13:51324->10.101.37.77:6379: use of closed network connection