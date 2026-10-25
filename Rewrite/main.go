package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	incus "github.com/lxc/incus/client"
	"github.com/lxc/incus/shared/api"
)

type ServerCreationRequest struct {
	Name     string
	MaxRAM   uint
	MaxCPU   uint
	MaxDisk  uint
	MaxProcs uint
}

func die(source string, err error) {
	log.Fatalf("%s: %s\n", source, err.Error())
}

func methodLock(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return r.Method == method
}

func main() {
	client, err := incus.ConnectIncusUnix("", nil)
	if err != nil {
		die("connectIncus", err)
	}

	http.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		if !methodLock(w, r, "GET") {
			return
		}
		fmt.Fprint(w, "1.0.0_alpha")
	})
	http.HandleFunc("/api/createServer", func(w http.ResponseWriter, r *http.Request) {
		if !methodLock(w, r, "POST") {
			return
		}
		var req ServerCreationRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, `{"success": false, "reason": "invalid or malformed JSON"}`, http.StatusBadRequest)
			return
		}

		cores := fmt.Sprintf("%d", int(math.Ceil(float64(req.MaxCPU)/100.0))) // keep req.maxCPU a float or Go explodes
		cpu_cfs_string := fmt.Sprintf("%dms/100ms", req.MaxCPU)
		ram_max_string := fmt.Sprintf("%dMiB", req.MaxRAM)
		disk_max_string := fmt.Sprintf("%dMiB", req.MaxDisk)
		proc_max_string := fmt.Sprintf("%d", req.MaxProcs)
		incusReq := api.InstancesPost{
			Name: req.Name,
			Type: api.InstanceTypeContainer,
			Source: api.InstanceSource{
				Type:     "image",
				Alias:    "ubuntu/26.04",
				Server:   "https://images.linuxcontainers.org",
				Protocol: "simplestreams",
			},
			Config: map[string]string{
				"limits.cpu":           cores,
				"limits.cpu.allowance": cpu_cfs_string,
				"limits.memory":        ram_max_string,
				"limits.processes":     proc_max_string,
			},
			// 2. Disk Size goes into Devices under the root block device
			Devices: map[string]map[string]string{
				"root": {
					"type": "disk",
					"path": "/",
					"pool": "default",
					"size": disk_max_string,
				},
			},
		}
		op, err := client.CreateInstance(incusReq)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success": false, "reason": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		err = op.Wait()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"success": false, "reason": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, `{"success": true}`)
	})
	http.ListenAndServe(":8080", nil)
}
