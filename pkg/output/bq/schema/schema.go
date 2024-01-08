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
    "name": "_Action",
    "type": "STRING"
  },
  {
    "mode": "NULLABLE",
    "name": "_Raw",
    "type": "STRING"
  },
  {
    "mode": "NULLABLE",
    "name": "_ActorDid",
    "type": "STRING"
  },
  {
    "fields": [
      {
        "mode": "REPEATED",
        "name": "AlsoKnownAs",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "DID",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "Handle",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "PublicKeyMultibase",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "Type",
                "type": "STRING"
              }
            ],
            "mode": "NULLABLE",
            "name": "atproto",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "Keys",
        "type": "RECORD"
      },
      {
        "fields": [
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "Type",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "URL",
                "type": "STRING"
              }
            ],
            "mode": "NULLABLE",
            "name": "atproto_pds",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "Services",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "_ActorIdentity",
    "type": "RECORD"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "avatar",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "banner",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "description",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "did",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "displayName",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "followersCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "followsCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "handle",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "indexedAt",
        "type": "TIMESTAMP"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cts",
            "type": "TIMESTAMP"
          },
          {
            "mode": "NULLABLE",
            "name": "neg",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "src",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "val",
            "type": "STRING"
          }
        ],
        "mode": "REPEATED",
        "name": "labels",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "postsCount",
        "type": "INTEGER"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "blockedBy",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "muted",
            "type": "BOOLEAN"
          }
        ],
        "mode": "NULLABLE",
        "name": "viewer",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "_ActorProfile",
    "type": "RECORD"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "avatar",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "banner",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "description",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "did",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "displayName",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "followersCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "followsCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "handle",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "indexedAt",
        "type": "TIMESTAMP"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cts",
            "type": "TIMESTAMP"
          },
          {
            "mode": "NULLABLE",
            "name": "neg",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "src",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "val",
            "type": "STRING"
          }
        ],
        "mode": "REPEATED",
        "name": "labels",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "postsCount",
        "type": "INTEGER"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "blockedBy",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "muted",
            "type": "BOOLEAN"
          }
        ],
        "mode": "NULLABLE",
        "name": "viewer",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "_BlockedProfile",
    "type": "RECORD"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "avatar",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "banner",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "description",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "did",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "displayName",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "followersCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "followsCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "handle",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "indexedAt",
        "type": "TIMESTAMP"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cts",
            "type": "TIMESTAMP"
          },
          {
            "mode": "NULLABLE",
            "name": "neg",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "src",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "val",
            "type": "STRING"
          }
        ],
        "mode": "REPEATED",
        "name": "labels",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "postsCount",
        "type": "INTEGER"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "blockedBy",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "muted",
            "type": "BOOLEAN"
          }
        ],
        "mode": "NULLABLE",
        "name": "viewer",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "_FollowedProfile",
    "type": "RECORD"
  },
  {
    "mode": "NULLABLE",
    "name": "_Item",
    "type": "STRING"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "avatar",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "did",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "displayName",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "handle",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "cid",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "cts",
                "type": "TIMESTAMP"
              },
              {
                "mode": "NULLABLE",
                "name": "neg",
                "type": "BOOLEAN"
              },
              {
                "mode": "NULLABLE",
                "name": "src",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "val",
                "type": "STRING"
              }
            ],
            "mode": "REPEATED",
            "name": "labels",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "blockedBy",
                "type": "BOOLEAN"
              },
              {
                "mode": "NULLABLE",
                "name": "muted",
                "type": "BOOLEAN"
              }
            ],
            "mode": "NULLABLE",
            "name": "viewer",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "author",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "cid",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "description",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "thumb",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "title",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              }
            ],
            "mode": "NULLABLE",
            "name": "external",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "alt",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "height",
                    "type": "INTEGER"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "width",
                    "type": "INTEGER"
                  }
                ],
                "mode": "NULLABLE",
                "name": "aspectRatio",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "fullsize",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "thumb",
                "type": "STRING"
              }
            ],
            "mode": "REPEATED",
            "name": "images",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "description",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "title",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "external",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "alt",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "height",
                        "type": "INTEGER"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "width",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "aspectRatio",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "fullsize",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "thumb",
                    "type": "STRING"
                  }
                ],
                "mode": "REPEATED",
                "name": "images",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "media",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "avatar",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "did",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "displayName",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "handle",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "cid",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "cts",
                        "type": "TIMESTAMP"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "neg",
                        "type": "BOOLEAN"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "src",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "val",
                        "type": "STRING"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "labels",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "blockedBy",
                        "type": "BOOLEAN"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "muted",
                        "type": "BOOLEAN"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "viewer",
                    "type": "RECORD"
                  }
                ],
                "mode": "NULLABLE",
                "name": "author",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "cid",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "description",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "thumb",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "title",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "external",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "alt",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "height",
                            "type": "INTEGER"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "width",
                            "type": "INTEGER"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "aspectRatio",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "fullsize",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "thumb",
                        "type": "STRING"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "images",
                    "type": "RECORD"
                  }
                ],
                "mode": "REPEATED",
                "name": "embeds",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "indexedAt",
                "type": "TIMESTAMP"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "avatar",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "did",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "displayName",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "handle",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "cid",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "cts",
                            "type": "TIMESTAMP"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "neg",
                            "type": "BOOLEAN"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "src",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "val",
                            "type": "STRING"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "labels",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "blockedBy",
                            "type": "BOOLEAN"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "muted",
                            "type": "BOOLEAN"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "viewer",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "author",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "alt",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "height",
                                "type": "INTEGER"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "width",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "aspectRatio",
                            "type": "RECORD"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "fullsize",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "thumb",
                            "type": "STRING"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "images",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "alt",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "height",
                                    "type": "INTEGER"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "width",
                                    "type": "INTEGER"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "aspectRatio",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "fullsize",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "thumb",
                                "type": "STRING"
                              }
                            ],
                            "mode": "REPEATED",
                            "name": "images",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "media",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "avatar",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "did",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "displayName",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "handle",
                                    "type": "STRING"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "blockedBy",
                                        "type": "BOOLEAN"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "muted",
                                        "type": "BOOLEAN"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "viewer",
                                    "type": "RECORD"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "author",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "cid",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "indexedAt",
                                "type": "TIMESTAMP"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "cid",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "cts",
                                    "type": "TIMESTAMP"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "neg",
                                    "type": "BOOLEAN"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "src",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "uri",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "val",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "REPEATED",
                                "name": "labels",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "uri",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "_type",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "createdAt",
                                    "type": "TIMESTAMP"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "_type",
                                        "type": "STRING"
                                      },
                                      {
                                        "fields": [
                                          {
                                            "mode": "NULLABLE",
                                            "name": "alt",
                                            "type": "STRING"
                                          },
                                          {
                                            "fields": [
                                              {
                                                "mode": "NULLABLE",
                                                "name": "height",
                                                "type": "INTEGER"
                                              },
                                              {
                                                "mode": "NULLABLE",
                                                "name": "width",
                                                "type": "INTEGER"
                                              }
                                            ],
                                            "mode": "NULLABLE",
                                            "name": "aspectRatio",
                                            "type": "RECORD"
                                          },
                                          {
                                            "fields": [
                                              {
                                                "mode": "NULLABLE",
                                                "name": "_type",
                                                "type": "STRING"
                                              },
                                              {
                                                "mode": "NULLABLE",
                                                "name": "mimeType",
                                                "type": "STRING"
                                              },
                                              {
                                                "fields": [
                                                  {
                                                    "mode": "NULLABLE",
                                                    "name": "_link",
                                                    "type": "STRING"
                                                  }
                                                ],
                                                "mode": "NULLABLE",
                                                "name": "ref",
                                                "type": "RECORD"
                                              },
                                              {
                                                "mode": "NULLABLE",
                                                "name": "size",
                                                "type": "INTEGER"
                                              }
                                            ],
                                            "mode": "NULLABLE",
                                            "name": "image",
                                            "type": "RECORD"
                                          }
                                        ],
                                        "mode": "REPEATED",
                                        "name": "images",
                                        "type": "RECORD"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "embed",
                                    "type": "RECORD"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "_type",
                                        "type": "STRING"
                                      },
                                      {
                                        "fields": [
                                          {
                                            "mode": "NULLABLE",
                                            "name": "val",
                                            "type": "STRING"
                                          }
                                        ],
                                        "mode": "REPEATED",
                                        "name": "values",
                                        "type": "RECORD"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "labels",
                                    "type": "RECORD"
                                  },
                                  {
                                    "mode": "REPEATED",
                                    "name": "langs",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "text",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "value",
                                "type": "RECORD"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "record",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "record",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "embeds",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "indexedAt",
                    "type": "TIMESTAMP"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "cid",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "cts",
                        "type": "TIMESTAMP"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "neg",
                        "type": "BOOLEAN"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "src",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "val",
                        "type": "STRING"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "labels",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "createdAt",
                        "type": "TIMESTAMP"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "alt",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "height",
                                    "type": "INTEGER"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "width",
                                    "type": "INTEGER"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "aspectRatio",
                                "type": "RECORD"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "_type",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "mimeType",
                                    "type": "STRING"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "_link",
                                        "type": "STRING"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "ref",
                                    "type": "RECORD"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "size",
                                    "type": "INTEGER"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "image",
                                "type": "RECORD"
                              }
                            ],
                            "mode": "REPEATED",
                            "name": "images",
                            "type": "RECORD"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "alt",
                                    "type": "STRING"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "height",
                                        "type": "INTEGER"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "width",
                                        "type": "INTEGER"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "aspectRatio",
                                    "type": "RECORD"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "_type",
                                        "type": "STRING"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "mimeType",
                                        "type": "STRING"
                                      },
                                      {
                                        "fields": [
                                          {
                                            "mode": "NULLABLE",
                                            "name": "_link",
                                            "type": "STRING"
                                          }
                                        ],
                                        "mode": "NULLABLE",
                                        "name": "ref",
                                        "type": "RECORD"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "size",
                                        "type": "INTEGER"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "image",
                                    "type": "RECORD"
                                  }
                                ],
                                "mode": "REPEATED",
                                "name": "images",
                                "type": "RECORD"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "media",
                            "type": "RECORD"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "cid",
                                    "type": "STRING"
                                  },
                                  {
                                    "mode": "NULLABLE",
                                    "name": "uri",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "record",
                                "type": "RECORD"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "record",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "embed",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "did",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "tag",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "uri",
                                "type": "STRING"
                              }
                            ],
                            "mode": "REPEATED",
                            "name": "features",
                            "type": "RECORD"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "byteEnd",
                                "type": "INTEGER"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "byteStart",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "index",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "facets",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "val",
                                "type": "STRING"
                              }
                            ],
                            "mode": "REPEATED",
                            "name": "values",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "labels",
                        "type": "RECORD"
                      },
                      {
                        "mode": "REPEATED",
                        "name": "langs",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "text",
                        "type": "STRING"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "value",
                    "type": "RECORD"
                  }
                ],
                "mode": "NULLABLE",
                "name": "record",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "createdAt",
                    "type": "TIMESTAMP"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "description",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "mimeType",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "_link",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "ref",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "size",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "thumb",
                            "type": "RECORD"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "title",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "external",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "alt",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "height",
                                "type": "INTEGER"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "width",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "aspectRatio",
                            "type": "RECORD"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "mimeType",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "_link",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "ref",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "size",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "image",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "images",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "embed",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "did",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "tag",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "features",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "byteEnd",
                            "type": "INTEGER"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "byteStart",
                            "type": "INTEGER"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "index",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "facets",
                    "type": "RECORD"
                  },
                  {
                    "mode": "REPEATED",
                    "name": "langs",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "cid",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "parent",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "cid",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "root",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "reply",
                    "type": "RECORD"
                  },
                  {
                    "mode": "REPEATED",
                    "name": "tags",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "text",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "value",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "record",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "embed",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "indexedAt",
        "type": "TIMESTAMP"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cts",
            "type": "TIMESTAMP"
          },
          {
            "mode": "NULLABLE",
            "name": "neg",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "src",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "val",
            "type": "STRING"
          }
        ],
        "mode": "REPEATED",
        "name": "labels",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "likeCount",
        "type": "INTEGER"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "createdAt",
            "type": "TIMESTAMP"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "description",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "mimeType",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_link",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "ref",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "size",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "thumb",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "title",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "external",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "alt",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "height",
                        "type": "INTEGER"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "width",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "aspectRatio",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "mimeType",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_link",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "ref",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "size",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "image",
                    "type": "RECORD"
                  }
                ],
                "mode": "REPEATED",
                "name": "images",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "description",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "title",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "external",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "alt",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "height",
                            "type": "INTEGER"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "width",
                            "type": "INTEGER"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "aspectRatio",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "mimeType",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_link",
                                "type": "STRING"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "ref",
                            "type": "RECORD"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "size",
                            "type": "INTEGER"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "image",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "images",
                    "type": "RECORD"
                  }
                ],
                "mode": "NULLABLE",
                "name": "media",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "cid",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "record",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "record",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "embed",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "did",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "tag",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "REPEATED",
                "name": "features",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "byteEnd",
                    "type": "INTEGER"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "byteStart",
                    "type": "INTEGER"
                  }
                ],
                "mode": "NULLABLE",
                "name": "index",
                "type": "RECORD"
              }
            ],
            "mode": "REPEATED",
            "name": "facets",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "val",
                    "type": "STRING"
                  }
                ],
                "mode": "REPEATED",
                "name": "values",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "labels",
            "type": "RECORD"
          },
          {
            "mode": "REPEATED",
            "name": "langs",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "parent",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "root",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "reply",
            "type": "RECORD"
          },
          {
            "mode": "NULLABLE",
            "name": "text",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "record",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "replyCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "repostCount",
        "type": "INTEGER"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  }
                ],
                "mode": "REPEATED",
                "name": "allow",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "createdAt",
                "type": "TIMESTAMP"
              },
              {
                "mode": "NULLABLE",
                "name": "post",
                "type": "STRING"
              }
            ],
            "mode": "NULLABLE",
            "name": "record",
            "type": "RECORD"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "threadgate",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "uri",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "replyDisabled",
            "type": "BOOLEAN"
          }
        ],
        "mode": "NULLABLE",
        "name": "viewer",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "_LikedPost",
    "type": "RECORD"
  },
  {
    "mode": "NULLABLE",
    "name": "_PulledTimestamp",
    "type": "TIMESTAMP"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "avatar",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "did",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "displayName",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "handle",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "cid",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "cts",
                "type": "TIMESTAMP"
              },
              {
                "mode": "NULLABLE",
                "name": "neg",
                "type": "BOOLEAN"
              },
              {
                "mode": "NULLABLE",
                "name": "src",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "val",
                "type": "STRING"
              }
            ],
            "mode": "REPEATED",
            "name": "labels",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "blockedBy",
                "type": "BOOLEAN"
              },
              {
                "mode": "NULLABLE",
                "name": "muted",
                "type": "BOOLEAN"
              }
            ],
            "mode": "NULLABLE",
            "name": "viewer",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "author",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "cid",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "description",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "thumb",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "title",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              }
            ],
            "mode": "NULLABLE",
            "name": "external",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "alt",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "height",
                    "type": "INTEGER"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "width",
                    "type": "INTEGER"
                  }
                ],
                "mode": "NULLABLE",
                "name": "aspectRatio",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "fullsize",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "thumb",
                "type": "STRING"
              }
            ],
            "mode": "REPEATED",
            "name": "images",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "avatar",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "did",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "displayName",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "handle",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "blockedBy",
                        "type": "BOOLEAN"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "muted",
                        "type": "BOOLEAN"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "viewer",
                    "type": "RECORD"
                  }
                ],
                "mode": "NULLABLE",
                "name": "author",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "cid",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "description",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "thumb",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "title",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "external",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "alt",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "height",
                            "type": "INTEGER"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "width",
                            "type": "INTEGER"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "aspectRatio",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "fullsize",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "thumb",
                        "type": "STRING"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "images",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "avatar",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "did",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "displayName",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "handle",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "cid",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "cts",
                                "type": "TIMESTAMP"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "neg",
                                "type": "BOOLEAN"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "src",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "uri",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "val",
                                "type": "STRING"
                              }
                            ],
                            "mode": "REPEATED",
                            "name": "labels",
                            "type": "RECORD"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "blockedBy",
                                "type": "BOOLEAN"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "muted",
                                "type": "BOOLEAN"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "viewer",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "author",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "cid",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "indexedAt",
                        "type": "TIMESTAMP"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "uri",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "createdAt",
                            "type": "TIMESTAMP"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "alt",
                                    "type": "STRING"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "height",
                                        "type": "INTEGER"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "width",
                                        "type": "INTEGER"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "aspectRatio",
                                    "type": "RECORD"
                                  },
                                  {
                                    "fields": [
                                      {
                                        "mode": "NULLABLE",
                                        "name": "_type",
                                        "type": "STRING"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "mimeType",
                                        "type": "STRING"
                                      },
                                      {
                                        "fields": [
                                          {
                                            "mode": "NULLABLE",
                                            "name": "_link",
                                            "type": "STRING"
                                          }
                                        ],
                                        "mode": "NULLABLE",
                                        "name": "ref",
                                        "type": "RECORD"
                                      },
                                      {
                                        "mode": "NULLABLE",
                                        "name": "size",
                                        "type": "INTEGER"
                                      }
                                    ],
                                    "mode": "NULLABLE",
                                    "name": "image",
                                    "type": "RECORD"
                                  }
                                ],
                                "mode": "REPEATED",
                                "name": "images",
                                "type": "RECORD"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "embed",
                            "type": "RECORD"
                          },
                          {
                            "mode": "REPEATED",
                            "name": "langs",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "text",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "value",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "record",
                    "type": "RECORD"
                  }
                ],
                "mode": "REPEATED",
                "name": "embeds",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "indexedAt",
                "type": "TIMESTAMP"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "createdAt",
                    "type": "TIMESTAMP"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "description",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "mimeType",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "_link",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "ref",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "size",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "thumb",
                            "type": "RECORD"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "title",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "external",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "alt",
                            "type": "STRING"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "height",
                                "type": "INTEGER"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "width",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "aspectRatio",
                            "type": "RECORD"
                          },
                          {
                            "fields": [
                              {
                                "mode": "NULLABLE",
                                "name": "_type",
                                "type": "STRING"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "mimeType",
                                "type": "STRING"
                              },
                              {
                                "fields": [
                                  {
                                    "mode": "NULLABLE",
                                    "name": "_link",
                                    "type": "STRING"
                                  }
                                ],
                                "mode": "NULLABLE",
                                "name": "ref",
                                "type": "RECORD"
                              },
                              {
                                "mode": "NULLABLE",
                                "name": "size",
                                "type": "INTEGER"
                              }
                            ],
                            "mode": "NULLABLE",
                            "name": "image",
                            "type": "RECORD"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "images",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "cid",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "uri",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "record",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "embed",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_type",
                            "type": "STRING"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "did",
                            "type": "STRING"
                          }
                        ],
                        "mode": "REPEATED",
                        "name": "features",
                        "type": "RECORD"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "byteEnd",
                            "type": "INTEGER"
                          },
                          {
                            "mode": "NULLABLE",
                            "name": "byteStart",
                            "type": "INTEGER"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "index",
                        "type": "RECORD"
                      }
                    ],
                    "mode": "REPEATED",
                    "name": "facets",
                    "type": "RECORD"
                  },
                  {
                    "mode": "REPEATED",
                    "name": "langs",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "text",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "value",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "record",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "embed",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "indexedAt",
        "type": "TIMESTAMP"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cts",
            "type": "TIMESTAMP"
          },
          {
            "mode": "NULLABLE",
            "name": "neg",
            "type": "BOOLEAN"
          },
          {
            "mode": "NULLABLE",
            "name": "src",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "val",
            "type": "STRING"
          }
        ],
        "mode": "REPEATED",
        "name": "labels",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "likeCount",
        "type": "INTEGER"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "createdAt",
            "type": "TIMESTAMP"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "description",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "mimeType",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_link",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "ref",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "size",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "thumb",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "title",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "external",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "alt",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "height",
                        "type": "INTEGER"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "width",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "aspectRatio",
                    "type": "RECORD"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_type",
                        "type": "STRING"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "mimeType",
                        "type": "STRING"
                      },
                      {
                        "fields": [
                          {
                            "mode": "NULLABLE",
                            "name": "_link",
                            "type": "STRING"
                          }
                        ],
                        "mode": "NULLABLE",
                        "name": "ref",
                        "type": "RECORD"
                      },
                      {
                        "mode": "NULLABLE",
                        "name": "size",
                        "type": "INTEGER"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "image",
                    "type": "RECORD"
                  }
                ],
                "mode": "REPEATED",
                "name": "images",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "record",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "embed",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "tag",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "REPEATED",
                "name": "features",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "byteEnd",
                    "type": "INTEGER"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "byteStart",
                    "type": "INTEGER"
                  }
                ],
                "mode": "NULLABLE",
                "name": "index",
                "type": "RECORD"
              }
            ],
            "mode": "REPEATED",
            "name": "facets",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "val",
                    "type": "STRING"
                  }
                ],
                "mode": "REPEATED",
                "name": "values",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "labels",
            "type": "RECORD"
          },
          {
            "mode": "REPEATED",
            "name": "langs",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "parent",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "cid",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "uri",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "root",
                "type": "RECORD"
              }
            ],
            "mode": "NULLABLE",
            "name": "reply",
            "type": "RECORD"
          },
          {
            "mode": "NULLABLE",
            "name": "text",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "record",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "replyCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "repostCount",
        "type": "INTEGER"
      },
      {
        "mode": "NULLABLE",
        "name": "uri",
        "type": "STRING"
      }
    ],
    "mode": "NULLABLE",
    "name": "_RepostedPost",
    "type": "RECORD"
  },
  {
    "mode": "NULLABLE",
    "name": "_Seq",
    "type": "INTEGER"
  },
  {
    "mode": "NULLABLE",
    "name": "_Type",
    "type": "STRING"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "mimeType",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_link",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "ref",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "size",
        "type": "INTEGER"
      }
    ],
    "mode": "NULLABLE",
    "name": "Avatar",
    "type": "RECORD"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "mimeType",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_link",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "ref",
        "type": "RECORD"
      },
      {
        "mode": "NULLABLE",
        "name": "size",
        "type": "INTEGER"
      }
    ],
    "mode": "NULLABLE",
    "name": "Banner",
    "type": "RECORD"
  },
  {
    "mode": "NULLABLE",
    "name": "CreatedAt",
    "type": "TIMESTAMP"
  },
  {
    "mode": "NULLABLE",
    "name": "Description",
    "type": "STRING"
  },
  {
    "mode": "NULLABLE",
    "name": "Did",
    "type": "STRING"
  },
  {
    "mode": "NULLABLE",
    "name": "DisplayName",
    "type": "STRING"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "description",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "mimeType",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_link",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "ref",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "size",
                "type": "INTEGER"
              }
            ],
            "mode": "NULLABLE",
            "name": "thumb",
            "type": "RECORD"
          },
          {
            "mode": "NULLABLE",
            "name": "title",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "external",
        "type": "RECORD"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "alt",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "height",
                "type": "INTEGER"
              },
              {
                "mode": "NULLABLE",
                "name": "width",
                "type": "INTEGER"
              }
            ],
            "mode": "NULLABLE",
            "name": "aspectRatio",
            "type": "RECORD"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "_type",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "mimeType",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_link",
                    "type": "STRING"
                  }
                ],
                "mode": "NULLABLE",
                "name": "ref",
                "type": "RECORD"
              },
              {
                "mode": "NULLABLE",
                "name": "size",
                "type": "INTEGER"
              }
            ],
            "mode": "NULLABLE",
            "name": "image",
            "type": "RECORD"
          }
        ],
        "mode": "REPEATED",
        "name": "images",
        "type": "RECORD"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "alt",
                "type": "STRING"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "height",
                    "type": "INTEGER"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "width",
                    "type": "INTEGER"
                  }
                ],
                "mode": "NULLABLE",
                "name": "aspectRatio",
                "type": "RECORD"
              },
              {
                "fields": [
                  {
                    "mode": "NULLABLE",
                    "name": "_type",
                    "type": "STRING"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "mimeType",
                    "type": "STRING"
                  },
                  {
                    "fields": [
                      {
                        "mode": "NULLABLE",
                        "name": "_link",
                        "type": "STRING"
                      }
                    ],
                    "mode": "NULLABLE",
                    "name": "ref",
                    "type": "RECORD"
                  },
                  {
                    "mode": "NULLABLE",
                    "name": "size",
                    "type": "INTEGER"
                  }
                ],
                "mode": "NULLABLE",
                "name": "image",
                "type": "RECORD"
              }
            ],
            "mode": "REPEATED",
            "name": "images",
            "type": "RECORD"
          }
        ],
        "mode": "NULLABLE",
        "name": "media",
        "type": "RECORD"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "fields": [
              {
                "mode": "NULLABLE",
                "name": "cid",
                "type": "STRING"
              },
              {
                "mode": "NULLABLE",
                "name": "uri",
                "type": "STRING"
              }
            ],
            "mode": "NULLABLE",
            "name": "record",
            "type": "RECORD"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "record",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "Embed",
    "type": "RECORD"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "val",
            "type": "STRING"
          }
        ],
        "mode": "REPEATED",
        "name": "values",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "Labels",
    "type": "RECORD"
  },
  {
    "mode": "REPEATED",
    "name": "Langs",
    "type": "STRING"
  },
  {
    "mode": "NULLABLE",
    "name": "LexiconTypeID",
    "type": "STRING"
  },
  {
    "mode": "NULLABLE",
    "name": "List",
    "type": "STRING"
  },
  {
    "fields": [
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "parent",
        "type": "RECORD"
      },
      {
        "fields": [
          {
            "mode": "NULLABLE",
            "name": "_type",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "cid",
            "type": "STRING"
          },
          {
            "mode": "NULLABLE",
            "name": "uri",
            "type": "STRING"
          }
        ],
        "mode": "NULLABLE",
        "name": "root",
        "type": "RECORD"
      }
    ],
    "mode": "NULLABLE",
    "name": "Reply",
    "type": "RECORD"
  },
  {
    "fields": [
      {
        "mode": "NULLABLE",
        "name": "_type",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "cid",
        "type": "STRING"
      },
      {
        "mode": "NULLABLE",
        "name": "uri",
        "type": "STRING"
      }
    ],
    "mode": "NULLABLE",
    "name": "Subject",
    "type": "RECORD"
  },
  {
    "mode": "NULLABLE",
    "name": "Text",
    "type": "STRING"
  }
]`

	schema, error := bigquery.SchemaFromJSON([]byte(rawSchema))

	if error != nil {
		log.Fatalf("unable to parse provided JSON BigQuery schema: %+v", error)
	}

	return schema
}
