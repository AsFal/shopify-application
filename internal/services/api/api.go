package api

import (
	"fmt"
	"log"
	"os"

	"github.com/AsFal/shopify-application/internal/pkg/classifier"
	"github.com/AsFal/shopify-application/internal/pkg/classifier/deepdetect"
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo/amazon"
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo/local"
	"github.com/AsFal/shopify-application/internal/pkg/search"
	essearch "github.com/AsFal/shopify-application/internal/pkg/search/elasticsearch"
	"github.com/AsFal/shopify-application/internal/pkg/tokenizer"
	estokenizer "github.com/AsFal/shopify-application/internal/pkg/tokenizer/elasticsearch"
)

type Service struct {
	imgrepo.ImgRepoClient
	classifier.Classifier
	search.SearchClient
	tokenizer.Tokenizer
}

func NewService() *Service {

	localImageFolder := os.Getenv("LOCAL_IMAGE_FOLDER")
	var imgRepoClient imgrepo.ImgRepoClient
	if localImageFolder != "" {
		// Use a local folder to store the images
		// Only used for local demo deployment or debugging purposes
		
		imgRepoClient = local.NewLocalImgRepo(localImageFolder, os.Getenv("DOCKER_HOST_VOLUME_PATH"))
	} else {
		session, err := amazon.ConnectAws(
			os.Getenv("ACCESS_KEY_ID"),
			os.Getenv("ACCES_KEY"),
			os.Getenv("REGION"), // TODO: Change this to a constant
		)
		if err != nil {
			log.Println("The AmazonS3 credentials provided are missing or incorrect.")
			log.Println("The API will not support Upload functionality.")
		}

		// TODO: The Connection function should verify that the bucket is valid
		imgRepoClient = amazon.NewAmazonS3Client(session, os.Getenv("BUCKET"))
	}

	deepdetectClassifier, err := deepdetect.NewDeepDetectClassifier(os.Getenv("DEEPDETECT_HOST"))
	if err != nil {
		log.Fatalf("Couldn't initialize the classifier service. Reason:\n %s", err.Error())
	}
	elasticsearchSearch := essearch.NewElasticsearchSearch(os.Getenv("ELASTICSEARCH_HOST"))
	elasticsearchTokenizer := estokenizer.NewTokenizer(os.Getenv("ELASTICSEARCH_HOST"))

	return &Service{
		ImgRepoClient: imgRepoClient,
		Classifier:    deepdetectClassifier,
		SearchClient:  elasticsearchSearch,
		Tokenizer:     elasticsearchTokenizer,
	}
}

func (s *Service) Run() {
	port := os.Getenv("PORT")
	log.Println(port)
	if port == "" {
		log.Fatal("Valid Port is required to run api")
	}
	log.Println("something")
	s.router().Run(fmt.Sprintf("0.0.0.0:%s", port))
}
