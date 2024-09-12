package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type Package struct {
	Name      string `json:"name"`
	Epoch     int    `json:"epoch"`
	Version   string `json:"version"`
	Release   string `json:"release"`
	Arch      string `json:"arch"`
	Disttag   string `json:"disttag"`
	Buildtime int    `json:"buildtime"`
	Source    string `json:"source"`
}

type Data struct {
	RequestArgs map[string]interface{} `json:"request_args"`
	Length      int                    `json:"length"`
	Packages    []Package              `json:"packages"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "branch-binary-packages",
	}

	var getPackagesCmd = &cobra.Command{
		Use: "get-packages",
		Run: func(cmd *cobra.Command, args []string) {
			branches := []string{"p10", "sisyphus"}
			for _, branch := range branches {
				apiURL := fmt.Sprintf("https://rdb.altlinux.org/api/export/branch_binary_packages/%s", branch)
				resp, err := http.Get(apiURL)
				if err != nil {
					fmt.Println("Ошибка при выполнении запроса:", err)
					os.Exit(1)
				}
				defer resp.Body.Close()

				filename := fmt.Sprintf("%s.json", branch)

				file, err := os.Create(filename)
				if err != nil {
					fmt.Println("Ошибка при создании файла:", err)
					return
				}
				defer file.Close()

				_, err = io.Copy(file, resp.Body)
				if err != nil {
					fmt.Println("Ошибка при копировании данных:", err)
					return
				}

				fmt.Printf("Данные успешно записаны в файл %s.json", branch)
			}
		},
	}

	rootCmd.AddCommand(getPackagesCmd)

	var comparePackagesCmd = &cobra.Command{
		Use: "compare-packages",
		Run: func(cmd *cobra.Command, args []string) {
			data1, err := os.ReadFile("p10.json")
			if err != nil {
				fmt.Println("Ошибка при чтении файла:", err)
			}

			data2, err := os.ReadFile("sisyphus.json")
			if err != nil {
				fmt.Println("Ошибка при чтении файла:", err)
			}

			var dataStruct1, dataStruct2 Data

			err = json.Unmarshal(data1, &dataStruct1)
			if err != nil {
				fmt.Println("Ошибка при разборе JSON:", err)
			}
			packages1 := dataStruct1.Packages

			var package1Names []string
			for _, pkg := range packages1 {
				package1Names = append(package1Names, pkg.Name)
			}

			fmt.Println(package1Names[:5])

			err = json.Unmarshal(data2, &dataStruct2)
			if err != nil {
				fmt.Println("Ошибка при разборе JSON:", err)
			}
			packages2 := dataStruct1.Packages

			var package2Names []string
			for _, pkg := range packages2 {
				package2Names = append(package2Names, pkg.Name)
			}

			fmt.Println(package2Names[:5])
		},
	}

	rootCmd.AddCommand(comparePackagesCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
