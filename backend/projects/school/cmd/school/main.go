package main

import (
 "log"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/launcher"
)

func main() {
 if err := launcher.Run(); err != nil { log.Fatal(err) }
}
