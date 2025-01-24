package repository

import (
	"allmygigs/internal/domain/entity"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// ArtistRepository gerencia operações com a tabela allmygigs
type ArtistRepository struct {
	Client dynamodbiface.DynamoDBAPI
}

// NewArtistRepository cria uma nova instância de ArtistRepository
func NewArtistRepository(client dynamodbiface.DynamoDBAPI) *ArtistRepository {
	return &ArtistRepository{
		Client: client,
	}
}

// GetArtistBySpotifyID busca um artista na tabela pelo Spotify ID
func (r *ArtistRepository) GetArtistBySpotifyID(spotifyID string) (*entity.Artist, error) {
	// Define a chave para a busca
	key := map[string]*dynamodb.AttributeValue{
		"artist_id": {
			S: aws.String(spotifyID),
		},
	}

	// Faz a consulta utilizando o cliente DynamoDB
	input := &dynamodb.GetItemInput{
		TableName: aws.String("allmygigs"),
		Key:       key,
	}

	result, err := r.Client.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("error retrieving artist from database: %v", err)
	}

	if result.Item == nil {
		return nil, nil // Retorna nil se o artista não foi encontrado
	}

	var artist entity.Artist
	err = dynamodbattribute.UnmarshalMap(result.Item, &artist)
	if err != nil {
		return nil, fmt.Errorf("error deserializing database response: %v", err)
	}

	return &artist, nil
}

// GetArtistsBySpotifyIDsBatch busca múltiplos artistas na tabela DynamoDB utilizando BatchGetItem
func (r *ArtistRepository) GetArtistsBySpotifyIDsBatch(spotifyIDs []string) ([]*entity.Artist, []string, error) {
	// Limite do DynamoDB para BatchGetItem é 100 itens por vez, então precisamos dividir em lotes
	var batchSize = 100
	var artists []*entity.Artist
	var missingIDs []string

	for i := 0; i < len(spotifyIDs); i += batchSize {
		// Define o intervalo de IDs a serem processados no lote atual
		end := i + batchSize
		if end > len(spotifyIDs) {
			end = len(spotifyIDs)
		}

		// Cria um array de chaves para o BatchGetItem
		keys := make([]map[string]*dynamodb.AttributeValue, end-i)
		for j, id := range spotifyIDs[i:end] {
			keys[j] = map[string]*dynamodb.AttributeValue{
				"artist_id": {
					S: aws.String(id),
				},
			}
		}

		// Faz a consulta utilizando o cliente DynamoDB
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

		// Cria uma lista temporária para os artistas encontrados nesta iteração
		var foundArtists []*entity.Artist

		// Extrai os itens encontrados e os adiciona ao array de artistas encontrados
		for _, item := range result.Responses["allmygigs"] {
			var artist entity.Artist
			err := dynamodbattribute.UnmarshalMap(item, &artist)
			if err != nil {
				return nil, nil, fmt.Errorf("error deserializing database response: %v", err)
			}
			foundArtists = append(foundArtists, &artist)
		}

		// Adiciona os artistas encontrados ao array final
		artists = append(artists, foundArtists...)

		// Verifica quais IDs não foram encontrados no DynamoDB e adiciona a missingIDs
		for _, id := range spotifyIDs[i:end] {
			found := false
			// Verifica se o ID foi encontrado nos artistas retornados
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

// SaveArtist salva um novo artista na tabela DynamoDB
func (r *ArtistRepository) SaveArtist(artist *entity.Artist) error {
	// Serializa o objeto Artist
	item, err := dynamodbattribute.MarshalMap(artist)
	if err != nil {
		return fmt.Errorf("error deserializing database artist: %v", err)
	}

	// Faz a inserção do item na tabela DynamoDB
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

// SaveArtistsBatch insere múltiplos artistas no DynamoDB em um único batch
func (r *ArtistRepository) SaveArtistsBatch(artistsInfo []map[string]interface{}) error {
	// Prepara os itens para a inserção no DynamoDB
	var writeRequests []*dynamodb.WriteRequest

	// Prepara as requisições de inserção
	for _, artist := range artistsInfo {
		// Converte o artista para o formato esperado pelo DynamoDB
		artistItem := map[string]*dynamodb.AttributeValue{
			"artist_id": {
				S: aws.String(fmt.Sprintf("%v", artist["artist_id"])), // Garantir que seja uma string
			},
			"artist_name": {
				S: aws.String(fmt.Sprintf("%v", artist["artist_name"])), // Garantir que seja uma string
			},
			"artist_image": {
				S: aws.String(fmt.Sprintf("%v", artist["artist_image"])), // Garantir que seja uma string
			},
		}

		// Cria uma requisição de inserção para o DynamoDB
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: artistItem,
			},
		})
	}

	// Faz a inserção no DynamoDB
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"allmygigs": writeRequests,
		},
	}

	// Chama o BatchWriteItem
	_, err := r.Client.BatchWriteItem(input)
	if err != nil {
		return fmt.Errorf("error saving artist to database: %v", err)
	}

	return nil
}
