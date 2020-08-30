package mdlog

import (
	"fmt"
	"log"
)

type Plan map[*Post]*Post

func (p Plan) Print() {
	if len(p) == 0 {
		fmt.Println("local posts is up to date with remote posts")
		return
	}

	log.Println("------- Change List -------")
	for remotePost, localPost := range p {
		log.Printf("<< Post [%s] >>", remotePost.ID)
		if remotePost.Title != localPost.Title {
			log.Printf("Title: %s -> %s\n", remotePost.Title, localPost.Title)
		}
		if remotePost.Content != localPost.Content {
			log.Println("Content: (diff is not supported yet)")
		}
		if len(remotePost.Tags) != len(localPost.Tags) {
			log.Printf("Tags: %s -> %s\n", remotePost.Tags, localPost.Tags)
		} else {
			var neq bool
			for i, v := range remotePost.Tags {
				if v != localPost.Tags[i] {
					neq = true
					break
				}
			}
			if neq {
				log.Printf("Tags: %s -> %s\n", remotePost.Tags, localPost.Tags)
			}
		}
		if remotePost.Status != localPost.Status {
			log.Printf("Status: %s -> %s\n", PostStatusText(remotePost.Status), PostStatusText(localPost.Status))
		}
		log.Println("---------------------------")
	}
}
