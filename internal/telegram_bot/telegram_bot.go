package telegram_bot

import (
	"fmt"
	"log"
	"repo_manager/internal/repo_manager"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"repo_manager/internal/config"
)

const (
	startMessage            = "Привет! Я бот для управления репозиториями. Используйте /help для получения списка команд."
	unknownCommandMessage   = "Неизвестная команда."
	emptyDevErrorMessage    = "Дев не может быть пустым."
	emptyRepoErrorMessage   = "Репозиторий не может быть пустым."
	emptyBranchErrorMessage = "Ветка не может быть пустой."
	invalidDevErrorMessage  = "Передано неверное имя дева."
	invalidRepoErrorMessage = "Передано неверное имя репозитория."
)

type commandsStatus struct {
	currentDev      string
	currentService  string
	serviceSelected bool
	devSelected     bool
}

type CommandHandler func(bot *tgbotapi.BotAPI, chatID int64, args []string)

var commandHandlers = map[string]CommandHandler{
	"start":        handleStart,
	"help":         handleHelp,
	"get_devs":     handleGetDevs,
	"select_dev":   selectDev,
	"select_repo":  selectRepository,
	"status":       getStatus,
	"back":         handleGoBack,
	"git_checkout": handleGitCheckout,
	"git_pull":     handleGitPull,
	"git_status":   handleGitStatus,
}

func Init(config *config.Config) {
	bot, err := tgbotapi.NewBotAPI(config.TelegramAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	ParseDevEnvironments()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			handleCommand(bot, update.Message)
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	command, args := parseCommand(msg.Text)
	handler, found := commandHandlers[command]
	if !found {
		sendMessage(bot, msg.Chat.ID, unknownCommandMessage)
		return
	}

	handler(bot, msg.Chat.ID, args)
}

func parseCommand(text string) (string, []string) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", nil
	}

	command := strings.TrimPrefix(parts[0], "/")
	args := parts[1:]

	return command, args
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

func handleStart(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	sendMessage(bot, chatID, startMessage)
}

func handleHelp(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	helpMessage := "Доступные команды:\n" +
		"/get_devs - показать список девов\n" +
		"/select_dev имя дева - выбрать дев\n" +
		"/status - показать выбранный дев и репозиторий\n" +
		"/back - отмена операции выбора дева или репозитория"
	sendMessage(bot, chatID, helpMessage)
}

func handleGetDevs(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	var names []string
	for _, devEnvironment := range devEnvironmentsContainer {
		names = append(names, devEnvironment.Name)
	}

	response := fmt.Sprintf("Dev Environments:\n%s", strings.Join(names, "\n"))
	sendMessage(bot, chatID, response)
}

func selectDev(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	if len(args) == 0 {
		sendMessage(bot, chatID, emptyDevErrorMessage)
		return
	}

	devName := args[0]

	if !isDevExist(devName) {
		sendMessage(bot, chatID, invalidDevErrorMessage)
		return
	}

	state := getState(chatID)
	state.devSelected = true
	state.currentDev = devName

	if state.serviceSelected {
		state.serviceSelected = false
		state.currentService = ""
	}

	response := fmt.Sprintf("Выбран дев: %s\nВыполните команду /select_repo для выбора репозитория\nДоступные репозитории:\n%s",
		devName, strings.Join(getRepositories(state.currentDev), "\n"))
	sendMessage(bot, chatID, response)
}

func selectRepository(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	state := getState(chatID)
	if !state.devSelected {
		sendMessage(bot, chatID, "Чтобы выбрать репозиторий, нужно сначала выбрать дева")
		return
	}

	if len(args) == 0 {
		sendMessage(bot, chatID, emptyRepoErrorMessage)
		return
	}

	repoName := args[0]

	if !isRepoExist(state.currentDev, repoName) {
		sendMessage(bot, chatID, invalidRepoErrorMessage)
		return
	}

	state.serviceSelected = true
	state.currentService = repoName

	response := fmt.Sprintf("Выбран репозиторий: %s\n"+
		"Доступны команды: \n/git_pull - обновление ветки\n/git_status - получение статуса ветки\n/git_checkout название ветки - смена ветки", repoName,
	)
	sendMessage(bot, chatID, response)
}

func getStatus(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	var selectedDev = "Не выбран"
	var selectedRepo = "Не выбран"
	state := getState(chatID)
	if state.devSelected {
		selectedDev = state.currentDev
	}

	if state.serviceSelected {
		selectedRepo = state.currentService
	}

	response := fmt.Sprintf("Дев: %s\nРепозиторий: %s", selectedDev, selectedRepo)
	sendMessage(bot, chatID, response)
}

func handleGoBack(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	state := getState(chatID)
	var response string

	if state.serviceSelected {
		state.serviceSelected = false
		state.currentService = ""
		response = "Операция с выбором репозитория отменена."
	} else if state.devSelected {
		state.devSelected = false
		state.currentDev = ""
		response = "Операция с выбором дева отменена."
	} else {
		response = "Нет активных операций для отмены."
	}

	sendMessage(bot, chatID, response)
}

func handleGitCheckout(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	if len(args) == 0 {
		sendMessage(bot, chatID, emptyBranchErrorMessage)
		return
	}

	var response string

	state := getState(chatID)
	if !state.serviceSelected {
		sendMessage(bot, chatID, "Сначала выберите дева и репозиторий")
		return
	}

	branch := args[0]

	repo := getRepo(state.currentService, state.currentDev)
	err := repo_manager.Checkout(repo.Path, branch)
	if err != nil {
		response = "Не удалось сменить ветку"
		log.Println(err)
	} else {
		response = "Ветка была успешно сменена"
	}

	sendMessage(bot, chatID, response)

	repo_manager.ApplyActions(repo)
}

func handleGitPull(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	state := getState(chatID)
	if !state.serviceSelected {
		sendMessage(bot, chatID, "Сначала выберите дев и репозиторий")
		return
	}

	repo := getRepo(state.currentService, state.currentDev)

	err := repo_manager.GitPull(repo.Path)

	var response string

	if err != nil {
		response = "Не удалось обновить ветку"
		log.Println(err)
	} else {
		response = "Ветка была успешно обновлена"
	}

	sendMessage(bot, chatID, response)

	repo_manager.ApplyActions(repo)
}

func handleGitStatus(bot *tgbotapi.BotAPI, chatID int64, args []string) {
	state := getState(chatID)
	if !state.serviceSelected {
		sendMessage(bot, chatID, "Сначала выберите дева и репозиторий")
		return
	}

	repo := getRepo(state.currentService, state.currentDev)

	var result, err = repo_manager.GitStatus(repo.Path)

	if err != nil {
		log.Println(err)
		sendMessage(bot, chatID, "Не удалось выполнить комманду git_status")
	}

	sendMessage(bot, chatID, result)
}
