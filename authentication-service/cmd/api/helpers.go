package main

import (
	"log"
)

func (authServer *AuthService) background(fn func()) {
	authServer.wg.Add(1)
	go func() {
		defer authServer.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s", err)
			}
		}()

		fn()
	}()
}
