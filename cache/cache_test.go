package cache

import (
	"fmt"
	"sync"
	"testing"
)

type TestUser struct {
	Id string
}

type mockedUserService struct {
	dbReads int
	db      map[string]*TestUser
}

func (s *mockedUserService) GetOne(key string) (*TestUser, error) {
	s.dbReads++
	user, ok := s.db[key]
	if !ok {
		return nil, fmt.Errorf("user %s not found", key)
	}
	return user, nil
}

func TestCache(t *testing.T) {
	tt := []struct {
		concurrentRequestsCount int
		usersCount              int
	}{
		{
			concurrentRequestsCount: 1000,
			usersCount:              100,
		},
		{
			concurrentRequestsCount: 10000,
			usersCount:              50,
		},
	}

	for _, tc := range tt {
		testName := fmt.Sprintf("%d users, %d requests", tc.concurrentRequestsCount, tc.usersCount)
		t.Run(testName, func(t *testing.T) {
			service := &mockedUserService{db: map[string]*TestUser{}}

			// seeed the database
			for i := 0; i < tc.usersCount; i++ {
				userId := fmt.Sprintf("user_%d", i+1)
				service.db[userId] = &TestUser{Id: userId}
			}

			serviceWithCache := NewCache[TestUser](service)

			// warm the cache
			for i := 0; i < tc.usersCount; i++ {
				userId := fmt.Sprintf("user_%d", i%tc.usersCount+1)
				user, err := serviceWithCache.GetOne(userId)
				if err != nil {
					t.Errorf("%s: error getting user %s\n", testName, userId)
				}

				if user.Id != userId {
					t.Errorf("%s: want %s user, got %s\n", testName, userId, user.Id)
				}
			}

			// perform concurrent requests
			var wg sync.WaitGroup
			for i := 0; i < tc.concurrentRequestsCount; i++ {
				wg.Add(1)
				userId := fmt.Sprintf("user_%d", i%tc.usersCount+1)

				go func(id string) {
					defer wg.Done()

					user, err := serviceWithCache.GetOne(id)
					if err != nil {
						t.Errorf("%s: error getting user %s\n", testName, id)
					}

					if user.Id != id {
						t.Errorf("%s: want %s user, got %s\n", testName, id, user.Id)
					}
				}(userId)
			}

			wg.Wait()

			if service.dbReads != tc.usersCount {
				t.Errorf("%s: got %d reads, want %d reads\n", testName, service.dbReads, tc.usersCount)
			}
		})
	}
}

func TestCacheWithoutWarm(t *testing.T) {
	tt := []struct {
		concurrentRequestsCount int
		usersCount              int
	}{
		{
			concurrentRequestsCount: 1000,
			usersCount:              100,
		},
		{
			concurrentRequestsCount: 10000,
			usersCount:              50,
		},
	}

	for _, tc := range tt {
		testName := fmt.Sprintf("%d users, %d requests", tc.concurrentRequestsCount, tc.usersCount)
		t.Run(testName, func(t *testing.T) {
			service := &mockedUserService{db: map[string]*TestUser{}}

			// seeed the database
			for i := 0; i < tc.usersCount; i++ {
				userId := fmt.Sprintf("user_%d", i+1)
				service.db[userId] = &TestUser{Id: userId}
			}

			serviceWithCache := NewCache[TestUser](service)

			// perform concurrent requests
			var wg sync.WaitGroup
			for i := 0; i < tc.concurrentRequestsCount; i++ {
				wg.Add(1)
				userId := fmt.Sprintf("user_%d", i%tc.usersCount+1)

				go func(id string) {
					defer wg.Done()

					user, err := serviceWithCache.GetOne(id)
					if err != nil {
						t.Errorf("%s: error getting user %s\n", testName, id)
					}

					if user.Id != id {
						t.Errorf("%s: want %s user, got %s\n", testName, id, user.Id)
					}
				}(userId)
			}

			wg.Wait()

			if service.dbReads != tc.usersCount {
				t.Errorf("%s: got %d reads, want %d reads\n", testName, service.dbReads, tc.usersCount)
			}
		})
	}
}
