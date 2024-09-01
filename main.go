package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

type User struct {
	id    int
	name  string
	email string
}

func main() {
	// 環境変数から接続情報を取得
	dbServer := os.Getenv("DB_SERVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 接続文字列を構築
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", dbServer, dbUser, dbPassword, dbName)

	// データベースに接続
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()

	// 接続確認
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %s", err)
	}
	fmt.Println("Connected to SQL Server!")

	// Ginのルーターを作成
	r := gin.Default()

	// エンドポイントの設定
	r.GET("/databases", func(c *gin.Context) {
		// データ選択のクエリを実行
		query := "SELECT name FROM sys.databases"
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		// 結果を格納するスライス
		var databases []string

		// 結果を表示
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
				return
			}
			databases = append(databases, name)
		}

		// エラーが発生していないか確認
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating rows"})
			return
		}

		// JSON形式で返す
		c.JSON(http.StatusOK, gin.H{"databases": databases})
	})

	// サーバーの起動
	r.Run(":8080")
}
