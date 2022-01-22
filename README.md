
## Setup

    minikube start
    kubectl apply -f k8s-env.yaml
    skaffold dev

    # test locks
    kubectl scale deployments.apps -n test consistelancer --replicas=5

    # to test changes in kubernetes api
    kubectl scale deployments.apps -n test test-web2 --replicas=2

## to do

+ Ограничить sa отдельной ролью
+ <s>Конфиг</s>
+ <s>Модель данных редис</s>
+ <s>Синкер конфига(апстримы) из редиса</s>
+ <s>Лок в редисе для воркера apiWatcher, ttl <5s</s>
+ Прокси в отдельной горутине
+ Логика выбор апстрима
+ Хелфчеки
+ Метрики
+ Unit тесты
+ CI pipeline (мультибранч для тестов)
+ redis cluster insted redis single node
+ <s>redis: 2022/01/08 09:08:17 pubsub.go:159: redis: discarding bad PubSub connection: read tcp 172.17.0.13:51324->10.101.37.77:6379: use of closed network connection</s>
+ <s>lock in redis always changing. try write own locker</s>
+ k8s.io/client-go instead http client
+ <s>skaffold local build instead dockerfile. try ko</s>
+ <s>defer close chan</s>
+ <s>change package layout</s>
