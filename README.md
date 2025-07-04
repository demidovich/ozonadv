# ozonadv

<a href="https://github.com/demidovich/ozonadv/releases"><img src="https://img.shields.io/github/release/demidovich/ozonadv.svg" alt="Latest Release"></a> [![Build status](https://github.com/demidovich/ozonadv/workflows/release/badge.svg)](https://github.com/demidovich/ozonadv/actions/workflows/release.yml) [![Build status](https://github.com/demidovich/ozonadv/workflows/develop/badge.svg)](https://github.com/demidovich/ozonadv/actions/workflows/develop.yml)

Консольная утилита извлечения статистики по объектам рекламных кампаний Ozon.

- [Для чего это?](#для-чего-это)
- [Установка](#установка)
- [Интерфейс](#интерфейс)
- [Создание отчета](#создание-отчета)
- [Реализация](#реализация)
- [Мотивация](#мотивация)

## Для чего это?

На текущий момент (апрель 2025) в UI рекламного кабинета Ozon не существует функционала, позволяющего выгрузить статистику по рекламным объектам нескольких кампаний. Если у вас десяток кампаний это неприятно, но терпимо. Если их количество переваливает за сотню и вам нужны еженедельные отчеты - у вас проблема.

ozonadv позволяет запомнить учетные данные рекламного кабинета и создавать на их основе запросы на выгрузки отчетов по одной или нескольким кампаниями. После запуска запроса в работу утилита будет постепенно создавать запросы по рекламным кампаниям, ждать их готовности, скачивать и создавать следующие. После того, как отчеты всех рекламных кампаний будут загружены они могут быть экспортированы в общие csv-файлы (в перспективе будут добавлены другие типы экспорта).

## Установка

Загрузить последнюю версию приложения.

- [macOS](https://github.com/demidovich/ozonadv/releases/latest/download/ozonadv_Darwin_x86_64.tar.gz)
- [Linux](https://github.com/demidovich/ozonadv/releases/latest/download/ozonadv_Linux_x86_64.tar.gz)
- [Windows](https://github.com/demidovich/ozonadv/releases/latest/download/ozonadv_Windows_x86_64.zip)

Распаковать и запустить в терминале.

## Интерфейс

Управления осуществляется через клавиатуру. Текущие описания клавиатурных комбинаций отображаются в нижней части экрана. 

> [!TIP]
> Сочетания клавиш для разных элементов интерфейса могут отличаться. 
> Следите за подсказкой в нижней части экрана.

Примеры интерфейса для Linux версии приложения.

Клавное меню<br>
<img src="./docs/assets/home-menu.png" alt="Главное меню">

Выбор рекламного кабинета<br>
<img src="./docs/assets/cabinets-menu.png" alt="Раздел Рекламных кабинетов">

Состояние отчета<br>
<img src="./docs/assets/stat-info.png" alt="Состояние отчета">

## Создание отчета

Создание отчета происходит из меню рекламного кабинета.

Для создания отчета необходимо определить его параметры:
* Название
* Тип статистики
* Начало интервала, дата
* Конец интервара, дата
* Группировка
* Список рекламных кампаний

Поиск рекламных кампаний для нового отчета<br>
<img src="./docs/assets/stat-new-choose-campaigns.png" alt="Выбор рекламных кампаний для нового отчета">

> [!TIP]
> Модерация кампании выполняется не только при ее создании, но и при редактировании. При фильтрации рекламных кампаний необходимо учитывать то, что на модерацию кампания может быть отправлена из любого активного состояния.

При поиске кампаний для отчета будут проигнорированы кампании, время работы которых не попадает в интервал отчета. Например:

```
Интервал:   2025-01-01 — 2025-01-31

Кампания 1: 2025-01-05 — 2025-01-10, попадает
Кампания 2: 2025-01-05 — 2025-02-10, попадает
Кампания 3: 2025-01-05 —           , попадает
Кампания 4: 2024-12-01 —           , попадает
Кампания 5: 2024-12-01 — 2025-10-01, попадает
Кампания 6: 2024-01-01 — 2024-01-31, не попадает
Кампания 7: 2025-02-01 — 2025-02-31, не попадает
```

Состояния отчета:

* **Не запускался** - ни один запрос на формирование статистики не был создан
* **В процессе** - не все файлы статистики кампаний скачаны
* **Готов к экспорту** - скачаны все файлы статистики кампаний

## Реализация

В настоящее время для выгрузки ретроспективных данных рекламных кабинетов Ozon действуют следующие ограничения:

1. До получения статистических данных необходимо сформировать запрос на их генерацию.
2. На один рекламный кабинет в единицу времени может обрабатываться 1-5 запросов генерации. Все запросы, выходящие за этот лимит будут получать HTTP 429.
3. Один отчет может формироваться до 10 минут. Чаще всего время генерации составляет 3 минуты. Но, в зависимости нагрузки на базу и количество конкурирующих (соседних) выгрузок в очереди, оно может неконтролируемо увеличиться.

Поскольку для генерации отчета одной рекламной кампании может потребоваться до 10 минут, а одновременное количество генерируемых отчетов небольшое, формирование и загрузка всех данных может занять больше часа. Это означает, что работа приложения может быть прервана в любой момент. Для ее возобновления требуется постоянное хранилище данных. так как утилита будет использоваться не разработчиками, docker с базой данных из доступных вариантов исчезает. По этой же причине исчезает возможность создавать конфигурационные файлы с токенами доступа в api.

В качестве хранилища приложение использует локальный диск. Данные хранятся в директории .ozonadv в домашнем каталоге пользователя в формате json.

Приложение не использует какие-либо модифицирующие запросы api. Весь код обмена с api можно посмотреть в директории проекта internal/infra/ozon.

Рано или поздно рекламные кабинеты Озон будут доработаны, в них появится функционал, позволяющий делать массовую выгрузку ретроспективных данных, а из API исчезнут такие сильные ограничения. Но на текущий момент (апрель 2025) дела обстоят именно так.

## Мотивация

Утилита написана в помощь супруге, которая мучилась с рутинной выгрузкой ретроспективых данных по объектам рекламных кампаний.

В качестве личной мотивации использовано желание написать консольное приложение на Go.
