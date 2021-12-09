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
	Category string   `json:"category"`
	Keywords []string `json:"keywords"`
}

// categorizeCmd represents the categorize command
var categorizeCmd = &cobra.Command{
	Use:   "categorize",
	Short: "categorizes given excel file according to given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		do()
	},
}

func init() {
	categorizeCmd.Flags().StringVarP(&excelFileName, "excel", "e", "", "excel file name")

	_ = cobra.MarkFlagRequired(categorizeCmd.Flags(), "excel")

	rootCmd.AddCommand(categorizeCmd)
}

func do() map[string]float64 {
	transactions, _ := readExcel()
	fmt.Println(transactions)

	categories, _ := loadCategories()
	keywords := createKeywordMap(categories)

	categoryDistribution := make(map[string]float64)

	for _, transaction := range transactions {
		found := false
		for keyword, category := range keywords {
			if strings.Contains(transaction.name, keyword) {
				categoryDistribution[category] = categoryDistribution[category] + transaction.price
				found = true
				break
			}
		}
		if found == false {
			//TODO: in this case as the user what to do
			fmt.Println()
			fmt.Printf("Couldn't find a category for: name: %s, price: %f, date: %s\n", transaction.name, transaction.price, transaction.date)
			fmt.Println("Which category do you want to place this expense?")

			fmt.Printf("%s:  ", "Your categories")
			for _, category := range categories.Categories {
				fmt.Printf("%s  ", category.Category)
			}

			fmt.Println()

			var catName string

			_, err := fmt.Scanf("%s \n", &catName)
			if err != nil {
				return nil
			}

			categoryDistribution[catName] = categoryDistribution[catName] + transaction.price

			fmt.Println("Which keyword in the expense name fits into the category?")

			var keyName string

			_, err = fmt.Scanf("%s \n", &keyName)
			if err != nil {
				return nil
			}

			index := indexOf(len(categories.Categories), func(i int) bool { return categories.Categories[i].Category == catName })

			if index > -1 {
				indexOfKey := indexOf(len(categories.Categories[index].Keywords), func(i int) bool { return categories.Categories[index].Keywords[i] == keyName })

				if indexOfKey == -1 {
					categories.Categories[index].Keywords = append(categories.Categories[index].Keywords, keyName)
				}
			} else {
				categories.Categories = append(categories.Categories, CategoryKeyValue{
					Category: catName,
					Keywords: []string{keyName},
				})
			}
		}
	}

	categoriesJson, _ := json.Marshal(categories)
	err := ioutil.WriteFile("config.json", categoriesJson, 0644)
	if err != nil {
		fmt.Errorf("%s", "error while writing into json")
		return nil
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
		if row[0] == "" {
			continue
		}
		name := strings.ToLower(row[0])
		priceWithDots := strings.Replace(row[2], ",", ".", -1)
		price, _ := strconv.ParseFloat(priceWithDots, 32)

		transactions = append(transactions, Transaction{
			name:  name,
			date:  row[1],
			price: math.Floor(float64(price)*100) / 100,
		})
	}

	return transactions, nil
}

func loadCategories() (Categories, error) {
	var cats Categories
	raw, err := ioutil.ReadFile(config)
	if err != nil {
		log.Println("Error occurred while reading config")
		return Categories{}, err
	}

	err = json.Unmarshal(raw, &cats)
	if err != nil {
		return Categories{}, err
	}

	return cats, nil
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

func indexOf(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
