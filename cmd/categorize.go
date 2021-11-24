package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var excelFileName string

type transaction struct {
	name  string
	date  string
	price float64
}

type categories struct {
	Categories []categoryKeyValue `json:"categories"`
}

type categoryKeyValue struct {
	Category string `json:"category"`
	Keywords []string `json:"keywords"`
}

// categorizeCmd represents the categorize command
var categorizeCmd = &cobra.Command{
	Use:   "categorize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		loadCategories()
		readExcel()
	},
}

func init() {
	categorizeCmd.Flags().StringVarP(&excelFileName, "excel", "e", "", "excel file name")

	_ = cobra.MarkFlagRequired(categorizeCmd.Flags(), "excel")

	rootCmd.AddCommand(categorizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// categorizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// categorizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readExcel() ([]transaction, error) {
	f, err := excelize.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var transactions []transaction

	rows, err := f.GetRows("sheet")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, row := range rows[2:] {
		priceWithDots := strings.Replace(row[2], ",", ".", -1)
		price, _ := strconv.ParseFloat(priceWithDots, 32)

		transactions = append(transactions, transaction{
			name:  row[0],
			date:  row[1],
			price: math.Floor(float64(price)*100) / 100,
		})
	}

	return transactions, nil
}

func loadCategories() (categories, error){
	var Categories categories
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Error occured while reading config")
		return categories{}, err
	}
	json.Unmarshal(raw, &Categories)

	return Categories, nil
}

func categorize(transactions []transaction, conf categories) {

}