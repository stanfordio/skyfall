package schema

import (
	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/bigquery"
)

func GetSchema() bigquery.Schema {
	// Generate this automatically from exported files using https://pypi.org/project/bigquery-schema-generator/
	// E.g., generate-schema < file.data.json
	rawSchema := `[
    {
      "mode": "NULLABLE",
      "name": "Action",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "CreatedAt",
      "type": "TIMESTAMP"
    },
    {
      "mode": "NULLABLE",
      "name": "Full",
      "type": "STRING"
    },
    {
      "fields": [
        {
          "fields": [
            {
              "mode": "NULLABLE",
              "name": "DID",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "DIDKey",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "Handle",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "PDS",
              "type": "STRING"
            }
          ],
          "mode": "NULLABLE",
          "name": "Actor",
          "type": "RECORD"
        },
        {
          "fields": [
            {
              "mode": "NULLABLE",
              "name": "Avatar",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "Description",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "DID",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "DisplayName",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "IndexedAt",
              "type": "TIMESTAMP"
            },
            {
              "mode": "NULLABLE",
              "name": "FollowersCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "FollowsCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "Handle",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "PostsCount",
              "type": "INTEGER"
            }
          ],
          "mode": "NULLABLE",
          "name": "BlockedProfile",
          "type": "RECORD"
        },
        {
          "fields": [
            {
              "mode": "NULLABLE",
              "name": "Avatar",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "Description",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "DID",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "DisplayName",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "IndexedAt",
              "type": "TIMESTAMP"
            },
            {
              "mode": "NULLABLE",
              "name": "FollowersCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "FollowsCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "Handle",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "PostsCount",
              "type": "INTEGER"
            }
          ],
          "mode": "NULLABLE",
          "name": "FollowedProfile",
          "type": "RECORD"
        },
        {
          "fields": [
            {
              "fields": [
                {
                  "mode": "NULLABLE",
                  "name": "Avatar",
                  "type": "STRING"
                },
                {
                  "mode": "NULLABLE",
                  "name": "DID",
                  "type": "STRING"
                },
                {
                  "mode": "NULLABLE",
                  "name": "DisplayName",
                  "type": "STRING"
                },
                {
                  "mode": "NULLABLE",
                  "name": "Handle",
                  "type": "STRING"
                }
              ],
              "mode": "NULLABLE",
              "name": "Author",
              "type": "RECORD"
            },
            {
              "mode": "NULLABLE",
              "name": "CID",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "CreatedAt",
              "type": "TIMESTAMP"
            },
            {
              "fields": [
                {
                  "fields": [
                    {
                      "mode": "NULLABLE",
                      "name": "Alt",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "BlobLink",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Height",
                      "type": "INTEGER"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "MimeType",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
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
                      "mode": "NULLABLE",
                      "name": "Description",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Title",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "URI",
                      "type": "STRING"
                    }
                  ],
                  "mode": "NULLABLE",
                  "name": "External",
                  "type": "RECORD"
                },
                {
                  "fields": [
                    {
                      "mode": "NULLABLE",
                      "name": "Alt",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "BlobLink",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Height",
                      "type": "INTEGER"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "MimeType",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
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
                      "mode": "NULLABLE",
                      "name": "CID",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Type",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "URI",
                      "type": "STRING"
                    }
                  ],
                  "mode": "NULLABLE",
                  "name": "Record",
                  "type": "RECORD"
                }
              ],
              "mode": "NULLABLE",
              "name": "Embed",
              "type": "RECORD"
            },
            {
              "mode": "REPEATED",
              "name": "Langs",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "LikeCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "ReplyCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "RepostCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "Text",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "URI",
              "type": "STRING"
            }
          ],
          "mode": "NULLABLE",
          "name": "LikedPost",
          "type": "RECORD"
        },
        {
          "fields": [
            {
              "mode": "NULLABLE",
              "name": "CreatedAt",
              "type": "TIMESTAMP"
            },
            {
              "fields": [
                {
                  "fields": [
                    {
                      "mode": "NULLABLE",
                      "name": "Alt",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "BlobLink",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Height",
                      "type": "INTEGER"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "MimeType",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
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
                      "mode": "NULLABLE",
                      "name": "Description",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Title",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "URI",
                      "type": "STRING"
                    }
                  ],
                  "mode": "NULLABLE",
                  "name": "External",
                  "type": "RECORD"
                },
                {
                  "fields": [
                    {
                      "mode": "NULLABLE",
                      "name": "Alt",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "BlobLink",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Height",
                      "type": "INTEGER"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "MimeType",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
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
                      "mode": "NULLABLE",
                      "name": "CID",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Type",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "URI",
                      "type": "STRING"
                    }
                  ],
                  "mode": "NULLABLE",
                  "name": "Record",
                  "type": "RECORD"
                }
              ],
              "mode": "NULLABLE",
              "name": "Embed",
              "type": "RECORD"
            },
            {
              "mode": "REPEATED",
              "name": "Langs",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "ReplyParentCID",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "Text",
              "type": "STRING"
            }
          ],
          "mode": "NULLABLE",
          "name": "Post",
          "type": "RECORD"
        },
        {
          "fields": [
            {
              "mode": "NULLABLE",
              "name": "Description",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "DisplayName",
              "type": "STRING"
            }
          ],
          "mode": "NULLABLE",
          "name": "Profile",
          "type": "RECORD"
        },
        {
          "fields": [
            {
              "fields": [
                {
                  "mode": "NULLABLE",
                  "name": "Avatar",
                  "type": "STRING"
                },
                {
                  "mode": "NULLABLE",
                  "name": "DID",
                  "type": "STRING"
                },
                {
                  "mode": "NULLABLE",
                  "name": "DisplayName",
                  "type": "STRING"
                },
                {
                  "mode": "NULLABLE",
                  "name": "Handle",
                  "type": "STRING"
                }
              ],
              "mode": "NULLABLE",
              "name": "Author",
              "type": "RECORD"
            },
            {
              "mode": "NULLABLE",
              "name": "CID",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "CreatedAt",
              "type": "TIMESTAMP"
            },
            {
              "fields": [
                {
                  "fields": [
                    {
                      "mode": "NULLABLE",
                      "name": "Alt",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "BlobLink",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Height",
                      "type": "INTEGER"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "MimeType",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
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
                      "mode": "NULLABLE",
                      "name": "Description",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Title",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "URI",
                      "type": "STRING"
                    }
                  ],
                  "mode": "NULLABLE",
                  "name": "External",
                  "type": "RECORD"
                },
                {
                  "fields": [
                    {
                      "mode": "NULLABLE",
                      "name": "Alt",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "BlobLink",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Height",
                      "type": "INTEGER"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "MimeType",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
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
                      "mode": "NULLABLE",
                      "name": "CID",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "Type",
                      "type": "STRING"
                    },
                    {
                      "mode": "NULLABLE",
                      "name": "URI",
                      "type": "STRING"
                    }
                  ],
                  "mode": "NULLABLE",
                  "name": "Record",
                  "type": "RECORD"
                }
              ],
              "mode": "NULLABLE",
              "name": "Embed",
              "type": "RECORD"
            },
            {
              "mode": "REPEATED",
              "name": "Langs",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "LikeCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "RepostCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "ReplyCount",
              "type": "INTEGER"
            },
            {
              "mode": "NULLABLE",
              "name": "Text",
              "type": "STRING"
            },
            {
              "mode": "NULLABLE",
              "name": "URI",
              "type": "STRING"
            }
          ],
          "mode": "NULLABLE",
          "name": "RepostedPost",
          "type": "RECORD"
        }
      ],
      "mode": "NULLABLE",
      "name": "Projection",
      "type": "RECORD"
    },
    {
      "mode": "NULLABLE",
      "name": "PulledTimestamp",
      "type": "TIMESTAMP"
    },
    {
      "mode": "NULLABLE",
      "name": "Seq",
      "type": "INTEGER"
    },
    {
      "mode": "NULLABLE",
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
