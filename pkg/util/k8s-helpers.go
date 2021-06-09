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

package util

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flacatus/che-inspector/pkg/common/client"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// CopyArtifactsFromPod copy files from a given pod to local directory where che-inspector it is executed
func CopyArtifactsFromPod(srcPath string, destPath string, podName string, namespace string, containerName string) error {
	restconfig, err, coreclient := client.InitRestClient()
	reader, outStream := io.Pipe()
	//todo some containers failed : tar: Refusing to write archive contents to terminal (missing -f option?) when execute `tar cf -` in container
	cmdArr := []string{"tar", "cf", "-", srcPath}
	req := coreclient.RESTClient().
		Get().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   cmdArr,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
	if err != nil {
		log.Fatalf("error %s\n", err)
		return err
	}
	go func() {
		defer outStream.Close()
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
			Tty:    false,
		})
	}()
	prefix := getPrefix(srcPath)
	prefix = path.Clean(prefix)
	prefix = stripPathShortcuts(prefix)
	destPath = path.Join(destPath, path.Base(prefix))
	err = untarAll(reader, destPath, prefix)
	return err
}

func ExecInContainer(podName string, container string, namespace string, command string) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}
	_, outStream := io.Pipe()
	restconfig, _, coreclient := client.InitRestClient()
	req := coreclient.RESTClient().Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		Param("container", container).
		Param("stdout", "true").
		Param("stderr", "true")

	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: outStream,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		return err
	}

	return nil
}

// Comment
func untarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		evaledPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evaledPath, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func getPrefix(file string) string {
	return strings.TrimLeft(file, "/")
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
func stripPathShortcuts(p string) string {
	newPath := path.Clean(p)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}
