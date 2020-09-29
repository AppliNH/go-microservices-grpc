package greet_server

import (
	"applinh/gogrpcudemy/greet/greetpb"
	"context"
	"testing"
)

func TestGreet(t *testing.T) {

	s := server{}

	tests := []struct {
		req  greetpb.Greeting
		want string
	}{
		{
			req: greetpb.Greeting{
				FirstName: "Thomas",
				LastName:  "Martin",
			},
			want: "Hello Thomas Martin",
		},
	}

	for _, tt := range tests {
		req := &greetpb.GreetRequest{Greeting: &tt.req}
		resp, err := s.Greet(context.Background(), req)
		if err != nil {
			t.Errorf("GreetTest got unexpected error %v", err)
		}
		if resp.Result != tt.want {
			t.Errorf("HelloText(%v)=%v, wanted %v", tt.req, resp.Result, tt.want)
		}
	}

}
