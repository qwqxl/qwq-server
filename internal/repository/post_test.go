package repository

import (
	"context"
	"fmt"
	"log"
	"qwqserver/internal/model"
	"qwqserver/pkg/database"
	"testing"
)

func TestPostRepoRun(t *testing.T) {
	t.Run("test-post-repo", func(t *testing.T) {
		// 初始化数据库连接
		cfg := &database.Config{
			Driver:   "mysql",
			DSN:      "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			LogLevel: "info",
		}

		if _, err := database.InitDB(cfg); err != nil {
			log.Fatalf("数据库初始化失败: %v", err)
		}
		defer database.Close()

		// 创建Post仓库
		postRepo, err := NewPostRepository()
		if err != nil {
			log.Fatal(err)
		}

		// 创建新帖子
		newPost := &model.Post{
			Title:     "Go语言最佳实践",
			Content:   "本文介绍Go语言开发中的最佳实践...",
			AuthorUID: 1,
			//CategoryID: 2,
		}

		if err := postRepo.Create(context.Background(), newPost); err != nil {
			log.Fatalf("创建帖子失败: %v", err)
		}
		fmt.Printf("创建帖子成功，ID: %d\n", newPost.ID)

		// 获取帖子
		post, err := postRepo.FindByID(context.Background(), newPost.ID)
		if err != nil {
			log.Fatal(err)
		}
		if post != nil {
			fmt.Printf("获取帖子: %s (作者ID: %d)\n", post.Title, post.AuthorUID)
		}

		// 增加浏览量
		if err := postRepo.IncrementViewCount(context.Background(), post.ID); err != nil {
			log.Printf("增加浏览量失败: %v", err)
		}

		// 增加点赞数
		if err := postRepo.IncrementLikeCount(context.Background(), post.ID); err != nil {
			log.Printf("增加点赞数失败: %v", err)
		}

		// 获取用户的所有帖子
		userPosts, total, err := postRepo.ListByUserID(context.Background(), 1, 1, 10)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("用户共有 %d 篇帖子，第一页: \n", total)
		for i, p := range userPosts {
			fmt.Printf("%d. %s\n", i+1, p.Title)
		}

		// 获取热门帖子
		popularPosts, err := postRepo.ListPopular(context.Background(), 7, 5)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("本周热门帖子:")
		for i, p := range popularPosts {
			//fmt.Printf("%d. %s (%d 赞)\n", i+1, p.Title, p.LikeCount)
			fmt.Printf("%d. %s (%d 赞)\n", i+1, p.Title, p) // X
		}

		// 置顶帖子
		if err := postRepo.PinPost(context.Background(), newPost.ID); err != nil {
			log.Printf("置顶帖子失败: %v", err)
		}

		// 在事务中更新帖子
		err = postRepo.WithTransaction(context.Background(), func(txRepo PostRepository) error {
			// 更新帖子内容
			post.Content = "更新后的内容..."
			if err := txRepo.Update(context.Background(), post); err != nil {
				return err
			}

			// 增加评论数
			if err := txRepo.IncrementCommentCount(context.Background(), post.ID); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Fatalf("事务执行失败: %v", err)
		}

		fmt.Println("帖子更新和评论数增加事务成功完成")
	})
}
