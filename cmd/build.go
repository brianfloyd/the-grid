package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"slices"
)

const TAILWIND_CLI_VERSION string = "v3.4.17"
const TEMPL_VERSION string = "v0.2.793"

var SUPPORTED_OS = []string{
	"windows",
	"darwin",
	"linux",
}

var SUPPORTED_ARCH = []string{
	"arm64",
	"amd64",
}

const TEMPL_DOWNLOAD_TEMPLATE = "https://github.com/a-h/templ/releases/download/%s/%s"
const TAILWIND_CLI_DOWNLOAD_TEMPLATE = "https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s"

type BuildSystemInfo struct {
	os           string
	architecture string
}

func main() {
	info, err := getBuildSystemInfo()
	if err != nil {
		return
	}

	err = ensureDirectories()
	if err != nil {
		return
	}

	templ, err := ensureTemplExists(info)
	if err != nil {
		return
	}

	tailwind, err := ensureTailwindExists(info)
	if err != nil {
		return
	}

	gopath, err := exec.LookPath("go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not locate go path. %s\n", err)
		return
	} else {
		fmt.Fprintf(os.Stdout, "Located go binary at %s\n", gopath)
	}

	fmt.Println("Generating tailwind css file.")
	cmd, err := runCommand(tailwind, "-o", "static/css/tailwind-gen.css")
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while waiting for tailwind to complete. %s\n", err)
		return
	}

	fmt.Println("Generating templ templates.")
	cmd, err = runCommand(templ, "generate")
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while waiting for templ generate to complete. %s\n", err)
		return
	}

	fmt.Println("Building application.")
	cmd, err = runCommand(gopath, "build", "-o", "the-grid", "cmd/main.go")
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while waiting for go build to complete. %s\n", err)
		return
	}

	fmt.Println("BUILD SUCCESS")
}

func runCommand(path string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred while starting program %s execution %v. %s\n", path, args, err)
		return nil, errors.New("could not start command")
	}
	return cmd, nil
}

func ensureTailwindExists(info BuildSystemInfo) (string, error) {
	tailwindPath := fmt.Sprintf("./build/toolchain/tailwind/%s", TAILWIND_CLI_VERSION)
	exists, err := doesFileExist(tailwindPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not validate if the tailwind toolchain version exists.\n")
		return "", err
	}

	if exists {
		exists, err = doesFileExist(tailwindPath + "/tailwind")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not validate if the tailwind binary exists.\n")
			return "", err
		}
	}

	if !exists {
		err := getTailwind(tailwindPath, info)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Fprintf(os.Stdout, "Tailwind executable exists at %s.\n", tailwindPath)
	}

	return tailwindPath + "/tailwind", nil
}

func ensureTemplExists(info BuildSystemInfo) (string, error) {
	templPath := fmt.Sprintf("./build/toolchain/templ/%s", TEMPL_VERSION)
	exists, err := doesFileExist(templPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not validate if the templ toolchain version exists.\n")
		return "", err
	}

	if exists {
		exists, err = doesFileExist(templPath + "/templ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not validate if the templ binary exists.\n")
			return "", err
		}
	}

	if !exists {
		err := getTempl(templPath, info)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Fprintf(os.Stdout, "Templ executable exists at %s.\n", templPath)
	}

	return templPath + "/templ", nil
}

func getTempl(path string, info BuildSystemInfo) error {
	fmt.Println("Downloading templ...")
	err := ensureDirectory(path, "templ version")
	if err != nil {
		return err
	}

	archive, err := downloadTempl(path, info)
	if err != nil {
		return err
	}

	file, err := os.Open(archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open templ tar.gz for reading. %s\n", err)
		return err
	}

	err = extractTarGz(path, file)
	if err != nil {
		return err
	}

	executablePath := path + "/templ"
	err = os.Chmod(executablePath, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not make the templ executable have exec permissions. %s\n", err)
		return err
	}
	return nil
}

func getTailwind(path string, info BuildSystemInfo) error {
	fmt.Println("Downloading tailwind...")
	err := ensureDirectory(path, "tailwind version")
	if err != nil {
		return err
	}

	path, err = downloadTailwind(path, info)
	if err != nil {
		return err
	}

	err = os.Chmod(path, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not make the tailwind executable have exec permissions. %s\n", err)
		return err
	}
	return nil
}

