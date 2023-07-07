package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func ConnectRedisDB() error {
	// Создание клиента Redis
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес и порт Redis сервера
		Password: "",               // Пароль (если требуется)
		DB:       0,                // Индекс базы данных
	})
	// Проверка соединения с Redis
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	log.Println("Connected to Redis:", pong)

	return nil
}

func GetRedis(key string) (string, error) {
	// Получение значения из Redis по ключу
	result, err := client.Get(context.Background(), key).Result()
	return result, err
}

func SetRedis(key string, value string, expiration time.Duration) error {
	// Установка значения в Redis с указанием срока жизни
	err := client.Set(context.Background(), key, value, expiration).Err()
	return err
}

func SelectRedis(index int) error {
	// Установка таблицы в Redis
	err := client.Do(context.Background(), "SELECT", index).Err()
	return err
}

func CloseRedisDB() error {
	if client == nil {
		return nil // Проверяем, что клиент Redis был инициализирован
	}
	// Закрытие соединения с Redis
	err := client.Close()
	if err != nil {
		return err
	}
	log.Println("Redis connection closed")
	return nil
}
