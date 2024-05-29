package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/brancobruyneel/go-fiovb/fiovb"
)

type Mode int

const (
	unknownMode Mode = iota
	readMode
	writeMode
	deleteMode
)

type Args struct {
	mode  Mode
	name  string
	value string
}

var (
	readFlag   = flag.Bool("read", false, "read mode")
	writeFlag  = flag.Bool("write", false, "write mode")
	deleteFlag = flag.Bool("delete", false, "delete mode")
	nameFlag   = flag.String("name", "", "name to write/read")
	valueFlag  = flag.String("value", "", "value to write")
)

func isRootUser() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	return os.Geteuid() == 0 && os.Getuid() == 0 && u.Username == "root"
}

func usageError(err string) {
	flag.Usage()
	log.Fatal(err)
}

func parseReadArgs() *Args {
	if *readFlag {
		if *nameFlag == "" {
			usageError("name is required when in read mode")
			return nil
		}

		return &Args{
			mode: readMode,
			name: *nameFlag,
		}
	}

	return nil
}

func parseWriteArgs() *Args {
	if *writeFlag {
		if *nameFlag == "" {
			usageError("name is required when in write mode")
			return nil
		} else if *valueFlag == "" {
			usageError("value is required when in write mode")
			return nil
		}

		return &Args{
			mode:  writeMode,
			name:  *nameFlag,
			value: *valueFlag,
		}
	}

	return nil
}

func parseDeleteArgs() *Args {
	if *deleteFlag {
		if *nameFlag == "" {
			usageError("name is required when in delete mode")
			return nil
		}

		return &Args{
			mode: deleteMode,
			name: *nameFlag,
		}
	}

	return nil
}

func parseArgs() Args {
	args := Args{mode: unknownMode}

	flag.Parse()

	if writeFlag == nil || readFlag == nil || nameFlag == nil || valueFlag == nil {
		usageError("one of the flags are nil")
		return args
	}

	if *writeFlag && *readFlag && *deleteFlag {
		usageError("write, read and delete modes are mutually exclusive")
		return args
	} else if !*writeFlag && !*readFlag && !*deleteFlag {
		usageError("mode selection required (i.e., --write, --read, --delete)")
		return args
	}

	if writeArgs := parseWriteArgs(); writeArgs != nil {
		return *writeArgs
	} else if readArgs := parseReadArgs(); readArgs != nil {
		return *readArgs
	} else if deleteArgs := parseDeleteArgs(); deleteArgs != nil {
		return *deleteArgs
	}

	return args
}

func read(args *Args) error {
	if args == nil {
		log.Fatal("read: args is nil")
	}

	fvb := fiovb.New()

	if err := fvb.Initialize(); err != nil {
		return err
	}

	value, err := fvb.Read(args.name)
	if err != nil {
		return err
	}

	if err := fvb.Finalize(); err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func write(args *Args) error {
	if args == nil {
		log.Fatal("write: args is nil")
	}

	fvb := fiovb.New()

	if err := fvb.Initialize(); err != nil {
		return err
	}

	if err := fvb.Write(args.name, args.value); err != nil {
		return err
	}

	if err := fvb.Finalize(); err != nil {
		return err
	}

	return nil
}

func delete(args *Args) error {
	if args == nil {
		log.Fatal("delete: args is nil")
	}

	fvb := fiovb.New()

	if err := fvb.Initialize(); err != nil {
		return err
	}

	if err := fvb.Delete(args.name); err != nil {
		return err
	}

	if err := fvb.Finalize(); err != nil {
		return err
	}

	return nil
}

func main() {
	if !isRootUser() {
		log.Fatal("permission denied")
		return
	}

	args := parseArgs()
	if args.mode == unknownMode {
		flag.Usage()
		return
	}

	switch args.mode {
	case readMode:
		if err := read(&args); err != nil {
			log.Fatal(err)
		}
	case writeMode:
		if err := write(&args); err != nil {
			log.Fatal(err)
		}
	case deleteMode:
		if err := delete(&args); err != nil {
			log.Fatal(err)
		}
	}
}
