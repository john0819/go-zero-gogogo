package main

import (
	"context"
	"fmt"
	"sync"
)

var synMap sync.Map

type UserSearchSystem interface {
	Search(ctx context.Context, query string) bool
}

type UserSearchSystemImpl struct {
	AdminName string
}

func (u *UserSearchSystemImpl) Search(ctx context.Context, query string) bool {
	full_query := fmt.Sprintf("%s + %s", u.AdminName, query)

	fmt.Printf("full_query: %s\n", full_query)
	return true
}

func main() {
	UserSearchSystem := &UserSearchSystemImpl{
		AdminName: "john",
	}
	synMap.LoadOrStore("user_search_system", UserSearchSystem)

	val, _ := synMap.Load("user_search_system")
	userSearchSystem := val.(*UserSearchSystemImpl)
	userSearchSystem.Search(context.Background(), "admin")
}
