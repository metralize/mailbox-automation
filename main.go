package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mvinturis/mailbox-automation/aol"
	"github.com/mvinturis/mailbox-automation/common/chromeuser"
	"github.com/mvinturis/mailbox-automation/common/models"

	"github.com/chromedp/chromedp"
)

const maxInstances = 1

func main() {
	rand.Seed(time.Now().UnixNano())

	waitChan := make(chan bool)
	work := make(chan *models.Seed, maxInstances)

	for i := 0; i < maxInstances; i++ {
		go func() {
			for {
				job := <-work
				startActivities(job)
			}
		}()
	}

	i := 0
	lines, err := readLines("seeds-aol.txt")
	if err != nil {
		fmt.Println("[ERROR] Error: %s", err)
		return
	}

	for _, line := range lines {
		i++
		tokens := strings.Split(line, ":")

		job := &models.Seed{
			Email:        tokens[0],
			Password:     tokens[1],
			RecoveryCode: "",
			LocalEmail:   "",
			ProxyIp:      "",
		}

		work <- job
		time.Sleep(time.Duration(6) * time.Second)
	}

	<-waitChan
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func startActivities(seed *models.Seed) {

	seed.ProfilePath = chromeuser.SetProfile(seed)

	// remove Headless option
	opts := append(chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.UserDataDir(seed.ProfilePath),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-web-security", "1"),
	)
	if seed.ProxyIp != "" {
		opts = append(opts,
			chromedp.Flag("proxy-server", "socks5://"+seed.ProxyIp+":1080"),
		)
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	var params models.TaskParams
	params.Keyword = "aol"

	runner := aol.NewRunner(seed, taskCtx)

	runner.Start("createNewSeed", &params)
}
