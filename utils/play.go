package utils

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github/luoyayu/goidx/config"
	"log"
	"os/exec"
	"runtime"
)

func Play(ctx context.Context, url, title string, auth string) {
	if auth == "" {
		auth = config.RootBasicAuth
	}
	header := fmt.Sprintf("--http-header-fields=authorization: %s", auth)
	playArgs := []string{header, "--autofit=640", url}
	//log.Println(strings.Join(playArgs, " "))

	cmd := exec.CommandContext(ctx, "mpv", playArgs...)
	switch runtime.GOOS {
	case "darwin":
	case "windows":
		cmd = exec.CommandContext(ctx, "mpv.exe", playArgs...)
	case "linux":
	}

	//out := &bytes.Buffer{}
	//cmd.Stdout = out

	err := cmd.Start()

	//go func(out *bytes.Buffer) {
	//	tm := time.NewTicker(time.Second * 5)
	//	for {
	//		select {
	//		case <-tm.C:
	//			return
	//		default:
	//			log.Println(out.String())
	//			time.Sleep(time.Second * 1)
	//		}
	//	}
	//}(out)

	select {
	case <-ctx.Done():
		if err == nil {
			StopMpvSafely(int32(cmd.Process.Pid))
		}
		return
	}
}

func StopMpvSafely(MpvPid int32) {
	if MpvPid != -1 {
		p, _ := process.NewProcess(MpvPid)
		if err := p.Kill(); err != nil {
			log.Println(err)
		}
	}
}
