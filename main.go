/***
Copyright © 2025 Louis Lefebvre <louislefebvre1999@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
***/

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/louislef299/aws-sso/cmd"
	sigs "github.com/louislef299/aws-sso/pkg/os"
)

func main() {
	// Create a context that intercepts SIGINT
	ctx, cancel := signal.NotifyContext(context.Background(), sigs.Signals...)
	go func() {
		<-ctx.Done()
		log.Println("received SIGINT; shutting down...")
		defer func() {
			cancel()
			os.Exit(0)
		}()
	}()

	cmd.Execute(ctx)
}
