package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Table struct {
	Columns []string
	Rows    []map[string]string
}

type Database struct {
	Tables map[string]Table
}

var (
	currentDB   string
	databases   = make(map[string]Database)
	storageFile = "db.json"
)

func saveDB() {
	data, _ := json.MarshalIndent(databases, "", "  ")
	_ = os.WriteFile(storageFile, data, 0644)
}

func loadDB() {
	data, err := os.ReadFile(storageFile)
	if err == nil {
		_ = json.Unmarshal(data, &databases)
	}
}

// --- Console Table Renderer ---
func printTable(columns []string, rows []map[string]string) {
	// column widths
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = len(col)
	}
	for _, row := range rows {
		for i, col := range columns {
			if len(row[col]) > widths[i] {
				widths[i] = len(row[col])
			}
		}
	}

	// border line
	border := "+"
	for _, w := range widths {
		border += strings.Repeat("-", w+2) + "+"
	}

	// header
	fmt.Println(border)
	fmt.Print("|")
	for i, col := range columns {
		fmt.Printf(" %-*s |", widths[i], col)
	}
	fmt.Println()
	fmt.Println(border)

	// rows
	for _, row := range rows {
		fmt.Print("|")
		for i, col := range columns {
			fmt.Printf(" %-*s |", widths[i], row[col])
		}
		fmt.Println()
	}
	fmt.Println(border)
}

func main() {
	loadDB()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔹 Mini SQL-like DB (with USE, SHOW, SELECT, DELETE, pretty tables)")
	for {
		fmt.Print("SQL> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		// CREATE DATABASE dbname;
		case "CREATE":
			if len(parts) >= 3 && strings.ToUpper(parts[1]) == "DATABASE" {
				dbName := parts[2]
				if _, exists := databases[dbName]; exists {
					fmt.Println("Database already exists.")
				} else {
					databases[dbName] = Database{Tables: make(map[string]Table)}
					fmt.Println("Database created:", dbName)
				}
				saveDB()
			} else if len(parts) >= 3 && strings.ToUpper(parts[1]) == "TABLE" {
				if currentDB == "" {
					fmt.Println("No database selected. Use `USE dbname;` first.")
					break
				}
				tableName := parts[2]
				cols := strings.Split(strings.Join(parts[3:], " "), ",")
				for i := range cols {
					cols[i] = strings.TrimSpace(cols[i])
				}
				db := databases[currentDB]
				db.Tables[tableName] = Table{Columns: cols, Rows: []map[string]string{}}
				databases[currentDB] = db
				fmt.Println("Table created:", tableName)
				saveDB()
			}

		// USE dbname;
		case "USE":
			if len(parts) >= 2 {
				dbName := parts[1]
				if _, ok := databases[dbName]; ok {
					currentDB = dbName
					fmt.Println("Using database:", dbName)
				} else {
					fmt.Println("Database does not exist.")
				}
			}

		// SHOW DATABASES;
		case "SHOW":
			if len(parts) >= 2 && strings.ToUpper(parts[1]) == "DATABASES" {
				fmt.Println("Databases:")
				for dbName := range databases {
					fmt.Println(" -", dbName)
				}
			} else if len(parts) >= 2 && strings.ToUpper(parts[1]) == "TABLES" {
				if currentDB == "" {
					fmt.Println("No database selected.")
					break
				}
				fmt.Println("Tables in", currentDB+":")
				for tName := range databases[currentDB].Tables {
					fmt.Println(" -", tName)
				}
			}

		// INSERT INTO table (col=val, col=val ...);
		case "INSERT":
			if currentDB == "" {
				fmt.Println("No database selected. Use `USE dbname;` first.")
				break
			}
			if len(parts) < 4 {
				fmt.Println("Usage: INSERT INTO table col=val col=val ...")
				break
			}
			tableName := parts[2]
			assignments := parts[3:]
			row := make(map[string]string)
			for _, assign := range assignments {
				kv := strings.SplitN(assign, "=", 2)
				if len(kv) == 2 {
					row[kv[0]] = kv[1]
				}
			}
			db := databases[currentDB]
			t := db.Tables[tableName]
			t.Rows = append(t.Rows, row)
			db.Tables[tableName] = t
			databases[currentDB] = db
			fmt.Println("Inserted into", tableName)
			saveDB()

		// SELECT * FROM table;
		case "SELECT":
			if currentDB == "" {
				fmt.Println("No database selected. Use `USE dbname;` first.")
				break
			}
			if len(parts) >= 4 && parts[1] == "*" && strings.ToUpper(parts[2]) == "FROM" {
				tableName := parts[3]
				if table, ok := databases[currentDB].Tables[tableName]; ok {
					if len(table.Rows) == 0 {
						fmt.Println("(empty)")
					} else {
						printTable(table.Columns, table.Rows)
					}
				} else {
					fmt.Println("Table does not exist:", tableName)
				}
			}

		// DELETE FROM table WHERE col=val;
		case "DELETE":
			if currentDB == "" {
				fmt.Println("No database selected. Use `USE dbname;` first.")
				break
			}
			if len(parts) < 5 || strings.ToUpper(parts[1]) != "FROM" || strings.ToUpper(parts[3]) != "WHERE" {
				fmt.Println("Usage: DELETE FROM table WHERE col=val")
				break
			}
			tableName := parts[2]
			condition := parts[4]
			kv := strings.SplitN(condition, "=", 2)
			if len(kv) != 2 {
				fmt.Println("Condition must be in format col=val")
				break
			}
			col, val := kv[0], kv[1]

			db := databases[currentDB]
			t, ok := db.Tables[tableName]
			if !ok {
				fmt.Println("Table does not exist:", tableName)
				break
			}

			// Filter rows
			newRows := []map[string]string{}
			deleted := 0
			for _, row := range t.Rows {
				if row[col] == val {
					deleted++
					continue
				}
				newRows = append(newRows, row)
			}
			t.Rows = newRows
			db.Tables[tableName] = t
			databases[currentDB] = db

			fmt.Printf("Deleted %d row(s) from %s\n", deleted, tableName)
			saveDB()

		case "EXIT", "QUIT":
			fmt.Println("Bye 👋")
			return
        // DROP TABLE tablename;
case "DROP":
    if currentDB == "" {
        fmt.Println("No database selected. Use `USE dbname;` first.")
        break
    }
    if len(parts) < 3 || strings.ToUpper(parts[1]) != "TABLE" {
        fmt.Println("Usage: DROP TABLE tablename")
        break
    }
    tableName := parts[2]
    db := databases[currentDB]
    if _, ok := db.Tables[tableName]; ok {
        delete(db.Tables, tableName)
        databases[currentDB] = db
        fmt.Println("Dropped table:", tableName)
        saveDB()
    } else {
        fmt.Println("Table does not exist:", tableName)
    }

		default:
			fmt.Println("Unknown command:", input)
		}
	}
}
