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

package clog

import (
	"github.com/sirupsen/logrus"
)

// LOGGER is a globally configured clog
var LOGGER = logrus.New()

// Comment
func init() {
	LOGGER.Formatter = new(logrus.TextFormatter)                                     // Default
	LOGGER.Formatter.(*logrus.TextFormatter).FullTimestamp = true                    // Enable timestamp
	LOGGER.Formatter.(*logrus.TextFormatter).TimestampFormat = "2006-01-02 15:04:05" // Customize timestamp format
	LOGGER.Level = logrus.TraceLevel
}
