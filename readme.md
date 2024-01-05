# Проект управления репозиториями с использованием бота Telegram

Этот проект представляет собой бота Telegram, разработанный на языке программирования Go, для управления репозиториями. Бот взаимодействует с пользователями через команды Telegram, предоставляя функционал по выбору девов, репозиториев, выполнению git-команд.

## Требования

- Установленный Go (версия 1.13 и выше)
- Наличие аккаунта Telegram и созданного бота

## Установка

* **Клонирование репозитория:**

   ```bash
   git clone https://github.com/nkorotkoff/repo_manager
   cd repo_manager
   ```
* **Установка зависимостей**

  ```bash
  go get ./...
   ```
* **Установка зависимостей**

    ```bash
    go run main.go
    ```

## Команды бота

- `/start`: Начало общения с ботом.
- `/help`: Показать список доступных команд.
- `/get_devs`: Показать список девов.
- `/select_dev имя_дева`: Выбрать дев.
- `/status`: Показать выбранный дев и репозиторий.
- `/back`: Отмена операции выбора дева или репозитория.
- `/git_checkout название_ветки`: Сменить ветку в выбранном репозитории.
- `/git_pull`: Обновить ветку в выбранном репозитории.
- `/git_status`: Получить статус ветки в выбранном репозитории.

## Административный интерфейс

Проект также включает административный интерфейс для управления данными девов. Здесь можно просматривать и изменять данные девов, которые бот использует для взаимодействия.