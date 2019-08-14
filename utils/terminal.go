package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var timeout = time.Second * 30

type sshTerminal struct {
	Session *ssh.Session
	exitMsg string
	stdout  io.Reader
	stdin   io.Writer
	stderr  io.Reader
}

func (t *sshTerminal) updateTerminalSize() {
	go func() {
		// SIGWINCH is sent to the process when the window size of the terminal has
		// changed.
		sigwinchCh := make(chan os.Signal, 1)
		signal.Notify(sigwinchCh, syscall.SIGWINCH)

		fd := int(os.Stdin.Fd())
		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		for {
			select {
			// The client updated the size of the local PTY. This change needs to occur
			// on the server side PTY as well.
			case sigwinch := <-sigwinchCh:
				if sigwinch == nil {
					return
				}
				currTermWidth, currTermHeight, err := terminal.GetSize(fd)

				// terminal size has not changed, don't do anything.
				if currTermHeight == termHeight && currTermWidth == termWidth {
					continue
				}

				t.Session.WindowChange(currTermHeight, currTermWidth)
				if err != nil {
					fmt.Printf("Unable to send window-change reqest: %s.", err)
					continue
				}

				termWidth, termHeight = currTermWidth, currTermHeight

			}
		}
	}()

}

func (t *sshTerminal) interactiveSession() error {

	defer func() {
		if t.exitMsg == "" {
			fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			fmt.Fprintln(os.Stdout, t.exitMsg)
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	t.updateTerminalSize()

	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return err
	}
	t.stderr, err = t.Session.StderrPipe()

	go io.Copy(os.Stderr, t.stderr)
	go io.Copy(os.Stdout, t.stdout)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n > 0 {
				_, err = t.stdin.Write(buf[:n])
				if err != nil {
					fmt.Println(err)
					t.exitMsg = err.Error()
					return
				}
			}
		}
	}()

	err = t.Session.Shell()
	if err != nil {
		return err
	}
	err = t.Session.Wait()
	if err != nil {
		return err
	}
	return nil
}

// RunTerminal run terminal in interactive mode
func RunTerminal(company string, c *ssh.Client) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	s := sshTerminal{
		Session: session,
	}
	SuccessPrintln(fmt.Sprintf("登陆%q的cds设备", company))
	return s.interactiveSession()
}

// RunCommands run give command
func RunCommands(c *ssh.Client, cmd string) error {
	return nil
}

//CopyFromFile Copies the contents of an os.File to a remote location, it will get the length of the file by looking it up from the filesystem
func CopyFromFile(c *ssh.Client, file os.File, remotePath string, permissions string) error {
	stat, _ := file.Stat()
	return Copy(c, &file, remotePath, permissions, stat.Size(), timeout)
}

// Copy Copies the contents of an io.Reader to a remote location
func Copy(c *ssh.Client, r io.Reader, remotePath string, permissions string, size int64, timeout time.Duration) error {
	filename := path.Base(remotePath)
	directory := path.Dir(remotePath)
	session, err := c.NewSession()

	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer wg.Done()
		w, err := session.StdinPipe()
		if err != nil {
			errCh <- err
			return
		}
		defer w.Close()

		_, err = fmt.Fprintln(w, "C"+permissions, size, filename)
		if err != nil {
			errCh <- err
			return
		}

		_, err = io.Copy(w, r)
		if err != nil {
			errCh <- err
			return
		}

		_, err = fmt.Fprint(w, "\x00")
		if err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer wg.Done()
		err := session.Run(fmt.Sprintf("%s -qt %s", "scp", directory))
		if err != nil {
			errCh <- err
			return
		}
	}()

	if waitTimeout(&wg, timeout) {
		return errors.New("timeout when upload files")
	}
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
