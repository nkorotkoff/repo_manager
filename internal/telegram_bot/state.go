package telegram_bot

import (
	"encoding/json"
	"log"
	"os"
	"repo_manager/internal/structs"
)

var (
	devEnvironmentsContainer []structs.DevEnvironment
	stateOfUsers             map[int64]*commandsStatus
)

func getRepo(repoName string, devName string) structs.Repository {
	for _, devEnvironment := range devEnvironmentsContainer {
		if devEnvironment.Name == devName {
			for _, repo := range devEnvironment.Repositories {
				if repo.Name == repoName {
					return repo
				}
			}
		}
	}

	return structs.Repository{}
}

func getState(chatID int64) *commandsStatus {
	if stateOfUsers == nil {
		stateOfUsers = make(map[int64]*commandsStatus)
	}

	state, exists := stateOfUsers[chatID]
	if !exists {
		initialState := commandsStatus{}
		stateOfUsers[chatID] = &initialState
		updatedState, _ := stateOfUsers[chatID]
		return updatedState
	}
	return state
}

func ParseDevEnvironments() {
	file, err := os.Open("data.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&devEnvironmentsContainer); err != nil {
		log.Fatal(err)
	}
}

func getRepositories(selectedDevEnvironment string) []string {
	var names []string
	var selectedDev structs.DevEnvironment
	for _, devEnvironment := range devEnvironmentsContainer {
		if selectedDevEnvironment == devEnvironment.Name {
			selectedDev = devEnvironment
		}
	}

	for _, repo := range selectedDev.Repositories {
		names = append(names, repo.Name)
	}

	return names
}

func isDevExist(devName string) bool {
	for _, devEnvironment := range devEnvironmentsContainer {
		if devEnvironment.Name == devName {
			return true
		}
	}
	return false
}

func isRepoExist(devName, repoName string) bool {
	for _, repo := range getRepositories(devName) {
		if repoName == repo {
			return true
		}
	}
	return false
}
