// Copyright (c) 2021 The Jaeger Authors.
// //
// // Copyright (c) 2021 Red Hat, Inc.
// // This program and the accompanying materials are made
// // available under the terms of the Eclipse Public License 2.0
// // which is available at https://www.eclipse.org/legal/epl-2.0/
// //
// // SPDX-License-Identifier: EPL-2.0
// //
// // Contributors:
// //   Red Hat, Inc. - initial API and implementation
// //

package version

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestVersion(t *testing.T) {
	cmd := AddVersionCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.Execute()
	_, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
}
