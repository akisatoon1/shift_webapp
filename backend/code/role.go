package main

const (
	RoleAdmin = 1 << iota
	RoleEmployee
)

func isAdmin(role int) bool {
	return role&RoleAdmin != 0
}
