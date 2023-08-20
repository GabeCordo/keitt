package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango/components/cluster"
	"github.com/GabeCordo/mango/core"
	"github.com/GabeCordo/mango/utils"
	"net/http"
	"time"
)

type ModuleRequestBody struct {
	ModulePath string `json:"path,omitempty"`
	ModuleName string `json:"module,omitempty"`
}

func (thread *Thread) moduleCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		thread.postModuleCallback(w, r)
	} else if r.Method == "DELETE" {
		thread.deleteModuleCallback(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (thread *Thread) postModuleCallback(w http.ResponseWriter, r *http.Request) {

	request := &ModuleRequestBody{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = common.RegisterModule(thread.C1, thread.ProvisionerResponseTable, request.ModulePath, thread.Config.Timeout)

	if errors.Is(err, utils.NoResponseReceived) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err != nil {
		w.WriteHeader(http.StatusConflict)
	}

	response := core.Response{Success: err == nil}
	if err != nil {
		response.Description = err.Error()
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

func (thread *Thread) deleteModuleCallback(w http.ResponseWriter, r *http.Request) {

	request := &ModuleRequestBody{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = common.DeleteModule(thread.C1, thread.ProvisionerResponseTable, request.ModuleName, thread.Config.Timeout)

	if errors.Is(err, utils.NoResponseReceived) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := core.Response{Success: err == nil}
	if err != nil {
		response.Description = err.Error()
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

type JSONResponse struct {
	Status      int    `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	Data        any    `json:"data,omitempty"`
}

type SupervisorConfigJSONBody struct {
	Module     string            `json:"module"`
	Cluster    string            `json:"cluster"`
	Config     cluster.Config    `json:"common"`
	Supervisor uint64            `json:"id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type SupervisorProvisionJSONResponse struct {
	Cluster    string `json:"cluster,omitempty"`
	Supervisor uint64 `json:"id,omitempty"`
}

func (thread *Thread) supervisorCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		thread.postSupervisorCallback(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (thread *Thread) postSupervisorCallback(w http.ResponseWriter, r *http.Request) {

	var request SupervisorConfigJSONBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if (r.Method != "GET") && (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = common.SupervisorProvision(thread.C1, thread.ProvisionerResponseTable,
		request.Module, request.Cluster, request.Metadata, &request.Config, thread.Config.Timeout)

	if errors.Is(err, utils.NoResponseReceived) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	response := core.Response{Success: err == nil}
	if err != nil {
		response.Description = err.Error()
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}

type DebugJSONBody struct {
	Action string `json:"action"`
}

type DebugJSONResponse struct {
	Duration time.Duration `json:"time-elapsed"`
	Success  bool          `json:"success"`
}

func (thread *Thread) debugCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		thread.getDebugCallback(w, r)
	} else if r.Method == "POST" {
		thread.postDebugCallback(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (thread *Thread) getDebugCallback(w http.ResponseWriter, r *http.Request) {
	// do nothing
}

func (thread *Thread) postDebugCallback(w http.ResponseWriter, r *http.Request) {

	request := &DebugJSONBody{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if (r.Method != "OPTIONS") && err != nil {
		fmt.Println("missing body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := core.Response{Success: true}

	if request.Action == "shutdown" {
		common.ShutdownCore(thread.Interrupt)
	}

	b, _ := json.Marshal(response)
	w.Write(b)
}
