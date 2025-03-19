# ozonadv
Консольная утилита извлечения данных по рекламным кампаниям Озон.

```shell
ozonadv [command] [options]

Команды:
    stat        Формирование и загрузка статистики по кампаниям
    stat:info   Статус формирования отчетов
    stat:pull   Получить незагруженные отчеты
```

Справка по команде

```shell
ozonadv [command] --help
```

Пример команды фомирования и загрузки статистики

```shell
ozonadv stat --from-date 2025-01-01 --to-date 2025-01-02
```

Документация Озон Performance Api

https://docs.ozon.ru/api/performance/
