package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/fatih/color"
	"github/luoyayu/goidx/api"
	"github/luoyayu/goidx/config"
	"github/luoyayu/goidx/utils"
	S "gopkg.in/abiosoft/ishell.v2"
	"net/url"
	"github.com/BurntSushi/toml"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var WorkDir = "/"
var WorkDrive *api.SDrive
var SharedDrives []*api.SDrive

func WorkDriveNotExist() bool {
	if WorkDrive == nil {
		fmt.Println(utils.COut("please select a drive by `dd`", color.FgRed))
		return true
	}
	return false
}

func init() {
	c := &config.Config{}
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		panic("no config.toml!")
	}
	config.WorkerHost = c.Host
	config.Password = c.Password
	config.UserName = c.User

	config.RootBasicAuth += base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(config.UserName, ":", config.Password)))
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	s := S.New()
	s.SetPrompt(utils.COut(">>> ", color.FgGreen))
	s.Interrupt(func(c *S.Context, count int, input string) {
		cancel()
		os.Exit(0)
	})

	s.Println("google drive index (gavinhty@gmail.com)")

	// dd
	s.AddCmd(&S.Cmd{Name: "dd",
		Aliases: []string{"dd"},
		Func: func(c *S.Context) {
			SharedDrives = api.GetSharedDrives("")
			if len(SharedDrives) == 0 {
				return
			}
			var choiceNames []string
			for i := range SharedDrives {
				choiceNames = append(choiceNames, SharedDrives[i].Name)
			}
			choice := c.MultiChoice(choiceNames, "select a shared drive")
			WorkDrive = SharedDrives[choice]

			c.Println("※ current working drive:", utils.COut(WorkDrive.Name, color.FgRed))
			c.Println()
		},
	})

	// pwd
	s.AddCmd(&S.Cmd{Name: "pwd", Func: func(c *S.Context) {
		if WorkDriveNotExist() {
			return
		}

		c.Println(utils.COut(WorkDrive.Name+":"+WorkDir, color.FgCyan))
		c.Println()
	}})

	// ls
	s.AddCmd(&S.Cmd{Name: "ls",
		Help: "list dir, not supporting show file info!",
		Func: func(c *S.Context) {
			if WorkDriveNotExist() {
				return
			}

			argPath := strings.Join(c.Args, " ")

			if !filepath.IsAbs(argPath) {
				argPath = filepath.Join(WorkDir, argPath)
			}

			exist, files := api.ShowDir(filepath.Clean(argPath)+"/", "", WorkDrive.Id)
			if !exist {
				c.Println(utils.COut("no such dir!", color.FgRed))
				return
			}
			// list file
			for _, f := range files {
				if f.IsFolder {
					c.Println(utils.COut(f.Name, color.FgGreen))
				} else {
					c.Println(utils.COut(f.Name, color.FgBlue))
				}
			}
			c.Println()
		},
	})

	/*s.AddCmd(&S.Cmd{Name: "ll",
		Func: func(c *S.Context) {
			if WorkDrive == nil {
				log.Println("please select a drive by `dd`")
				return
			}
		},
	})*/

	// cd
	s.AddCmd(&S.Cmd{Name: "cd",
		Func: func(c *S.Context) {
			if WorkDriveNotExist() {
				return
			}

			// select mode
			if len(c.Args) == 0 {
				select2dir(c)
				c.Println()
				return
			}

			argPath := strings.Join(c.Args, " ")

			if !filepath.IsAbs(argPath) {
				argPath = filepath.Join(WorkDir, argPath)
			}

			if exist, _ := api.ShowDir(filepath.Clean(argPath)+"/", "", WorkDrive.Id); exist {
				WorkDir = filepath.Clean(argPath)
			} else {
				c.Println(utils.COut("no such dir!", color.FgRed))
			}
		},
	})

	// play
	s.AddCmd(&S.Cmd{
		Name: "play",
		Func: func(c *S.Context) {
			if len(c.Args) == 0 { // select mode
				var choicePlayable []string
				_, files := api.ShowDir(filepath.Clean(WorkDir)+"/", "", WorkDrive.Id)
				var playableFiles []*api.File

				for _, f := range files {
					if f.IsPlayable {
						playableFiles = append(playableFiles, f)
						choicePlayable = append(choicePlayable, f.Name)
					}
				}

				if len(playableFiles) > 0 {
					choice := c.MultiChoice(choicePlayable, "select to play")
					if choice >= 0 && choice < len(playableFiles) {
						selected := playableFiles[choice]
						//log.Println("play for", selected.Name)

						url_ := (&url.URL{
							Scheme:   "https",
							Host:     config.WorkerHost,
							Path:     WorkDir + "/" + selected.Name,
							RawQuery: url.Values{"rootId": []string{WorkDrive.Id}}.Encode(),
						}).String()
						go utils.Play(ctx, url_, selected.Name, "")
						c.Println("You can type `stop` to kill mpv!")
					}
				} else {
					c.Println(utils.COut("nothing is playable!", color.FgRed))
				}
			} else { // play the folder
				// TODO:
			}
		},
	})

	s.AddCmd(&S.Cmd{
		Name: "stop",
		Func: func(c *S.Context) {
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
		},
	})
	s.Run()
}

func select2dir(c *S.Context) {
	_, files := api.ShowDir(filepath.Clean(WorkDir)+"/", "", WorkDrive.Id)
	choiceDirs := []string{".."}
	folderFiles := []*api.File{{Name: ".."}}

	for _, f := range files {
		if f.IsFolder {
			folderFiles = append(folderFiles, f)
			choiceDirs = append(choiceDirs, f.Name)
		}
	}
	choiceDirs = append(choiceDirs, utils.COut(">>> Quit select mode", color.FgRed))

	if len(folderFiles) > 0 {
		choice := c.MultiChoice(choiceDirs, utils.COut("WorkDir: "+WorkDir, color.FgCyan))
		if choice == len(choiceDirs)-1 {
			return
		}

		if choice >= 0 && choice < len(folderFiles) {
			WorkDir = path.Join(WorkDir, choiceDirs[choice])
			select2dir(c)
		}
	}
	return
}
