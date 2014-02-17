package util

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateTarGz(output io.Writer, dir string) error {
	if doesExist, _ := Exists(dir); !doesExist {
		return errors.New("Source directory does not exist")
	}

	fmt.Println("- Archiving and compressing dir:", dir)

	// Buffer the writer
	buf := bufio.NewWriter(output)
	defer buf.Flush()

	// Gzip on the way through
	gzw := gzip.NewWriter(buf)
	defer gzw.Close()

	// Tar the directory
	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Add each file to the tar archive
	visit := func(path string, f os.FileInfo, err error) error {
		// Skip directories
		if f.Mode().IsDir() {
			return nil
		}

		// Build archive from our root directory
		new_path := path[len(filepath.Dir(dir)):]
		if len(new_path) == 0 {
			return nil
		}

		// Open the current file
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		// Build header
		if h, err := tar.FileInfoHeader(f, new_path); err != nil {
			log.Fatalln(err)
		} else {
			h.Name = new_path
			if err = tw.WriteHeader(h); err != nil {
				log.Fatalln(err)
			}
		}

		// Copy file to tar achive
		if length, err := io.Copy(tw, fr); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf("--> %s (%d)\n", fr.Name(), length)
		}

		return nil
	}

	// Recursively add files in root directory
	filepath.Walk(dir, visit)

	fmt.Println("")
	return nil
}

func ExtractTarGz(dest string, file io.Reader, stripTopDir bool) error {
	target, err := filepath.Abs(dest)

	if doesExist, _ := Exists(target); !doesExist {
		os.MkdirAll(target, 0755)
	}

	if err != nil {
		return err
	}

	if file == nil {
		return errors.New("Unable to extract empty file")
	}

	// Buffer the reader
	buf := bufio.NewReader(file)

	// UnGzip on the way through
	gzr, _ := gzip.NewReader(buf)
	defer gzr.Close()

	// Get a tar reader
	tr := tar.NewReader(gzr)

	fmt.Println("- Extracting to", target, "...")
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		hdrName := hdr.Name
		if stripTopDir {
			strip := hdr.Name[0 : strings.Index(hdr.Name[1:len(hdr.Name)], string(os.PathSeparator))+1]
			hdrName = strings.TrimPrefix(hdr.Name, strip)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			dir := filepath.Join(target, hdrName)
			fmt.Printf(">>> %s\n", dir)

			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
			continue

		case tar.TypeReg, tar.TypeRegA:
			// ok
		default:
			fmt.Fprintf(os.Stderr, "**error: %v\n", hdr.Typeflag)
			return err
		}
		oname := filepath.Join(target, hdrName)
		fmt.Printf("--> %s (%d)\n", oname, hdr.Size)
		dir := filepath.Dir(oname)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		o, err := os.OpenFile(oname,
			os.O_WRONLY|os.O_CREATE,
			hdr.FileInfo().Mode(),
		)
		if err != nil {
			return err
		}
		defer o.Close()

		_, err = io.Copy(o, tr)
		if err != nil {
			return err
		}
		o.Sync()
		o.Close()
	}

	fmt.Println("")
	return nil
}
