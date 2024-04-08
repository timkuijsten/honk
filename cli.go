package main

import (
	"fmt"
	"os"
	"strconv"
)

type cmd struct {
	help     string
	callback func(args []string)
}

var commands = map[string]cmd{
	"init": {
		help: "initialize honk",
		callback: func(args []string) {
			initdb()
		},
	},
	"upgrade": {
		help: "upgrade honk",
		callback: func(args []string) {
			upgradedb()
		},
	},
	"version": {
		help: "print version",
		callback: func(args []string) {
			fmt.Println(softwareVersion)
			os.Exit(0)
		},
	},
	"admin": {
		help: "admin interface",
		callback: func(args []string) {
			adminscreen()
		},
	},
	"import": {
		help: "import data into honk",
		callback: func(args []string) {
			if len(args) != 4 {
				errx("import username honk|mastodon|twitter srcdir")
			}
			importMain(args[1], args[2], args[3])
		},
	},
	"export": {
		help: "export data from honk",
		callback: func(args []string) {
			if len(args) != 3 {
				errx("export username destdir")
			}
			export(args[1], args[2])
		},
	},
	"devel": {
		help: "turn devel on/off",
		callback: func(args []string) {
			if len(args) != 2 {
				errx("need an argument: devel (on|off)")
			}
			switch args[1] {
			case "on":
				setconfig("devel", 1)
			case "off":
				setconfig("devel", 0)
			default:
				errx("argument must be on or off")
			}
		},
	},
	"setconfig": {
		help: "set honk config",
		callback: func(args []string) {
			if len(args) != 3 {
				errx("need an argument: setconfig key val")
			}
			var val interface{}
			var err error
			if val, err = strconv.Atoi(args[2]); err != nil {
				val = args[2]
			}
			setconfig(args[1], val)
		},
	},
	"adduser": {
		help: "add a user to honk",
		callback: func(args []string) {
			adduser()
		},
	},
	"deluser": {
		help: "delete a user from honk",
		callback: func(args []string) {
			if len(args) < 2 {
				errx("usage: honk deluser username")
			}
			deluser(args[1])
		},
	},
	"chpass": {
		help: "change password of an account",
		callback: func(args []string) {
			if len(args) < 2 {
				errx("usage: honk chpass username")
			}
			chpass(args[1])
		},
	},
	"follow": {
		help: "follow an account",
		callback: func(args []string) {
			if len(args) < 3 {
				errx("usage: honk follow username url")
			}
			user, err := butwhatabout(args[1])
			if err != nil {
				errx("user %s not found", args[1])
			}
			var meta HonkerMeta
			mj, _ := jsonify(&meta)
			honkerid, flavor, err := savehonker(user, args[2], "", "presub", "", mj)
			if err != nil {
				errx("had some trouble with that: %s", err)
			}
			if flavor == "presub" {
				followyou(user, honkerid, true)
			}
		},
	},
	"unfollow": {
		help: "unfollow an account",
		callback: func(args []string) {
			if len(args) < 3 {
				errx("usage: honk unfollow username url")
			}
			user, err := butwhatabout(args[1])
			if err != nil {
				errx("user not found")
			}

			honkerid, err := gethonker(user.ID, args[2])
			if err != nil {
				errx("sorry couldn't find them")
			}
			unfollowyou(user, honkerid, true)
		},
	},
	"sendmsg": {
		help: "send a raw activity",
		callback: func(args []string) {
			if len(args) < 4 {
				errx("usage: honk sendmsg username filename rcpt")
			}
			user, err := butwhatabout(args[1])
			if err != nil {
				errx("user %s not found", args[1])
			}
			data, err := os.ReadFile(args[2])
			if err != nil {
				errx("can't read file: %s", err)
			}
			deliverate(user.ID, args[3], data)
		},
	},
	"cleanup": {
		help: "clean up stale data from database",
		callback: func(args []string) {
			arg := "30"
			if len(args) > 1 {
				arg = args[1]
			}
			cleanupdb(arg)
		},
	},
	"storefiles": {
		help: "store attachments as files",
		callback: func(args []string) {
			setconfig("usefilestore", 1)
		},
	},
	"storeblobs": {
		help: "store attachments as blobs",
		callback: func(args []string) {
			setconfig("usefilestore", 0)
		},
	},
	"extractblobs": {
		help: "extract blobs to file store",
		callback: func(args []string) {
			extractblobs()
		},
	},
	"unplug": {
		help: "disconnect from a dead server",
		callback: func(args []string) {
			if len(args) < 2 {
				errx("usage: honk unplug servername")
			}
			name := args[1]
			unplugserver(name)
		},
	},
	"backup": {
		help: "backup honk",
		callback: func(args []string) {
			if len(args) < 2 {
				errx("usage: honk backup dirname")
			}
			name := args[1]
			svalbard(name)
		},
	},
	"ping": {
		help: "ping from user to user/url",
		callback: func(args []string) {
			if len(args) < 3 {
				errx("usage: honk ping (from username) (to username or url)")
			}
			name := args[1]
			targ := args[2]
			user, err := butwhatabout(name)
			if err != nil {
				errx("unknown user %s", name)
			}
			ping(user, targ)
		},
	},
	"extractchatkey": {
		help: "extract secret chat key from user",
		callback: func(args []string) {
			if len(args) < 3 || args[2] != "yesimsure" {
				errx("usage: honk extractchatkey [username] yesimsure")
			}
			user, _ := butwhatabout(args[1])
			if user == nil {
				errx("user not found")
			}
			fmt.Printf("%s\n", user.Options.ChatSecKey)
			user.Options.ChatSecKey = ""
			j, err := jsonify(user.Options)
			if err == nil {
				db := opendatabase()
				_, err = db.Exec("update users set options = ? where username = ?", j, user.Name)
			}
			if err != nil {
				elog.Printf("error bouting what: %s", err)
			}
		},
	},
	"run": {
		help: "run honk",
		callback: func(args []string) {
			serve()
		},
	},
	"backend": {
		help: "run backend",
		callback: func(args []string) {
			backendServer()
		},
	},
	"test": {
		help: "run test",
		callback: func(args []string) {
			ElaborateUnitTests()
		},
	},
}
