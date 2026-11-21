package sandbox

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestContainerPolicy(t *testing.T) {
	args := strings.Join(containerArgs("test", "image", "run"), " ")
	for _, required := range []string{"--network none", "--memory 256m", "--memory-swap 256m", "--pids-limit 32", "--user 65534:65534", "--read-only", "--cap-drop ALL", "no-new-privileges", "--ulimit cpu=15:15"} {
		if !strings.Contains(args, required) {
			t.Fatalf("missing %s", required)
		}
	}
	for _, forbidden := range []string{"--privileged", "--env-file", "docker.sock", "--volume", "--mount"} {
		if strings.Contains(args, forbidden) {
			t.Fatalf("unsafe %s", forbidden)
		}
	}
}

// These are real attack submissions. CI must enable them after building the image.
func TestAttacks(t *testing.T) {
	if os.Getenv("SANDBOX_TESTS") != "1" {
		t.Skip("set SANDBOX_TESTS=1 with Docker available")
	}
	for _, test := range []struct{ name, expected string }{
		{"network", "BLOCKED"}, {"cpu", "Time Limit Exceeded"}, {"wall", "Time Limit Exceeded"},
		{"memory", "Memory Limit Exceeded"}, {"process", "BLOCKED"}, {"output", "Output Limit Exceeded"},
		{"user", "BLOCKED"}, {"mounts", "BLOCKED"}, {"secrets", "BLOCKED"},
	} {
		t.Run(test.name, func(t *testing.T) {
			source, err := os.ReadFile("testdata/" + test.name + ".cpp")
			if err != nil {
				t.Fatal(err)
			}
			result, err := (Runner{Wall: 8 * time.Second}).Run(context.Background(), source, "", "run")
			if err != nil {
				t.Fatal(err)
			}
			if test.expected == "BLOCKED" {
				if result.Verdict != "Executed" || strings.TrimSpace(result.Output) != "BLOCKED" {
					t.Fatalf("attack escaped or test malfunctioned: %+v", result)
				}
			} else if result.Verdict != test.expected {
				t.Fatalf("want %s, got %+v", test.expected, result)
			}
			if test.name == "wall" {
				if err := exec.Command("docker", "inspect", result.Container).Run(); err == nil {
					t.Fatal("wall timeout left the container alive")
				}
			}
			if len(result.Output) > 64*1024 {
				t.Fatal("output limit exceeded")
			}
		})
	}
}
