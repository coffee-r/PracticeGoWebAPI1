package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

	r.GET("/users", func(c *gin.Context) {
		query := "SELECT * FROM USERS"
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	r.POST("/users", func(c *gin.Context) {
		var user User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// データベースにユーザーを挿入するクエリ
		query := `
			INSERT INTO USERS (name, email) 
			VALUES (?, ?);
			SELECT SCOPE_IDENTITY();
		`
		row := db.QueryRow(query, user.Name, user.Email)

		var id int
		if err := row.Scan(&id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting inserted ID"})
			return
		}

		user.ID = id
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		// URLパラメータからユーザーIDを取得
		idStr := c.Param("id")

		// IDを整数に変換
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user User

		user.ID = id

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		query := "UPDATE USERS SET name = ?, email = ? where id = ?"
		result, err := db.Exec(query, user.Name, user.Email, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update user"})
			return
		}

		// 更新された行数を確認
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting affected rows"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// 削除成功のレスポンス
		c.JSON(http.StatusOK, gin.H{"message": "User put successfully"})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		// URLパラメータからユーザーIDを取得
		idStr := c.Param("id")

		// IDを整数に変換
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// データベースからユーザーを削除するクエリ
		query := "DELETE FROM USERS WHERE id = ?"
		result, err := db.Exec(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
			return
		}

		// 削除された行数を確認
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting affected rows"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// 削除成功のレスポンス
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

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
