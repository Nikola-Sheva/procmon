package monitor

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	usersMutex sync.Mutex
)

func loadUsers() {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	data, err := os.ReadFile("users.json")
	if err == nil {
		json.Unmarshal(data, &users)
	} else {
		log.Println("users.json не найден, будет создан новый")
	}
}

func saveUsers() {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Printf("Ошибка при сохранении users.json: %v", err)
		return
	}
	os.WriteFile("users.json", data, 0644)
}

func GetUsers() map[string]int64 {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	// Возвращаем копию, чтобы не менять оригинал
	copy := make(map[string]int64)
	for k, v := range users {
		copy[k] = v
	}
	return copy
}

func AddUser(username string, chatID int64) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	users[username] = chatID
	saveUsers()
}

// monitor/users.go или auth.go
func RemoveUser(username string) bool {
	if _, ok := users[username]; ok {
		delete(users, username)
		saveUsers()
		return true
	}
	return false
}

// ListUsers возвращает список всех usernames
func ListUsers() []string {
	userList := make([]string, 0, len(users))
	for username := range users {
		userList = append(userList, username)
	}
	return userList
}
