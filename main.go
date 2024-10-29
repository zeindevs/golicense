package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

var Licenses = []string{
	"Apache-2.0",
	"MIT",
	"ISC",
	"BSD-3-Clause",
	"BSD-2-Clause",
	"BSD-1-Clause",
	"Unlicense",
	"WTFPL",
	"GLWTPL",
}

//go:embed license
var LicenseFS embed.FS

type GitConfig struct {
	Name  string
	Email string
}

func runCmd(cmds ...string) (string, error) {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	text := bytes.TrimSpace(out)
	return string(text), nil
}

func getGitConfig() (*GitConfig, error) {
	name, err := runCmd("git", "config", "--get", "user.name")
	if err != nil {
		return nil, err
	}
	email, err := runCmd("git", "config", "--get", "user.email")
	if err != nil {
		return nil, err
	}
	return &GitConfig{
		Name:  name,
		Email: email,
	}, nil
}

func readLicense(name string) (string, error) {
	text, err := LicenseFS.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func generateLicense(name, email, license string) error {
	text, err := readLicense("license/LICENSE-" + strings.ToUpper(license))
	if err != nil {
		return err
	}
	text = strings.ReplaceAll(text, "{year}", strconv.Itoa(time.Now().Year()))
	text = strings.ReplaceAll(text, "{name}", name)
	if email == "" {
		text = strings.ReplaceAll(text, "{email}", "")
	} else {
		text = strings.ReplaceAll(text, "{email}", fmt.Sprintf("<%s>", email))
	}
	if err := os.WriteFile("LICENSE", []byte(text), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write license: %s", err.Error())
	}
	return nil
}

func promptInput(label string) (string, error) {
	prompt := &promptui.Prompt{
		Label: label,
	}
	value, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return value, nil
}

func promptSelect(label string, items []string) (string, error) {
	prompt := &promptui.Select{
		Label: label,
		Items: items,
	}
	_, value, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return value, nil
}

func errHandler(err error) {
	if err.Error() == "^C" {
		fmt.Println("Error: operation was interrupted by the user")
	} else {
		fmt.Printf("prompt failed %v\n", err)
	}
}

func main() {
	config, err := getGitConfig()
	if err != nil {
		fmt.Printf("failed to get git config %v\n", err.Error())
		return
	}
	name, err := promptSelect("Select your name", []string{config.Name, "No, I will input my name."})
	if err != nil {
		errHandler(err)
		return
	}
	if name == "No, I will input my name." {
		if name, err = promptInput("What is your name"); err != nil {
			errHandler(err)
			return
		}
	}
	email, err := promptSelect("Select your email", []string{config.Email, "No, I will input my email."})
	if err != nil {
		errHandler(err)
		return
	}
	if email == "No, I will input my email." {
		if email, err = promptInput("What is your email"); err != nil {
			errHandler(err)
			return
		}
	}
	license, err := promptSelect("Select License", Licenses)
	if err != nil {
		errHandler(err)
		return
	}
	if err := generateLicense(name, email, license); err != nil {
		fmt.Printf("failed to generate license %v\n", err)
		return
	}
	fmt.Println("Successfully generated license file 🎉")
}
