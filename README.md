# Организация клиент-серверного взаимодейсвтия между Golang и PostgreSQL

Цель работы — получение первичных навыков в организации долгосрочного хранения данных с использованием PostgreSQL и Golang

В рамках данной лабораторной работы предлагается продолжить изучение Golang и познакомиться с набором стандартных библиотек, используемых для организации клиент-серверного взаимодействия между Golang и Postgresql, где в роли клиента выступает сервис Golang, а в роли сервера СУБД Postgresql.
 
## Основные теоретические сведения

Как установить PostgreSQL: https://tproger.ru/articles/osnovy-postgresql-dlya-nachinayushhih--ot-ustanovki-do-pervyh-zaprosov-250851

Как работать с PostgreSQL из-под Golang: https://golangdocs.com/golang-postgresql-example

Данная лабораторная работа базируется на результатах 6-й лабораторной работы, где было реализовано 3 сервиса: `count`, `hello` и `query`. Поскольку данные сервисы достаточно лаконичны (имеют крайне маленький объем кодовой базы), будем их далее называть МИКРОсервисами.
 
## Порядок выполнения

Для успешного выполнения лабораторной работы необходимо проделать следующие шаги:

0. Установить и настроить PostgreSQL
1. Ознакомиться с теоретическими сведениями
2. Сделать форк данного репозитория в GitHub, склонировать получившуюся копию локально, создать от мастера ветку dev и переключиться на неё
3. Перекопировать код сервисов, полученный в ходе выполнения 6-й лабораторной работы, в соответствующие поддиректории в директории cmd (кроме кода сервиса `hello`, т.к. он уже реализован в качестве примера)
4. Доработать сервисы таким образом, чтобы они использовали для хранения данных СУБД PostgreSQL. Каждый сервис должен как добавлять новые данные в БД (`insert`/`update`), так и доставать их для предоставления пользователю (`select`)
5. Проверить свой код линтерами с помощью команды `make lint`
6. Сделать отчёт и поместить его в директорию docs
7. Зафиксировать изменения, сделать коммит и отправить получившееся состояние ветки дев в личный форк данного репозитория в GitHub
8. Через интерфейс GitHub создать Pull Request dev --> master
9. На защите лабораторной работы продемонстрировать открытый Pull Request. PR должен быть направлен в master ветку форка, а не исходного репозитория

## Содержание отчета

1. Титульный лист
2. Цель работы
3. Задание
4. Ход работы со скриншотами и листингами результатов 
5. Заключение
6. Ответы на контрольные вопросы (если они есть)
7. Список использованных источников
