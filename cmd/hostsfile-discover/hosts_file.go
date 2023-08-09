package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"sort"
	"strings"
)

type hostfile_entry struct {
	Ip   net.IP
	Name string
}

func (h hostfile_entry) String() string {
	return fmt.Sprintf("%s=%s", h.Name, h.Ip.String())
}

func getEntriesMatching(hostname string) ([]hostfile_entry, error) {
	entries, err := getHostFileEntries()
	if err != nil {
		return []hostfile_entry{}, err
	}

	hostRegExp := regexp.MustCompile("^[^\\.]*\\." + hostname + "\\.?$")
	matching := filter(entries, func(h hostfile_entry) bool { return hostRegExp.Match([]byte(h.Name)) })
	sort.Slice(matching, func(i, j int) bool {
		return matching[i].Name < matching[j].Name
	})

	return matching, nil
}

func getHostFileEntries() ([]hostfile_entry, error) {
	content, err := ioutil.ReadFile(CONFIG.HOSTS_FILE_PATH())

	if err != nil {
		return []hostfile_entry{}, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	var entries []hostfile_entry

	for scanner.Scan() {
		entriesFromLine := processHostfileLine(CONFIG.TLD(), scanner.Text())
		entries = append(entries, entriesFromLine...)

	}

	return entries, nil
}

func processHostfileLine(tld string, line string) []hostfile_entry {
	fields := strings.Fields(line)

	if len(fields) < 2 {
		return []hostfile_entry{}
	}

	ipAddress := net.ParseIP(fields[0])

	if ipAddress == nil || ipAddress.IsLoopback() || ipAddress.IsMulticast() {
		return []hostfile_entry{}
	}

	var entries []hostfile_entry
	for i := 1; i < len(fields); i++ {
		name := fields[i]

		tldRegExp := regexp.MustCompile("\\.?" + tld + "$")
		if tldRegExp.Match([]byte(name)) {
			entries = append(entries, hostfile_entry{Ip: ipAddress, Name: name})
		}
	}
	return entries
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
