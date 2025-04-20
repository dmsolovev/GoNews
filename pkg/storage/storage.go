package storage

// // Post - публикация.
// type Post struct {
// 	ID          int
// 	Title       string
// 	Content     string
// 	AuthorID    int
// 	AuthorName  string
// 	CreatedAt   int64
// 	PublishedAt int64
// }
type Post struct {
	ID          int    `bson:"id" json:"id"`
	AuthorID    int    `bson:"authorid" json:"authorid"`
	AuthorName  string `bson:"authorname" json:"authorname"`
	Title       string `bson:"title" json:"title"`
	Content     string `bson:"content" json:"content"`
	CreatedAt   int64  `bson:"createdat" json:"createdat"`
	PublishedAt int64  `bson:"publishedat" json:"publishedat"`
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
}
