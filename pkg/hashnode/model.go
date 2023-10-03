package hashnode

type CreateArticleRequest struct {
	Query     string `json:"query"`
	Variables struct {
		Input struct {
			Title               string `json:"title"`
			ContentMarkdown     string `json:"contentMarkdown"`
			Tags                []Tag  `json:"tags"`
			CoverImageURL       string `json:"coverImageURL"`
			IsPartOfPublication struct {
				PublicationId string `json:"publicationId"`
			} `json:"isPartOfPublication"`
		} `json:"input"`
	} `json:"variables"`
}
type Tag struct {
	Id   string `json:"_id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

const createPostRequestQuery = "mutation createStory($input: CreateStoryInput!){ createStory(input: $input){ code success message } }"

func NewCreateArticleRequest() CreateArticleRequest {
	return CreateArticleRequest{Query: createPostRequestQuery}
}

type GetPostRequest struct {
	Query     string `json:"query"`
	Variables struct {
		Slug     string `json:"slug"`
		Hostname string `json:"hostname"`
	} `json:"variables"`
}

const getPostRequestQuery = `
query Post($slug: String!, $hostname:String! ) {
    post(
        slug: $slug, 
        hostname: $hostname
    ) {
        slug
		title
    }
}
`

func NewGetPostRequest() GetPostRequest {
	return GetPostRequest{Query: getPostRequestQuery}
}

type Post struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}
