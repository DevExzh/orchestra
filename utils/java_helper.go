package utils

import (
	"archive/zip"
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	InvalidJarFile              = errors.New("invalid jar file")
	MainClassNotFoundInManifest = errors.New("main class not found in the manifest file")
	ManifestNotFound            = errors.New("manifest file not found")
	JVMNotFound                 = errors.New("runnable JVM not found. Maybe the JVM version is too low to run the file. ")
	NoAvailableRuntime          = errors.New("no runtime available to run the jar file")
	InvalidClassVersion         = errors.New("invalid class version")
)

// ScanJarFile Scans specific jar file at the given path and then returns the path of main class file
func ScanJarFile(path string) (string, error) {
	if f, err := os.Stat(path); err == nil {
		if ret := f.IsDir() || filepath.Ext(path) != ".jar"; ret {
			return "", InvalidJarFile
		}
	} else {
		return "", err
	}
	reader, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer func(reader *zip.ReadCloser, e *error) {
		er := reader.Close()
		if er != nil {
			e = &er
		}
	}(reader, &err)
	var file fs.File
	file, err = reader.Open("META-INF/MANIFEST.MF")
	if err == nil {
		// Found manifest file
		zr := bufio.NewReader(file) // Read single file inside the zip file
		var mainClass string
	readFileInZip:
		for {
			var str string
			str, err = zr.ReadString('\n')
			switch {
			case err == io.EOF:
				break readFileInZip
			case err != nil:
				return "", err
			}
			str = strings.TrimSpace(str)
			kvs := strings.Split(str, ":")
			if len(kvs) == 2 {
				var kv []string // Slice storing the key-value data
				for _, v := range strings.Split(str, ":") {
					kv = append(kv, strings.TrimSpace(v)) // Trim all spaces
				}
				if kv[0] == "Main-Class" {
					mainClass = kv[1]
				}
			}
		}
		if mainClass == "" {
			return "", MainClassNotFoundInManifest
		}
		if er := file.Close(); er != nil {
			return "", er
		}
		f, e := reader.Open(strings.ReplaceAll(mainClass, ".", "/") + ".class")
		if e != nil {
			return "", e
		}
		classfile, er := os.Create(filepath.Join(
			os.TempDir(),
			strconv.FormatInt(time.Now().Unix(), 10)+mainClass))
		if er != nil {
			return "", er
		}
		_, copyerr := io.Copy(classfile, f)
		if copyerr != nil {
			return "", copyerr
		}
		classfile.Close()
		return classfile.Name(), nil
	}
	return "", ManifestNotFound
}

// IsLegalClassFile Checks if a specified file is a legal class file,
// if the file is legal returns the major version of the file,
// otherwise return 0
func IsLegalClassFile(path string) uint16 {
	f, _ := os.Open(path)
	defer func(file *os.File, p string) {
		f.Close()
		os.Remove(path)
	}(f, path)
	bytes := make([]byte, 4)
	f.Read(bytes)
	magic := binary.BigEndian.Uint32(bytes)
	if magic != 0xcafebabe {
		return 0
	}
	f.Read(make([]byte, 2)) // Ignore the minor version
	bytes = make([]byte, 2)
	f.Read(bytes)
	major := binary.BigEndian.Uint16(bytes)
	if major > 0x2c {
		return major
	} else {
		return 0
	}
}

// ConvertStringToClassVersionNumber Converts the version string to the version number of Java class file
func ConvertStringToClassVersionNumber(version string) uint16 {
	vs := strings.Split(version, ".")
	if len(vs) == 0 {
		return 0
	}
	major, err := strconv.ParseUint(vs[0], 10, 16)
	if err != nil {
		return 0
	}
	switch {
	case major > 4:
		return uint16(major + 0x2c)
	case major == 1:
		if minor, er := strconv.ParseUint(vs[0], 10, 16); er != nil {
			return uint16(minor + 0x2c)
		} else {
			return 0
		}
	default:
		return 0
	}
}

// ConvertClassVersionNumberToString Similar to ConvertStringToClassVersionNumber, but do the opposite operation
func ConvertClassVersionNumberToString(major uint16) string {
	switch major {
	case 0x3F:
		return "Java SE 19"
	case 0x3E:
		return "Java SE 18"
	case 0x3D:
		return "Java SE 17"
	case 0x3C:
		return "Java SE 16"
	case 0x3B:
		return "Java SE 15"
	case 0x3A:
		return "Java SE 14"
	case 0x39:
		return "Java SE 13"
	case 0x38:
		return "Java SE 12"
	case 0x37:
		return "Java SE 11"
	case 0x36:
		return "Java SE 10"
	case 0x35:
		return "Java SE 9"
	case 0x34:
		return "Java SE 8"
	case 0x33:
		return "Java SE 7"
	case 0x32:
		return "Java SE 6.0"
	case 0x31:
		return "Java SE 5.0"
	case 0x30:
		return "JDK 1.4"
	case 0x2F:
		return "JDK 1.3"
	case 0x2E:
		return "JDK 1.2"
	case 0x2D:
		return "JDK 1.1"
	default:
		return "Invalid version"
	}
}

func FindRunnableJVM(version uint16, jvm []RuntimeVersion) (string, error) {
	for _, j := range jvm {
		if version <= ConvertStringToClassVersionNumber(j.Name) {
			return j.Path, nil
		}
	}
	return "", JVMNotFound
}

// IsRunnableJarFile Checks if the specific jar file at the given path is runnable.
// If the jar is runnable, it will return (string, nil), otherwise ("", error)
func IsRunnableJarFile(path string, runtimes []RuntimeVersion) (string, error) {
	if len(runtimes) == 0 {
		return "", NoAvailableRuntime
	}
	cls, err := ScanJarFile(path)
	if err != nil {
		return "", err
	}
	if major := IsLegalClassFile(cls); major == 0 {
		return "", InvalidClassVersion
	} else {
		return FindRunnableJVM(major, runtimes)
	}
}
