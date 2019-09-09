package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var invalidMarkers = []string{": invalid pseudo-version:", ": invalid version:"}

func main() {
	count := 1
	for count > 0 {
		count = 0

		_, stderr, err := goModDownload()

		if err != nil {
			fmt.Println(stderr)
			lines := strings.Split(stderr, "\n")

			for _, l := range lines {
				l = strings.TrimSpace(l)
				idx := containsAnyMarker(l)
				if idx > -1 {
					count++
					addReplaceToGoMod(l[:idx])
				}
			}
		}
	}
}

func goModDownload() (stdout string, stderr string, err error) {
	fmt.Print("download...")
	cmd := exec.Command("go", "mod", "download")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	fmt.Println("OK")
	stdout = outb.String()
	stderr = errb.String()
	return
}

func containsAnyMarker(s string) int {
	for _, marker := range invalidMarkers {
		idx := strings.Index(s, marker)
		if idx > -1 {
			return idx
		}
	}
	return -1
}

func addReplaceToGoMod(dep string) {

	lib := strings.Split(dep, "@")
	ver := strings.Split(dep, "-")

	rep := fmt.Sprintf("%s=%s@%s", dep, lib[0], ver[len(ver)-1])
	log.Println("replacing: ", rep)
	//fmt.Println("replace:", rep)
	_, err := exec.Command("go", "mod", "edit", "-replace", rep).Output()
	if err != nil {
		log.Println("Error while running go mod edit:", err)
	}

}

func goModTidy() {
	_, err := exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		log.Println("Error while running go mod tidy:", err)
	}
}
