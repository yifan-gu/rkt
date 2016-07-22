// Copyright 2015 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/coreos/rkt/api/v1alpha"
	"github.com/coreos/rkt/store/imagestore"
	"golang.org/x/net/context"
)

// v1alphaAPIReadOnlyServer implements v1alpha.PublicWriteOnlyAPI interface.
type v1alphaWriteOnlyAPIServer struct {
	store *imagestore.Store
}

var _ v1alpha.PublicWriteOnlyAPIServer = &v1alphaWriteOnlyAPIServer{}

func newV1alphaWriteOnlyAPIServer(s *imagestore.Store) (*v1alphaWriteOnlyAPIServer, error) {
	return &v1alphaWriteOnlyAPIServer{store: s}, nil
}

func (s *v1alphaWriteOnlyAPIServer) FetchImage(ctx context.Context, request *v1alpha.FetchImageRequest) (*v1alpha.FetchImageResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *v1alphaWriteOnlyAPIServer) RemoveImage(ctx context.Context, request *v1alpha.RemoveImageRequest) (*v1alpha.RemoveImageResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
