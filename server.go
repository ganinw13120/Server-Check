package main

import (
	"fmt"
	"os"
	"server-health/handler"
	"server-health/repository"
	"server-health/router"
	"server-health/service"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTimeZone() error {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	time.Local = location
	return nil
}

func fiberConfig() fiber.Config {
	return fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Books",
	}
}

func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}
}

func setupFiber() error {
	app := fiber.New(fiberConfig())
	app.Use(cors.New(corsConfig()))
	app.Use(recover.New())
	redis := setupRedis()

	bot, err := setupBot()
	if err != nil {
		return err
	}
	db, err := setupDatabase()
	if err != nil {
		return err
	}

	healthRepository := repository.NewHealthRepository(db)
	wishListRepository := repository.NewWishListRepository(db)

	healthService := service.NewHealthService(healthRepository, wishListRepository, bot)

	healthHandler := handler.NewHealthHandler(healthService, redis, bot)

	router.New(app, healthHandler)
	err = app.Listen(":" + os.Getenv("PORT"))

	return err
}

func setupBot() (*linebot.Client, error) {
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	return bot, err
}

func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func setupDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	err = setupTimeZone()
	if err != nil {
		panic(err)
	}
	err = setupFiber()
	if err != nil {
		panic(err)
	}
}
