package main

import (
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"os/exec"
	"strings"
)

var (
	etcdURL = ""
	etcdEnv = ""
	debug   = flag.Bool("debug", false, "Debug flag")

	etcdClient = etcd.Client{}
)

func init() {
	flag.StringVar(&etcdURL, "etcdurl", "", "URL for etcd, e.g. http://192.168.59.103:4001")
	flag.StringVar(&etcdEnv, "etcdenv", "", "Name of the etcdEnv")

	flag.Parse()

	if etcdURL == "" {
		etcdURL = os.Getenv("ETCD_URL")
	}

	if etcdEnv == "" {
		etcdEnv = os.Getenv("ETCD_ENV")
	}

	etcdClient = *etcd.NewClient([]string{etcdURL})
}

func main() {
	m1, err1 := readDir(&etcdClient, "/environment")
	m2, err2 := readDir(&etcdClient, fmt.Sprintf("/%s/environment", etcdEnv))

	env := make(map[string]string)

	if err1 == nil {
		env = m1
	}

	if err2 == nil {
		env = merge(env, m2)
	}

	err := run(arrayify(env))
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

	vars := flag.Args()
	var args = make([]string, 0)

	if len(vars) > 1 {
		args = vars[1:len(vars)]
	}

	cmd := exec.Command(vars[0], args...)
	cmd.Env = osEnv
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error starting '%+v'.\n%+v\n", vars[1], err)
		return err
	}

	return nil
}

func arrayify(m map[string]string) []string {
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

func readDir(etcdClient *etcd.Client, directory string) (map[string]string, error) {
	response, err := etcdClient.Get(directory, false, false)

	if err != nil {
		if *debug {
			fmt.Printf("No %s directory found\n", directory)
		}
		return nil, err
	}

	result := make(map[string]string)

	for _, node := range response.Node.Nodes {
		var path = strings.SplitAfter((*node).Key, "/")
		var key = path[len(path)-1]
		var value = (*node).Value

		result[key] = value
	}

	return result, nil
}
