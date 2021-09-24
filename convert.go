package main

import (
	"archive/zip"
	"log"
	"os"
    "io"
    "strings"
    "strconv"
    "io/fs"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "regexp"
	"path/filepath"
	"github.com/signintech/gopdf"
)

func sortdir(dir []fs.DirEntry) []fs.DirEntry{
    // create memory for sorted directory
    sorted_dir := make([]fs.DirEntry, len(dir))

    // define regex
    re := regexp.MustCompile(`.+_(\d+?)\..+?`)
    // put files
    for _, file := range dir{
        match_result := re.FindAllStringSubmatch(file.Name(), -1)
        file_number_str := match_result[0][1]
        file_number, _ := strconv.Atoi(file_number_str)
        sorted_dir[file_number] = file
    }

    // return sorted_dir
    return sorted_dir
}

func dir2pdf(dir_path, path_to_dir_pdf, path_to_dir_old string) string{

    // generate path to pdf
    pdf_path := filepath.Join(
        path_to_dir_pdf,filepath.Base(dir_path) + ".pdf")

    // read files and sort
    files, err := os.ReadDir(dir_path)
    files = sortdir(files)
    if err != nil {
        panic(err)
    }

    generate_pdf(dir_path, pdf_path, files)

    // move to old directory
    path_to_old := filepath.Join(path_to_dir_old, filepath.Base(dir_path))
    err = os.Rename(dir_path, path_to_old)
    if err != nil{
        panic(err)
    }

    // return pdf path string
    return pdf_path
}

func generate_pdf(dir_path string, pdf_path string, files []fs.DirEntry){
    // define regex
    re := regexp.MustCompile(`(?i)(.+\.(jpg|png))`)

    // init GoPdf
    pdf := gopdf.GoPdf{}

    // generate pdf file
    gopdf_started := false
    for _, file := range files{
        file_path := filepath.Join(dir_path, file.Name())
        if re.MatchString(file_path){
            img, err := os.Open(file_path)
            if err != nil{
                panic(err)
            }
            defer img.Close()
            img_conf, _, err := image.DecodeConfig(img)
            if err != nil {
                panic(err)
            }
            rect := gopdf.Rect{W: float64(img_conf.Width), H: float64(img_conf.Height)}
            if gopdf_started  != true{
                pdf.Start(gopdf.Config{PageSize: rect})
                gopdf_started = true
            }
            pageOpt := gopdf.PageOption{PageSize: &rect}
            pdf.AddPageWithOption(pageOpt)
            pdf.Image(file_path, 0, 0, &rect)
        }
    }
    pdf.WritePdf(pdf_path)
}

func zip2dir(zip_path, path_to_dir_old string) string{
    dir_name := strings.Replace(zip_path, ".zip", "", -1)
    unzip(zip_path, dir_name)
    path_to_old := filepath.Join(path_to_dir_old, filepath.Base(zip_path))
    err := os.Rename(zip_path, path_to_old)
    if err != nil{
        panic(err)
    }
    return dir_name
}

func unzip(src, dest string) error {
    // check exist extracted directory
    if _, err := os.Stat(dest); err != nil{
        // if not exist, create directory
        err := os.Mkdir(dest, 0777)
        if err != nil{
            panic(err)
        }
    }

    // unzip
    r, err := zip.OpenReader(src)
    defer r.Close()
    if err != nil {
        log.Print("zip.OpenReader error")
        panic(err)
    }
    for _, f := range r.File {
        rc, err  := f.Open()
        if err != nil {
            log.Print("f.Open error")
            panic(err)
        }
        defer rc.Close()
        path := filepath.Join(dest, f.Name)
        if f.FileInfo().IsDir(){
            os.MkdirAll(path, f.Mode())
        }else{
            f, err := os.OpenFile(
                path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
                if err != nil {
                    log.Print("os.OpenFile error")
                panic(err)
                }
            defer f.Close()

            _, err = io.Copy(f, rc)
            if err != nil {
                log.Print("io.Copy error")
                panic(err)
            }
        }
    }
    return nil
}
