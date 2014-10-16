package main

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"os/exec"
	"strings"
)

var (
	etcdURL    = ""
	etcdEnv    = ""
	etcdClient = etcd.Client{}
)

func init() {
	etcdURL = os.Getenv("ETCD_URL")
	etcdEnv = os.Getenv("ETCD_ENV")
	etcdClient = *etcd.NewClient([]string{etcdURL})
}

func main() {
	m1 := readDir(&etcdClient, "/environment")
	m2 := readDir(&etcdClient, fmt.Sprintf("/%s/environment", etcdEnv))

	err := run(arraify(merge(m1, m2)))
	if err != nil {
		// TODO: Figure out the correct exit code
		os.Exit(1)
	}
}

func run(env []string) error {
	osEnv := os.Environ()

	for _, e := range env {
		osEnv = append(osEnv, e)
	}

	vars := os.Args
	var args = make([]string, 0)

	if len(vars) > 2 {
		args = vars[2:len(vars)]
	}

	cmd := exec.Command(vars[1], args...)
	cmd.Env = osEnv
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error starting '%+v'.\n%+v\n", vars[1], err)
		return err
	}

	return nil
}

func arraify(m map[string]string) []string {
	var result = make([]string, 0)
	for key, value := range m {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}

func merge(first, second map[string]string) map[string]string {
	result := first
	for key, value := range second {
		first[key] = value
	}
	return result
}

func readDir(etcdClient *etcd.Client, directory string) map[string]string {
	result := make(map[string]string)

	response, err := etcdClient.Get(directory, false, false)

	if err != nil {
		fmt.Printf("No %s directory found\n", directory)
		return result
	}

	for _, node := range response.Node.Nodes {
		var path = strings.SplitAfter((*node).Key, "/")
		var key = path[len(path)-1]
		var value = (*node).Value

		result[key] = value
	}

	return result
}