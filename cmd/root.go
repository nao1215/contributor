package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type authorInfo struct {
	name       string
	mail       string
	addLine    int
	deleteLine int
}

var rootCmd = &cobra.Command{
	Use:   "coder",
	Short: "output the ranking table of people who wrote a lot of code (only support git)",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(coder(cmd, args))
	},
}

func exitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// Execute start command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitError(err)
	}
}

func coder(cmd *cobra.Command, args []string) int {
	if !canUseGitCommand() {
		fmt.Fprint(os.Stderr, "coder: this system does not install git command.")
		return 1
	}

	if err := cdGitRootDir(); err != nil {
		fmt.Fprint(os.Stderr, "coder: can not change current directory. are you in the .git project?")
		return 1
	}

	authors, err := authorsInfo()
	if err != nil {
		fmt.Fprint(os.Stderr, "coder: can not get authors information")
		return 1
	}
	printTable(authors)
	return 0
}

func printTable(author []authorInfo) {
	data := [][]string{}
	for _, a := range author {
		data = append(data, []string{a.name, a.mail, strconv.Itoa(a.addLine), strconv.Itoa(a.deleteLine)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Email", "+(append)", "-(delete)"})
	table.SetAutoWrapText(false)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// canUseGitCommand check whether git command install in the system.
func canUseGitCommand() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// cdGitRootDir change current directory to git project root.
func cdGitRootDir() error {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return err
	}
	return os.Chdir(strings.Split(string(out), "\n")[0])
}

// exists check whether file or directory exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return (err == nil)
}

func getAuthorsAlphabeticalOrder() ([]string, error) {
	out, err := exec.Command("git", "log", "--pretty=format:%an<%ae>").Output()
	if err != nil {
		return nil, err
	}

	list := strings.Split(string(out), "\n")
	list = removeDuplicate(list)
	sort.Strings(list)
	return list, nil
}

func authorsInfo() ([]authorInfo, error) {
	authorInfos := []authorInfo{}
	authors, err := getAuthorsAlphabeticalOrder()
	if err != nil {
		fmt.Fprintf(os.Stderr, "coder: %s\n", err.Error())
		return nil, err
	}

	rex := regexp.MustCompile(`<[^<]*@.*>$`) // e-mail address
	for _, v := range authors {
		mailWithAngleBrackets := rex.FindString(v)
		tmp := strings.Replace(mailWithAngleBrackets, "<", "", 1)
		mail := strings.Replace(tmp, ">", "", 1)

		a := authorInfo{
			name: strings.Replace(v, mailWithAngleBrackets, "", 1),
			mail: mail,
		}
		authorInfos = append(authorInfos, a)
	}

	for _, v := range authorInfos {
		out, err := exec.Command("git", "log", "--author="+v.mail, "--numstat", "--pretty=", "--no-merges", "main").Output()
		if err != nil {
			out, err = exec.Command("git", "log", "--author=\""+v.mail+"\"", "--numstat", "--pretty=", "--no-merges", "master").Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "coder: %s\n", err.Error())
				return nil, err
			}
		}

		list := strings.Split(string(out), "\n")
		for _, line := range list {
			list := strings.Fields(line)
			// 0=append line num, 1=delete line num, 2=file name
			if len(list) == 3 {
				add, err := atoi(list[0])
				if err != nil {
					return nil, err
				}
				delete, err := atoi(list[1])
				if err != nil {
					return nil, err
				}
				v.addLine += add
				v.deleteLine += delete
			}
		}
	}
	return sortInOrderOfMostCodesWritten(authorInfos), nil
}

func sortInOrderOfMostCodesWritten(a []authorInfo) []authorInfo {
	// key=author, value=append LOC
	authorMap := map[authorInfo]int{}

	for _, v := range a {
		authorMap[v] = v.addLine
	}

	type kv struct {
		Key   authorInfo
		Value int
	}

	var ss []kv
	for k, v := range authorMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	result := []authorInfo{}
	for _, kv := range ss {
		result = append(result, kv.Key)
	}
	return result
}

func atoi(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprint(os.Stderr, "coder: can not convert line from string to integer")
		return 0, err
	}
	return i, nil
}

// removeDuplicate removes duplicates in the slice.
func removeDuplicate(list []string) []string {
	results := make([]string, 0, len(list))
	encountered := map[string]bool{}
	for i := 0; i < len(list); i++ {
		if !encountered[list[i]] {
			encountered[list[i]] = true
			results = append(results, list[i])
		}
	}
	return results
}
