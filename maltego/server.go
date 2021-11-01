package maltego

/*
   Gondor - Go Maltego Transform Framework
   Copyright (C) 2021 Maxime Landon

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// TransformServer - A server holding all its registered Transforms,
// serving them, either through HTTP or through local invocation.
type TransformServer struct {
	config     *globalConfig         // A configuration applying to all transforms, or the server.
	transforms map[string]*Transform // All user-registered transforms
	hs         http.Server
	mux        *http.ServeMux
	mutex      *sync.RWMutex // Concurrency
}

// NewTransformServer - Create a new Transform Server instance,
// optionally passing a Maltego configuration file (for global
// transform settings, HTTP security details, etc)
func NewTransformServer(config *globalConfig) *TransformServer {
	ts := &TransformServer{
		config: config,
		hs:     http.Server{},
		mux:    http.NewServeMux(),
		mutex:  &sync.RWMutex{},
	}

	// Make a default config if needed
	if config == nil {

	}

	return ts
}

// RegisterTransform - Once you have declared/instantiated a Transform
// in your code, you must register it to a Server with this function.
// The path at which the Transform is available is automatically set
// from its properties, and this should match any exported Config.
func (ts *TransformServer) RegisterTransform(t *Transform) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	// Map the transform to the server
	ts.transforms["transform.Namespace"] = t

	// And to the HTTP server
	ts.mux.HandleFunc("transform.Namespace", ts.transformHandler)

	return
}

// ListenAndServe - The Transform Server starts serving its content, pulling from the current
// state of its configuration: target address, TLS configuration, transforms settings, etc.
func (ts *TransformServer) ListenAndServe() (err error) {

	// Bind the mux handler to the server
	ts.hs.Handler = ts.mux

	return
}

// ListenAndServeTLS - The Transform Server starts serving its content, with an optional TLS
// configuration passed as argument. If nil, will default on its present configuration state.
func (ts *TransformServer) ListenAndServeTLS(addr string, tlsConfig *tls.Config) (err error) {

	// Bind the mux handler to the server
	ts.hs.Handler = ts.mux

	return
}

// GetTransform - Find the Transform corresponding to an HTTP URL path.
func (ts *TransformServer) GetTransform(path string) *Transform {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	return ts.transforms[path]
}

// transformHandler - Handle a request to run a Transform from a Maltego Client: unmarshal the Request,
// pass it to a Transform, run the latter and return its output, regardless of the outcome.
func (ts *TransformServer) transformHandler(w http.ResponseWriter, r *http.Request) {

	// Get the transform transform keyed with the request path
	transform := ts.GetTransform(r.URL.Path)
	if transform == nil {
		http.Error(w, "Did not found Transform for required URL path", http.StatusNoContent)
		return
	}

	// Get the request body, and return if failed or empty
	r.ParseForm()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		http.Error(w, "Error: No form Data in Request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the Maltego Request into its type.
	var request = Message{}
	err = xml.Unmarshal(data, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new Transform instance based on the model.
	instance := transform.newInstanceFromRequest(request)

	// Run the transform.
	err = transform.run(instance)

	// Marshal its output (success or failure)
	response, err := instance.marshalOutput(err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Finally, write the output to the HTTP response
	fmt.Fprintf(w, string(response))
}
