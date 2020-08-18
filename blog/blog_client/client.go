package blog_client

import (
	"applinh/gogrpcudemy/blog/blogpb"
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
)

func StartClient() {

	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50053", opts)
	if err != nil {
		log.Fatalf("Failed to connect %v \n", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	// Create blog
	blog := &blogpb.Blog{
		AuthorId: "AppliNH",
		Title:    "1st Blog",
		Content:  "Nothing much.",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v \n", err)
	}
	fmt.Printf("Done ! %v \n", res)
	blogId := res.GetBlog().GetId()
	fmt.Println(strings.Repeat("_", 25))
	// Read blog

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5bdc29e661b75adcac496cf4"})
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}
	fmt.Println(strings.Repeat("_", 25))
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	resRead, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		log.Fatalf("Error: %v \n", err)
	}

	fmt.Printf("Blog: %v \n", resRead.GetBlog())

}
