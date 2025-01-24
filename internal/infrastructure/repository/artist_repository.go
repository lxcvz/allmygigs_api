package repository

import (
	"allmygigs/internal/domain/entity"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type ArtistRepository struct {
	Client dynamodbiface.DynamoDBAPI
}

func NewArtistRepository(client dynamodbiface.DynamoDBAPI) *ArtistRepository {
	return &ArtistRepository{
		Client: client,
	}
}

func (r *ArtistRepository) GetArtistBySpotifyID(spotifyID string) (*entity.Artist, error) {
	key := map[string]*dynamodb.AttributeValue{
		"artist_id": {
			S: aws.String(spotifyID),
		},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("allmygigs"),
		Key:       key,
	}

	result, err := r.Client.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("error retrieving artist from database: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var artist entity.Artist
	err = dynamodbattribute.UnmarshalMap(result.Item, &artist)
	if err != nil {
		return nil, fmt.Errorf("error deserializing database response: %v", err)
	}

	return &artist, nil
}

func (r *ArtistRepository) GetArtistsBySpotifyIDsBatch(spotifyIDs []string) ([]*entity.Artist, []string, error) {
	var batchSize = 100
	var artists []*entity.Artist
	var missingIDs []string

	for i := 0; i < len(spotifyIDs); i += batchSize {
		end := i + batchSize
		if end > len(spotifyIDs) {
			end = len(spotifyIDs)
		}

		keys := make([]map[string]*dynamodb.AttributeValue, end-i)
		for j, id := range spotifyIDs[i:end] {
			keys[j] = map[string]*dynamodb.AttributeValue{
				"artist_id": {
					S: aws.String(id),
				},
			}
		}

		input := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				"allmygigs": {
					Keys: keys,
				},
			},
		}

		result, err := r.Client.BatchGetItem(input)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving artist from database: %v", err)
		}

		var foundArtists []*entity.Artist

		for _, item := range result.Responses["allmygigs"] {
			var artist entity.Artist
			err := dynamodbattribute.UnmarshalMap(item, &artist)
			if err != nil {
				return nil, nil, fmt.Errorf("error deserializing database response: %v", err)
			}
			foundArtists = append(foundArtists, &artist)
		}

		artists = append(artists, foundArtists...)

		for _, id := range spotifyIDs[i:end] {
			found := false
			for _, artist := range foundArtists {
				if artist.ArtistID == id {
					found = true
					break
				}
			}
			if !found {
				missingIDs = append(missingIDs, id)
			}
		}
	}

	fmt.Println(missingIDs)

	return artists, missingIDs, nil
}

func (r *ArtistRepository) SaveArtist(artist *entity.Artist) error {
	item, err := dynamodbattribute.MarshalMap(artist)
	if err != nil {
		return fmt.Errorf("error deserializing database artist: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("allmygigs"),
		Item:      item,
	}

	_, err = r.Client.PutItem(input)
	if err != nil {
		return fmt.Errorf("error saving artist to database: %v", err)
	}

	return nil
}

func (r *ArtistRepository) SaveArtistsBatch(artistsInfo []map[string]interface{}) error {
	var writeRequests []*dynamodb.WriteRequest

	for _, artist := range artistsInfo {
		artistItem := map[string]*dynamodb.AttributeValue{
			"artist_id": {
				S: aws.String(fmt.Sprintf("%v", artist["artist_id"])),
			},
			"artist_name": {
				S: aws.String(fmt.Sprintf("%v", artist["artist_name"])),
			},
			"artist_image": {
				S: aws.String(fmt.Sprintf("%v", artist["artist_image"])),
			},
		}

		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: artistItem,
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"allmygigs": writeRequests,
		},
	}

	_, err := r.Client.BatchWriteItem(input)
	if err != nil {
		return fmt.Errorf("error saving artist to database: %v", err)
	}

	return nil
}
