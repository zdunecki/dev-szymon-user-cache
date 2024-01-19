package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/dev-szymon/user-cache/cache"
	"github.com/dev-szymon/user-cache/service"
)

var usersCount int = 100
var concurrentRequestsCount int = 10000

func main() {
	userService := service.NewUserService(usersCount)

	cachedUserService := cache.NewCache(userService)

	var wg sync.WaitGroup
	for i := 0; i < concurrentRequestsCount; i++ {
		wg.Add(1)
		userId := fmt.Sprintf("user_%d", i%usersCount+1)

		go func(id string) {
			defer wg.Done()
			user, err := cachedUserService.GetOne(userId)
			if err != nil {
				log.Fatalf("error getting user %s", id)
			}
			fmt.Printf("Found user: %s\n", user.Id)
		}(userId)
	}

	wg.Wait()

	// number of db hits will be more than `usersCount` due to sync.RWMutex usage,
	// however should still be less than `concurrentRequestCount`
	fmt.Printf("database service reached %d times\n", userService.DbHits)
}
