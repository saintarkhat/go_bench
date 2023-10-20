package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
type Browser struct {
	Name string
	Seen bool
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile("@")
	seenBrowsers := []Browser{}
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	users := make([]map[string]interface{}, 0)
	for _, line := range lines {
		user := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}
			if strings.Contains(browser, "Android") {
				isAndroid = true
				found := false
				for j := range seenBrowsers {
					if seenBrowsers[j].Name == browser {
						found = true
						break
					}
				}
				if !found {
					seenBrowsers = append(seenBrowsers, Browser{Name: browser, Seen: true})
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				found := false
				for j := range seenBrowsers {
					if seenBrowsers[j].Name == browser {
						found = true
						break
					}
				}
				if !found {
					seenBrowsers = append(seenBrowsers, Browser{Name: browser, Seen: true})
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", uniqueBrowsers)
}