func downloadTempl(path string, info BuildSystemInfo) (string, error) {
	osType := ""
	arch := ""
	if info.os == "windows" {
		osType = "Windows"
	} else if info.os == "darwin" {
		osType = "Darwin"
	} else if info.os == "linux" {
		osType = "Linux"
	} else {
		fmt.Fprintf(os.Stderr, "Templ does not support the os '%s'.\n", info.os)
		return "", errors.New("unsupported templ os")
	}

	if info.architecture == "amd64" {
		arch = "x86_64"
	} else if info.architecture == "386" {
		arch = "i386"
	} else if info.architecture == "arm64" {
		arch = "arm64"
	} else {
		fmt.Fprintf(os.Stderr, "Templ does not support the arch '%s'.\n", info.architecture)
		return "", errors.New("unsupported templ arch")
	}

	archive := fmt.Sprintf("templ_%s_%s.tar.gz", osType, arch)
	url := fmt.Sprintf(TEMPL_DOWNLOAD_TEMPLATE, TEMPL_VERSION, archive)
	fmt.Printf("Templ download url: %s\n", url)

	path = path + "/" + archive
	err := download(path, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not download templ at path '%s' and url '%s'.\n", path, url)
		return "", errors.New("could not download templ")
	}

	return path, nil
}

func downloadTailwind(path string, info BuildSystemInfo) (string, error) {
	osType := ""
	arch := ""
	extension := ""
	if info.os == "windows" {
		osType = "windows"
		extension = ".exe"
	} else if info.os == "darwin" {
		osType = "macos"
	} else if info.os == "linux" {
		osType = "linux"
	} else {
		fmt.Fprintf(os.Stderr, "Tailwind does not support the os '%s'.\n", info.os)
		return "", errors.New("unsupported tailwind os")
	}
	if info.architecture == "amd64" {
		arch = "x64"
	} else if info.architecture == "arm64" {
		arch = "arm64"
	} else {
		fmt.Fprintf(os.Stderr, "Tailwind does not support the arch '%s'.\n", info.architecture)
		return "", errors.New("unsupported tailwind arch")
	}

	executable := fmt.Sprintf("tailwindcss-%s-%s%s", osType, arch, extension)
	url := fmt.Sprintf(TAILWIND_CLI_DOWNLOAD_TEMPLATE, TAILWIND_CLI_VERSION, executable)
	fmt.Printf("Tailwind download url: %s\n", url)

	path = path + "/tailwind"

	err := download(path, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not download tailwind at path '%s' and url '%s'.\n", path, url)
		return "", errors.New("could not download tailwind")
	}

	return path, nil
}

func download(path, url string) error {
	out, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not prepare file '%s' for download.\n", path)
		return errors.New("could not prepare file for download")
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not download file from '%s' to '%s'. %s\n", url, path, err)
		return errors.New("could not download file")
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not copy response body from '%s' to '%s'. %s\n", url, path, err)
		return errors.New("could not copy response body")
	}

	if n == 0 {
		fmt.Fprintf(os.Stderr, "Failed to copy enough response body bytes from '%s' to '%s'.\n", url, path)
		return errors.New("could not copy enough bytes")
	}

	return nil
}

func getBuildSystemInfo() (BuildSystemInfo, error) {
	osType, arch := runtime.GOOS, runtime.GOARCH
	if !slices.Contains(SUPPORTED_OS, osType) {
		fmt.Fprintf(os.Stderr, "Unsupported build platform os %s\n.", osType)
		return BuildSystemInfo{}, errors.New("unsupported build platform os")
	}

	if !slices.Contains(SUPPORTED_ARCH, arch) {
		fmt.Fprintf(os.Stderr, "Unsupported build platform architecture %s\n.", arch)
		return BuildSystemInfo{}, errors.New("unsupported build platform architecture")
	}

	return BuildSystemInfo{
		os:           osType,
		architecture: arch,
	}, nil
}

func ensureDirectories() error {
	err := ensureDirectory("./build", "build")
	if err != nil {
		return err
	}
	err = ensureDirectory("./build/toolchain", "toolchain")
	if err != nil {
		return err
	}
	err = ensureDirectory("./build/toolchain/templ", "templ")
	if err != nil {
		return err
	}
	err = ensureDirectory("./build/toolchain/tailwind", "tailwind")
	if err != nil {
		return err
	}
	return nil
}

func ensureDirectory(directory string, name string) error {
	exists, err := doesFileExist(directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not check if the %s directory exists. %s\n", name, err)
		return err
	}

	if !exists {
		err := os.Mkdir(directory, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create the %s directory. %s\n", name, err)
			return err
		}
	}
	return nil
}

func doesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func extractTarGz(path string, gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open gzip reader for the stream. %s\n", err)
		return errors.New("could not open gzip reader for the stream")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "extract tar gz next() failed: %s", err)
			return errors.New("next() failed")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "extract tar gz mkdir() failed: %s", err)
				return errors.New("mkdir() failed")
			}
		case tar.TypeReg:
			outFile, err := os.Create(path + "/" + header.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "extract tar gz create() failed: %s", err)
				return errors.New("create() failed")
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				fmt.Fprintf(os.Stderr, "extract tar gz copy() failed: %s", err)
				return errors.New("copy() failed")
			}
			outFile.Close()
		default:
			fmt.Fprintf(os.Stderr, "extract tar gz unknown type: %s in %s.\n", header.Typeflag, header.Name)
			return errors.New("unkown type")
		}
	}
	return nil
}
