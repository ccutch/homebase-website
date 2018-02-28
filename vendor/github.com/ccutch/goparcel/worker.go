package goparcel

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Worker at workspace does not exist")
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

type Worker struct {
	workspace string
	PublicURL string
	Entry     string
}

func New(entry string) *Worker {
	return &Worker{
		workspace: "/tmp/" + genid(),
		Entry:     entry,
		PublicURL: "/",
	}
}

func Open(workspace, entry string) (*Worker, error) {
	worker := &Worker{
		Entry:     entry,
		workspace: workspace,
		PublicURL: "/",
	}

	if _, err := os.Stat(workspace); os.IsNotExist(err) {
		return worker, &NotFoundError{}
	}

	return worker, nil
}

func genid() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]

	}

	return string(b)
}
func (w *Worker) Close() {
	err := exec.Command("rm", "-rf", w.workspace).Run()
	if err != nil {
		panic(err)
	}
}

func (w Worker) Workspace() string {
	return w.workspace
}

func (w *Worker) EnsureWorkspace(err error) {
	if err != nil && IsNotFoundError(err) {
		if err = w.SetupWorkspace(); err != nil {
			panic(err)
		}
	}
}

func (w *Worker) SetupWorkspace() error {
	err := os.Mkdir(w.workspace, 0777)
	if err != nil {
		return fmt.Errorf("Error create workspace %s %v", w.workspace, err)
	}

	wd, _ := os.Getwd()
	d := path.Join(wd, filepath.Dir(w.Entry))
	fmt.Println(d)

	if _, err = w.RunHere("ln -s " + d + " " + w.workspace); err != nil {
		return fmt.Errorf("Error linking source directory %v", err)
	}

	if err = w.Run("npm init -y"); err != nil {
		return fmt.Errorf("Error setting up npm %v", err)
	}

	if err = w.Run("npm i --save parcel-bundler"); err != nil {
		return fmt.Errorf("Error setting up npm %v", err)
	}
	return nil
}

func (w *Worker) RunHere(command string) ([]byte, error) {
	p := strings.Fields(command)
	cmd := exec.Command(p[0], p[1:]...)
	return cmd.Output()
}

func (w *Worker) Run(command string) error {
	p := strings.Fields(command)
	cmd := exec.Command(p[0], p[1:]...)
	cmd.Dir = w.workspace
	o, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(o)
	go func() {
		for scanner.Scan() {
			t := scanner.Text()
			if t != "" {
				log.Printf("[parcel] %s\n", t)
			}
		}
	}()

	return cmd.Run()
}

func (w *Worker) Start() error {
	err := w.Run("node_modules/.bin/parcel watch --public-url " + w.PublicURL + " " + w.Entry)
	fmt.Println(err)

	return err
}

func (w *Worker) File(name string) (*os.File, error) {
	return os.Open(w.workspace + "/dist/" + name)
}

func (w *Worker) FileServer() http.Handler {
	fs := http.FileServer(http.Dir(w.workspace + "/dist"))
	return http.StripPrefix(w.PublicURL, fs)
}
