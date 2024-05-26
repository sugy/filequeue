package filequeue

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock_filequeue "github.com/sugy/filequeue/mock"
)

// mockCommand helps to mock the exec.Command
func mockCommand(path string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", path}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is used as a helper process to simulate exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i := 0; i < len(args); i++ {
		if args[i] == "--" {
			args = args[i+1:]
			break
		}
	}
	switch args[0] {
	case "echo":
		if len(args) > 1 {
			os.Stdout.WriteString(strings.Join(args[1:], " ") + "\n")
		}
	}
	os.Exit(0)
}

func TestExecuteRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdExec := mock_filequeue.NewMockCommandExecutor(ctrl)
	mockCmdExec.EXPECT().
		Command("echo", "hello", "world").
		Return(mockCommand("echo", "hello", "world"))

	// Replace the global CmdExec with the mock
	CmdExec = mockCmdExec

	tests := []struct {
		name       string
		path       string
		args       []string
		stdin      string
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{
			name:       "Valid echo command",
			path:       "echo",
			args:       []string{"hello", "world"},
			wantStdout: "hello world\n",
			wantStderr: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newExecute(tt.path, tt.args)
			cmd.stdin = tt.stdin

			err := cmd.run()
			if (err != nil) != tt.wantErr {
				t.Fatalf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if cmd.stdout != tt.wantStdout {
				t.Fatalf("stdout = %v, want %v", cmd.stdout, tt.wantStdout)
			}
			if cmd.stderr != tt.wantStderr {
				t.Fatalf("stderr = %v, want %v", cmd.stderr, tt.wantStderr)
			}
		})
	}
}
