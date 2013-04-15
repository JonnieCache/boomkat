package main

import (
	"archive/zip"
	goopt "github.com/droundy/goopt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

func boomkat(name string, no_cog bool, no_delete bool) {
	usr, err := user.Current()
  if err != nil {
	  log.Fatal(err)
  }

	dirname := path.Base(name)
	dirname = strings.Replace(dirname, "flac_", "", 1)
	dirname = strings.Replace(dirname, "mp3_", "", 1)
	dirname = strings.Replace(dirname, "_", " ", -1)
	dirname = dirname[:len(dirname)-4]
	dirname = strings.TrimRight(dirname, " ")
	path := usr.HomeDir + "/Music/" + dirname

	os.Mkdir(path, 0755)

	zip, err := zip.OpenReader(name)
	if err != nil {
		log.Fatal(err)
	}
	defer zip.Close()

	for _, file_in_zip := range zip.File {
		new_filename := strings.Replace(file_in_zip.Name, "_", " ", -1)
		if new_filename[1] == "-"[0] {
			new_filename = "0" + new_filename
		}
		log.Println(new_filename)
		file_path := path + "/" + new_filename

		file_handler_in_zip, err := file_in_zip.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer file_handler_in_zip.Close()

		file_to_write, err := os.Create(file_path)
		if err != nil {
			log.Fatal(err)
		}
		defer file_to_write.Close()

		_, err = io.Copy(file_to_write, file_handler_in_zip)
		if err != nil {
			log.Fatal(err)
		}
	}
	if !no_delete {
		os.Remove(name)
	}
	if !no_cog {
		cmd := exec.Command("open", path, "-a", "cog")
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = cmd.Wait()
	}
}

func main() {
	var no_cog = goopt.Flag([]string{"-C", "--no-cog"}, []string{"-c", "--cog"},
		"skip opening in cog", "open in cog")
	var no_delete = goopt.Flag([]string{"-D", "--no-delete"}, []string{"-d", "--delete"},
		"skip deleting original zip", "delete original zip")
	goopt.Parse(nil)

	boomkat(goopt.Args[0], *no_cog, *no_delete)

}
