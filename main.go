package main

import(
    "log"
    "os"
    "io"
    "regexp"
    "path/filepath"
    "github.com/urfave/cli/v2"
)

func move_media(src_path string, dest_path string){
    re := regexp.MustCompile(`(?i)(.+\.(jpg|png|gif|mp4))`)
    files, err := os.ReadDir(src_path)
    if err != nil{
        panic(err)
    }
    os.Mkdir(dest_path, 0766)
    var transfer_file_path string
    for _, file := range files{
        file_path := filepath.Join(src_path, file.Name())
        if file.IsDir(){
            transfer_file_path = dir2pdf(file_path)
        }else if filepath.Ext(file.Name()) == ".zip"{
            dir_path := zip2dir(file_path)
            transfer_file_path = dir2pdf(dir_path)
        }else if re.MatchString(file_path){
            log.Print("match file: ", file_path)
            transfer_file_path = file_path
        }else{
            continue
        }
        log.Print(transfer_file_path)
        transfer_file, err := os.Open(transfer_file_path)
        dest_file_path := filepath.Join(dest_path, filepath.Base(transfer_file_path))
        if err != nil{
            panic(err)
        }
        defer transfer_file.Close()

        dst, err := os.Create(dest_file_path)
        if err != nil{
            panic(err)
        }
        defer dst.Close()

        _, err = io.Copy(dst, transfer_file)
        if err != nil{
            panic(err)
        }
    }
}

func main(){
    var dest_path string
    var src_path string
    app := &cli.App{
        Name: "mecent",
        Usage: "multi jpeg files convert to pdf. and media files send other directory",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:   "dest_path",
                Aliases: []string{"d"},
                Value:  "./dest",
                Usage:   "set destinaion path",
                Destination: &dest_path,
            },
            &cli.StringFlag{
                Name:   "src_path",
                Aliases: []string{"s"},
                Value:  "./src",
                Usage:   "set source path",
                Destination: &src_path,
            },
            &cli.BoolFlag{
                Name: "remove",
            },
        },
        Action: func(c *cli.Context) error {
            move_media(src_path, dest_path)
            if c.Bool("remove") {
                log.Print("remove directory...: ", src_path)
                err := os.RemoveAll(src_path)
                if err != nil{
                    panic(err)
                }
                os.Mkdir(src_path, 766)
            }
            return nil
        },
    }

    err := app.Run(os.Args)
    if err != nil{
        log.Fatal(err)
    }
}
