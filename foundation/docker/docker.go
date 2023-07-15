package docker

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
)

// Container tracks information about the docker container started for tests.
type Container struct {
	Id   string
	Host string
}

// StartContainer starts the specified container for running tests.
func StartContainer(t *testing.T, img, port string, args ...string) *Container {
	// build args
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, img)

	// execute
	var out bytes.Buffer
	cmd := exec.Command("docker", arg...)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start container %s: %v", img, err)
	}

	// get id of container (just using the first 12 bytes of the container)
	id := out.String()[:12]

	// run docker inspect to retrieve port of the container
	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container %s: %v", id, err)
	}

	// parse json to find port
	var doc []map[string]any
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}
	ip, randPort := extractIpPort(t, doc, port)

	// build container
	c := Container{
		Id:   id,
		Host: net.JoinHostPort(ip, randPort),
	}

	// log info
	t.Logf("image:		%s", img)
	t.Logf("container:	%s", c.Id)
	t.Logf("host:		%s", c.Host)

	return &c
}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T, id string) {
	// stop container
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}
	t.Logf("stopped: %s", id)

	// remove container
	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}
	t.Logf("removed: %s", id)
}

// DumpContainerLogs outputs logs from the running docker container.
func DumpContainerLogs(t *testing.T, id string) {
	// dump logs
	out, err := exec.Command("docker", "logs", id).CombinedOutput()
	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}
	t.Logf("logs for %s\n%s", id, out)
}

// extractIpPort retrives host and port from started container
func extractIpPort(t *testing.T, doc []map[string]any, port string) (string, string) {
	// get network settings
	ns, ok := doc[0]["NetworkSettings"]
	if !ok {
		t.Fatalf("could not get network settings")
	}

	// get ports
	ports, ok := ns.(map[string]any)["Ports"]
	if !ok {
		t.Fatalf("could not get network port settings")
	}

	// get tcp
	tcp, ok := ports.(map[string]any)[port+"/tcp"]
	if !ok {
		t.Fatalf("could not get network tcp port settings")
	}

	// get list
	list, ok := tcp.([]any)
	if !ok {
		t.Fatalf("could not get network tcp port list")
	}

	// parse
	var hostIp, hostPort string
	for _, l := range list {
		data, ok := l.(map[string]any)
		if !ok {
			t.Fatalf("could not get network tcp list data")
		}

		hostIp = data["HostIp"].(string)
		if hostIp != "::" {
			hostPort = data["HostPort"].(string)
		}
	}

	return hostIp, hostPort
}
