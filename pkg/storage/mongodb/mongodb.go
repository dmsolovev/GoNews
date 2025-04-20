package mongodb

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

func New(connstr string) (*Store, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}

	db := client.Database("gonews")
	return &Store{
		client: client,
		db:     db,
	}, nil
}

// func (s *Store) Posts() ([]storage.Post, error) {
// 	collection := s.db.Collection("posts")
// 	fmt.Println(collection)
// 	cursor, err := collection.Find(context.Background(), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

//		var posts []storage.Post
//		if err = cursor.All(context.Background(), &posts); err != nil {
//			return nil, err
//		}
//		return posts, nil
//	}
func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Collection("posts")

	// Логируем количество документов
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Total documents in collection: %d\n", count)

	// Получаем все документы
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Читаем в структуру Post
	var posts []storage.Post
	for cursor.Next(context.Background()) {
		var post storage.Post
		if err := cursor.Decode(&post); err != nil {
			fmt.Printf("Error decoding document: %v\n", err)
			continue
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Логируем результат
	fmt.Printf("Decoded posts: %+v\n", posts)

	return posts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.db.Collection("posts")
	_, err := collection.InsertOne(context.Background(), p)
	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	collection := s.db.Collection("posts")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"id": p.ID},
		bson.M{"$set": p},
	)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	collection := s.db.Collection("posts")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": p.ID})
	return err
}
