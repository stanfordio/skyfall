package schema

import (
	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/bigquery"
)

func GetSchema() bigquery.Schema {
	// Generated automatically from exported files using https://pypi.org/project/bigquery-schema-generator/
	// E.g., generate-schema < file.data.json
	// Note for updates you need exact raw schema from bq show --schema
	rawSchema := `[
	{
		"name": "Action",
		"type": "STRING"
	},
	{
		"name": "CreatedAt",
		"type": "TIMESTAMP"
	},
	{
		"name": "Full",
		"type": "STRING"
	},
	{
		"fields": [
		{
			"fields": [
			{
				"name": "DID",
				"type": "STRING"
			},
			{
				"name": "DIDKey",
				"type": "STRING"
			},
			{
				"name": "Handle",
				"type": "STRING"
			},
			{
				"name": "PDS",
				"type": "STRING"
			}
			],
			"name": "Actor",
			"type": "RECORD"
		},
		{
			"fields": [
			{
				"name": "Avatar",
				"type": "STRING"
			},
			{
				"name": "Description",
				"type": "STRING"
			},
			{
				"name": "DID",
				"type": "STRING"
			},
			{
				"name": "DisplayName",
				"type": "STRING"
			},
			{
				"name": "FollowersCount",
				"type": "INTEGER"
			},
			{
				"name": "FollowsCount",
				"type": "INTEGER"
			},
			{
				"name": "Handle",
				"type": "STRING"
			},
			{
				"name": "PostsCount",
				"type": "INTEGER"
			},
			{
				"name": "IndexedAt",
				"type": "TIMESTAMP"
			}
			],
			"name": "BlockedProfile",
			"type": "RECORD"
		},
		{
			"fields": [
			{
				"name": "Avatar",
				"type": "STRING"
			},
			{
				"name": "Description",
				"type": "STRING"
			},
			{
				"name": "DID",
				"type": "STRING"
			},
			{
				"name": "DisplayName",
				"type": "STRING"
			},
			{
				"name": "FollowersCount",
				"type": "INTEGER"
			},
			{
				"name": "FollowsCount",
				"type": "INTEGER"
			},
			{
				"name": "Handle",
				"type": "STRING"
			},
			{
				"name": "PostsCount",
				"type": "INTEGER"
			},
			{
				"name": "IndexedAt",
				"type": "TIMESTAMP"
			}
			],
			"name": "FollowedProfile",
			"type": "RECORD"
		},
		{
			"fields": [
			{
				"fields": [
				{
					"name": "Avatar",
					"type": "STRING"
				},
				{
					"name": "DID",
					"type": "STRING"
				},
				{
					"name": "DisplayName",
					"type": "STRING"
				},
				{
					"name": "Handle",
					"type": "STRING"
				},
				{
					"name": "IndexedAt",
					"type": "TIMESTAMP"
				}
				],
				"name": "Author",
				"type": "RECORD"
			},
			{
				"name": "CID",
				"type": "STRING"
			},
			{
				"name": "CreatedAt",
				"type": "TIMESTAMP"
			},
			{
				"fields": [
				{
					"fields": [
					{
						"name": "Alt",
						"type": "STRING"
					},
					{
						"name": "BlobLink",
						"type": "STRING"
					},
					{
						"name": "Height",
						"type": "INTEGER"
					},
					{
						"name": "MimeType",
						"type": "STRING"
					},
					{
						"name": "Width",
						"type": "INTEGER"
					}
					],
					"mode": "REPEATED",
					"name": "EmbedRecordMedia",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "Description",
						"type": "STRING"
					},
					{
						"name": "Title",
						"type": "STRING"
					},
					{
						"name": "URI",
						"type": "STRING"
					}
					],
					"name": "External",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "Alt",
						"type": "STRING"
					},
					{
						"name": "BlobLink",
						"type": "STRING"
					},
					{
						"name": "Height",
						"type": "INTEGER"
					},
					{
						"name": "MimeType",
						"type": "STRING"
					},
					{
						"name": "Width",
						"type": "INTEGER"
					}
					],
					"mode": "REPEATED",
					"name": "Images",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "CID",
						"type": "STRING"
					},
					{
						"name": "Type",
						"type": "STRING"
					},
					{
						"name": "URI",
						"type": "STRING"
					}
					],
					"name": "Record",
					"type": "RECORD"
				}
				],
				"name": "Embed",
				"type": "RECORD"
			},
			{
				"mode": "REPEATED",
				"name": "Langs",
				"type": "STRING"
			},
			{
				"name": "LikeCount",
				"type": "INTEGER"
			},
			{
				"name": "RepostCount",
				"type": "INTEGER"
			},
			{
				"mode": "REPEATED",
				"name": "Hashtags",
				"type": "STRING"
			},
			{
				"mode": "REPEATED",
				"name": "URLs",
				"type": "STRING"
			},
			{
				"name": "Text",
				"type": "STRING"
			},
			{
				"name": "URI",
				"type": "STRING"
			}
			],
			"name": "LikedPost",
			"type": "RECORD"
		},
		{
			"fields": [
			{
				"name": "CreatedAt",
				"type": "TIMESTAMP"
			},
			{
				"fields": [
				{
					"fields": [
					{
						"name": "Alt",
						"type": "STRING"
					},
					{
						"name": "BlobLink",
						"type": "STRING"
					},
					{
						"name": "Height",
						"type": "INTEGER"
					},
					{
						"name": "MimeType",
						"type": "STRING"
					},
					{
						"name": "Width",
						"type": "INTEGER"
					}
					],
					"mode": "REPEATED",
					"name": "EmbedRecordMedia",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "Description",
						"type": "STRING"
					},
					{
						"name": "Title",
						"type": "STRING"
					},
					{
						"name": "URI",
						"type": "STRING"
					}
					],
					"name": "External",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "Alt",
						"type": "STRING"
					},
					{
						"name": "BlobLink",
						"type": "STRING"
					},
					{
						"name": "Height",
						"type": "INTEGER"
					},
					{
						"name": "MimeType",
						"type": "STRING"
					},
					{
						"name": "Width",
						"type": "INTEGER"
					}
					],
					"mode": "REPEATED",
					"name": "Images",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "CID",
						"type": "STRING"
					},
					{
						"name": "Type",
						"type": "STRING"
					},
					{
						"name": "URI",
						"type": "STRING"
					}
					],
					"name": "Record",
					"type": "RECORD"
				}
				],
				"name": "Embed",
				"type": "RECORD"
			},
			{
				"mode": "REPEATED",
				"name": "Langs",
				"type": "STRING"
			},
			{
				"name": "ReplyParentCID",
				"type": "STRING"
			},
			{
				"mode": "REPEATED",
				"name": "Hashtags",
				"type": "STRING"
			},
			{
				"mode": "REPEATED",
				"name": "URLs",
				"type": "STRING"
			},
			{
				"name": "Text",
				"type": "STRING"
			}
			],
			"name": "Post",
			"type": "RECORD"
		},
		{
			"fields": [
			{
				"name": "Description",
				"type": "STRING"
			},
			{
				"name": "DisplayName",
				"type": "STRING"
			}
			],
			"name": "Profile",
			"type": "RECORD"
		},
		{
			"fields": [
			{
				"fields": [
				{
					"name": "Avatar",
					"type": "STRING"
				},
				{
					"name": "DID",
					"type": "STRING"
				},
				{
					"name": "DisplayName",
					"type": "STRING"
				},
				{
					"name": "Handle",
					"type": "STRING"
				},
				{
					"name": "IndexedAt",
					"type": "TIMESTAMP"
				}
				],
				"name": "Author",
				"type": "RECORD"
			},
			{
				"name": "CID",
				"type": "STRING"
			},
			{
				"name": "CreatedAt",
				"type": "TIMESTAMP"
			},
			{
				"fields": [
				{
					"fields": [
					{
						"name": "Alt",
						"type": "STRING"
					},
					{
						"name": "BlobLink",
						"type": "STRING"
					},
					{
						"name": "Height",
						"type": "INTEGER"
					},
					{
						"name": "MimeType",
						"type": "STRING"
					},
					{
						"name": "Width",
						"type": "INTEGER"
					}
					],
					"mode": "REPEATED",
					"name": "EmbedRecordMedia",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "Description",
						"type": "STRING"
					},
					{
						"name": "Title",
						"type": "STRING"
					},
					{
						"name": "URI",
						"type": "STRING"
					}
					],
					"name": "External",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "Alt",
						"type": "STRING"
					},
					{
						"name": "BlobLink",
						"type": "STRING"
					},
					{
						"name": "Height",
						"type": "INTEGER"
					},
					{
						"name": "MimeType",
						"type": "STRING"
					},
					{
						"name": "Width",
						"type": "INTEGER"
					}
					],
					"mode": "REPEATED",
					"name": "Images",
					"type": "RECORD"
				},
				{
					"fields": [
					{
						"name": "CID",
						"type": "STRING"
					},
					{
						"name": "Type",
						"type": "STRING"
					},
					{
						"name": "URI",
						"type": "STRING"
					}
					],
					"name": "Record",
					"type": "RECORD"
				}
				],
				"name": "Embed",
				"type": "RECORD"
			},
			{
				"mode": "REPEATED",
				"name": "Langs",
				"type": "STRING"
			},
			{
				"name": "LikeCount",
				"type": "INTEGER"
			},
			{
				"name": "RepostCount",
				"type": "INTEGER"
			},
			{
				"mode": "REPEATED",
				"name": "Hashtags",
				"type": "STRING"
			},
			{
				"mode": "REPEATED",
				"name": "URLs",
				"type": "STRING"
			},
			{
				"name": "Text",
				"type": "STRING"
			},
			{
				"name": "URI",
				"type": "STRING"
			}
			],
			"name": "RepostedPost",
			"type": "RECORD"
		}
		],
		"name": "Projection",
		"type": "RECORD"
	},
	{
		"name": "PulledTimestamp",
		"type": "TIMESTAMP"
	},
	{
		"name": "Seq",
		"type": "INTEGER"
	},
	{
		"name": "Type",
		"type": "STRING"
	}
  ]`

	schema, error := bigquery.SchemaFromJSON([]byte(rawSchema))

	if error != nil {
		log.Fatalf("unable to parse provided JSON BigQuery schema: %+v", error)
	}

	return schema
}
