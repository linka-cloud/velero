/*
Copyright the Velero contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package backend

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/kopia/kopia/repo/blob"
)

// Store defines the methods for Kopia to establish a connection to
// the backend storage
type Store interface {
	// Setup setups the variables to a specific backend storage
	Setup(ctx context.Context, flags map[string]string, logger logrus.FieldLogger) error

	// Connect connects to a specific backend storage with the storage variables
	Connect(ctx context.Context, isCreate bool, logger logrus.FieldLogger) (blob.Storage, error)
}
