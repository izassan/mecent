package main

import (
	"archive/zip"
	"log"
	"os"
    "io"
    "unicode/utf8"
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
    sorted_dir := make([]fs.DirEntry, len(dir))
    for _, file := range dir{
        pagenum := 0
        rate := 1
        started := false
        b := []byte(file.Name())
        for len(b) > 0{
            r, size := utf8.DecodeLastRune(b)
            char :=  string(r)
            if char == "."{
                started = true
            }else if started && "0" <= char && char <= "9"{
                int_char, _ := strconv.Atoi(char)
                pagenum = pagenum + int_char * rate
                rate = rate * 10
            }else if char == "_"{
                break
            }
            b = b[:len(b)-size]
        }
        sorted_dir[pagenum] = file
    }
    return sorted_dir
}

func dir2pdf(dir_path string) string{
    re := regexp.MustCompile(`(?i)(.+\.(jpg|png))`)
    pdf := gopdf.GoPdf{}
    pdf_path := filepath.Join(
        filepath.Dir(dir_path),filepath.Base(dir_path) + ".pdf")
    files, err := os.ReadDir(dir_path)
    files = sortdir(files)
    if err != nil {
        panic(err)
    }
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
    log.Print("convert dir to pdf: ", dir_path)
    return pdf_path
}

func zip2dir(zip_path string) string{
    log.Print("convert zip to dir: ", zip_path)
    dir_name := strings.Replace(zip_path, ".zip", "", -1)
    os.Mkdir(dir_name, 0766)
    unzip(zip_path, dir_name)
    return dir_name
}

func unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
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
        log.Print("path: ", path)
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
