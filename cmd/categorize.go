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
const config = "config.json"

var excelFileName string

type Transaction struct {
	name  string
	date  string
	price float64
}

type Categories struct {
	Categories []CategoryKeyValue `json:"categories"`
}

type CategoryKeyValue struct {
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
		do()
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

func do() map[string]float64 {
	transactions, _ := readExcel()
	keywords, _ := loadCategories()

	categoryDistribution := make(map[string]float64)

	for _, transaction := range transactions{
		found := false
		for keyword, category := range keywords{
			if strings.Contains(transaction.name, keyword){
				categoryDistribution[category] = categoryDistribution[category] + transaction.price
				found = true
				break
			}
		}
		if found == false{
			//TODO: in this case as the user what to do
			categoryDistribution["other"] = categoryDistribution["other"] + transaction.price
		}
	}

	fmt.Println(categoryDistribution)
	return categoryDistribution
}

func readExcel() ([]Transaction, error) {
	f, err := excelize.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var transactions []Transaction

	rows, err := f.GetRows("sheet")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, row := range rows[2:] {
		priceWithDots := strings.Replace(row[2], ",", ".", -1)
		price, _ := strconv.ParseFloat(priceWithDots, 32)
		name := strings.ToLower(row[0])

		transactions = append(transactions, Transaction{
			name:  name,
			date:  row[1],
			price: math.Floor(float64(price)*100) / 100,
		})
	}

	return transactions, nil
}

func loadCategories() (map[string]string, error){
	var cats Categories
	raw, err := ioutil.ReadFile(config)
	if err != nil {
		log.Println("Error occurred while reading config")
		return nil, err
	}

	err = json.Unmarshal(raw, &cats)
	if err != nil {
		return nil, err
	}

	return createKeywordMap(cats), nil
}

func createKeywordMap(cats Categories) map[string]string {
	keywordMap := make(map[string]string)

	for _, value := range cats.Categories {
		category := value.Category
		for _, keyword := range value.Keywords {
			keywordMap[keyword] = category
		}
	}

	return keywordMap
}