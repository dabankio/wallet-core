package eth

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ganache-cli --account="0xc356cfe48d1ddcd2320b62553fe8739978d0478e2e940d6c30190e7637f51c76,200"

const ganacheCliPort = 8545

// RunGanacheCli 启动ganache-cli 用以测试,返回杀死ganache-cli的函数
func RunGanacheCli() (func(), error) {
	if cmdIsPortContainsNameRunning(ganacheCliPort, "ganache-cli") {
		return nil, fmt.Errorf("ganache 似乎已经运行在%d端口了,不先杀掉的话数据可能有问题", ganacheCliPort)
	}

	closeChan := make(chan struct{})

	cmd := exec.Command("ganache-cli")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Args)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		fmt.Println("Wait for message to kill ganache-cli")
		<-closeChan
		fmt.Println("Received message,killing ganache-cli regtest")

		if e := cmd.Process.Kill(); e != nil {
			fmt.Println("关闭 ganache-cli 时发生异常", e)
		}
		fmt.Println("关闭 ganache-cli 完成")
		closeChan <- struct{}{}
	}()

	// err = cmd.Wait()
	fmt.Println("等待2.5秒,让 ganache-cli 启动")
	time.Sleep(time.Millisecond * 2500)
	return func() {
		closeChan <- struct{}{}
	}, nil
}

// 进程名包含$name的端口$port是否在运行中
func cmdIsPortContainsNameRunning(port uint, name string) bool {
	if strings.Contains(runtime.GOOS, "linux") {
		checkPortCmd := exec.Command("netstat", "-ntpl")

		cmdPrint := cmdThenPrint(checkPortCmd)
		if strings.Contains(cmdPrint, strconv.Itoa(int(port))) && strings.Contains(cmdPrint, name) {
			return true
		}
		return false
	} else if strings.Contains(runtime.GOOS, "darwin") {
		checkPortCmd := exec.Command("lsof", "-i", "tcp:18443")
		cmdPrint := cmdThenPrint(checkPortCmd)
		if strings.Contains(cmdPrint, strconv.Itoa(int(port))) && strings.Contains(cmdPrint, "bitcoin") {
			return true
		}
		return false
	} else {
		panic("其他平台尚未實現")
	}
}

// 执行cmd,然后将输出打印到控制台并return
func cmdThenPrint(cmd *exec.Cmd) string {
	fmt.Println("[CMD]", cmd.Args)
	stderr, err := cmd.StderrPipe()
	panicIf(err, "Failed to get stderr pip ")

	stdout, err := cmd.StdoutPipe()
	panicIf(err, fmt.Sprintf("Failed to get stdout pipe %v", err))
	err = cmd.Start()
	panicIf(err, fmt.Sprintf("Failed to start cmd %v", err))
	b, err := ioutil.ReadAll(stdout)
	panicIf(err, fmt.Sprintf("Failed to read cmd (%v) stdout, %v", cmd, err))
	out := string(b)
	fmt.Println(out)

	bo, err := ioutil.ReadAll(stderr)
	panicIf(err, "Failed to read stderr")
	out += string(bo)

	cmd.Wait()
	stdout.Close()
	stderr.Close()
	return strings.TrimSpace(out)
}

func panicIf(e error, msg string) {
	if e != nil {
		panic(fmt.Errorf("【ERR】 %s %v", msg, e))
	}
}
