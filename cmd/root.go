package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/nao1215/contributor/internal/completion"
	"github.com/nao1215/contributor/internal/print"
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
	Use:   "contributor",
	Short: "output the ranking table of people who wrote a lot of code (only support git)",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(contributor(cmd, args))
	},
}

// Execute start command.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	completion.DeployShellCompletionFileIfNeeded(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		print.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().BoolP("markdown", "m", false, "Change output format to markdown")
	rootCmd.Flags().BoolP("file", "f", false, "Generate Contributors.md file")
}

func contributor(cmd *cobra.Command, args []string) int {
	markdownFlag, err := cmd.Flags().GetBool("markdown")
	if err != nil {
		print.Err("can not parse command line argument (--markdown)")
		return 1
	}

	fileFlag, err := cmd.Flags().GetBool("file")
	if err != nil {
		print.Err("can not parse command line argument (--file)")
		return 1
	}

	if !canUseGitCommand() {
		print.Err("this system does not install git command.")
		return 1
	}

	if err := cdGitRootDir(); err != nil {
		print.Err("can not change current directory. are you in the .git project?")
		return 1
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
	s.Start()
	authors, err := authorsInfo()
	if err != nil {
		print.Err("can not get authors information")
		return 1
	}
	s.Stop()

	if fileFlag {
		f, err := os.OpenFile("Contributors.md", os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			print.Err(err)
			return 1
		}
		printTable(f, authors, true)

		if err := f.Close(); err != nil {
			print.Err(err)
			return 1
		}
		print.Info("Generate Contributors.md")
	} else {
		printTable(os.Stdout, authors, markdownFlag)
	}
	return 0
}

func printTable(out *os.File, author []authorInfo, markdown bool) {
	tableData := [][]string{}
	names := []string{}
	emails := []string{}
	for _, a := range author {
		if contains(names, a.name) || contains(emails, a.mail) {
			continue
		}
		tableData = append(tableData, []string{a.name, a.mail, strconv.Itoa(a.addLine), strconv.Itoa(a.deleteLine)})
		names = append(names, a.name)
		emails = append(emails, a.mail)
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"Name", "Email", "+(append)", "-(delete)"})
	table.SetAutoWrapText(false)
	if markdown {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	}

	for _, v := range tableData {
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
		print.Err(err)
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

	defalutBranch, err := defaultBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "contributor: %s\n", err.Error())
		return nil, err
	}

	result := []authorInfo{}
	for _, v := range authorInfos {
		out, err := exec.Command("git", "log", "--author="+v.mail, "--numstat", "--pretty=", "--no-merges", defalutBranch).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "contributor: %s\n", err.Error())
			return nil, err
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
				v.addLine = v.addLine + add
				v.deleteLine = v.deleteLine + delete
			}
		}
		result = append(result, v)
	}
	return sortInOrderOfMostCodesWritten(result), nil
}

func defaultBranch() (string, error) {
	out, err := exec.Command("git", "remote", "show", "origin").Output()
	if err != nil {
		return "", errors.New("can not get default branch name")
	}

	list := strings.Split(string(out), "\n")
	for _, v := range list {
		v = strings.TrimSpace(v)
		if strings.Contains(v, "HEAD branch:") {
			v = strings.TrimLeft(v, "HEAD branch:")
			return strings.TrimSpace(v), nil
		}
	}
	return "", errors.New("can not get default branch name")
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
	if s == "-" {
		return 0, nil // this case is binary upload.
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, "contributor: can not convert line from string to integer")
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

// contains returns whether the specified data is contained in the slice.
func contains(slice interface{}, elem interface{}) bool {
	rvList := reflect.ValueOf(slice)

	if rvList.Kind() == reflect.Slice {
		for i := 0; i < rvList.Len(); i++ {
			item := rvList.Index(i).Interface()
			if !reflect.TypeOf(elem).ConvertibleTo(reflect.TypeOf(item)) {
				continue
			}
			target := reflect.ValueOf(elem).Convert(reflect.TypeOf(item)).Interface()
			if ok := reflect.DeepEqual(item, target); ok {
				return true
			}
		}
	}
	return false
}
