package main

type contextKey string

const (
	authUserKey               = "authenticatedUserID"   // key used for an authenticated user in Session Manager
	isAuthenticatedContextKey = contextKey(authUserKey) // custom type wrapping authUserKey string
)
