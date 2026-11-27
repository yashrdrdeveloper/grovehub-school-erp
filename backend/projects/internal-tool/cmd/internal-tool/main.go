package main

import (
 "log"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/launcher"
)

func main() {
 if err := launcher.Run(); err != nil { log.Fatal(err) }
}
