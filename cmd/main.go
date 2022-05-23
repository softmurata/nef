/*
// test code for hello.go
package cmd

func Command() string {

	cmd := "git commit -m first commit"

	return cmd
	
}
*/

package main

import (
	"github.com/softmurata/internal/logger"
	// nef_service "github.com/softmurata/nef/pkg/service"
)

// var NEF = &nef_service.NEF{}

func main(){
	/*
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.AppLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()
	*/

	logger.AppLog.Infof("Main Function")

}
