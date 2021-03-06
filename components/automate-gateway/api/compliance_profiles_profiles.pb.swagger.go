package api

func init() {
	Swagger.Add("compliance_profiles_profiles", `{
  "swagger": "2.0",
  "info": {
    "title": "components/automate-gateway/api/compliance/profiles/profiles.proto",
    "version": "version not set"
  },
  "schemes": [
    "http",
    "https"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/compliance/market/read/{name}/version/{version}": {
      "get": {
        "operationId": "ReadFromMarket",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/v1Profile"
            }
          }
        },
        "parameters": [
          {
            "name": "name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "version",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "owner",
            "in": "query",
            "required": false,
            "type": "string"
          }
        ],
        "tags": [
          "ProfilesService"
        ]
      }
    },
    "/compliance/profiles/read/{owner}/{name}/version/{version}": {
      "get": {
        "operationId": "Read",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/v1Profile"
            }
          }
        },
        "parameters": [
          {
            "name": "owner",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "version",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "ProfilesService"
        ]
      }
    },
    "/compliance/profiles/search": {
      "post": {
        "operationId": "List",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/v1Profiles"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/v1Query"
            }
          }
        ],
        "tags": [
          "ProfilesService"
        ]
      }
    },
    "/compliance/profiles/{owner}/{name}/version/{version}": {
      "delete": {
        "operationId": "Delete",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "properties": {}
            }
          }
        },
        "parameters": [
          {
            "name": "owner",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "version",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "tags": [
          "ProfilesService"
        ]
      }
    }
  },
  "definitions": {
    "QueryOrderType": {
      "type": "string",
      "enum": [
        "ASC",
        "DESC"
      ],
      "default": "ASC"
    },
    "protobufAny": {
      "type": "object",
      "properties": {
        "type_url": {
          "type": "string"
        },
        "value": {
          "type": "string",
          "format": "byte"
        }
      }
    },
    "runtimeStreamError": {
      "type": "object",
      "properties": {
        "grpc_code": {
          "type": "integer",
          "format": "int32"
        },
        "http_code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        },
        "http_status": {
          "type": "string"
        },
        "details": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/protobufAny"
          }
        }
      }
    },
    "v1Attribute": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "options": {
          "$ref": "#/definitions/v1Option"
        }
      }
    },
    "v1CheckMessage": {
      "type": "object",
      "properties": {
        "file": {
          "type": "string"
        },
        "line": {
          "type": "integer",
          "format": "int32"
        },
        "column": {
          "type": "integer",
          "format": "int32"
        },
        "control_id": {
          "type": "string"
        },
        "msg": {
          "type": "string"
        }
      }
    },
    "v1CheckResult": {
      "type": "object",
      "properties": {
        "summary": {
          "$ref": "#/definitions/v1ResultSummary"
        },
        "errors": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1CheckMessage"
          }
        },
        "warnings": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1CheckMessage"
          }
        }
      }
    },
    "v1Chunk": {
      "type": "object",
      "properties": {
        "data": {
          "type": "string",
          "format": "byte"
        },
        "position": {
          "type": "string",
          "format": "int64"
        }
      }
    },
    "v1Control": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "code": {
          "type": "string"
        },
        "desc": {
          "type": "string"
        },
        "impact": {
          "type": "number",
          "format": "float"
        },
        "title": {
          "type": "string"
        },
        "source_location": {
          "$ref": "#/definitions/v1SourceLocation"
        },
        "results": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Result"
          }
        },
        "refs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Ref"
          }
        },
        "tags": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "v1Dependency": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "url": {
          "type": "string"
        },
        "path": {
          "type": "string"
        },
        "git": {
          "type": "string"
        },
        "branch": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        },
        "commit": {
          "type": "string"
        },
        "version": {
          "type": "string"
        },
        "supermarket": {
          "type": "string"
        },
        "github": {
          "type": "string"
        },
        "compliance": {
          "type": "string"
        }
      }
    },
    "v1Group": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "title": {
          "type": "string"
        },
        "controls": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "v1ListFilter": {
      "type": "object",
      "properties": {
        "values": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "type": {
          "type": "string"
        }
      }
    },
    "v1Metadata": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "version": {
          "type": "string"
        },
        "content_type": {
          "type": "string"
        }
      }
    },
    "v1Option": {
      "type": "object",
      "properties": {
        "description": {
          "type": "string"
        },
        "default": {
          "type": "string"
        }
      }
    },
    "v1Profile": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "title": {
          "type": "string"
        },
        "maintainer": {
          "type": "string"
        },
        "copyright": {
          "type": "string"
        },
        "copyright_email": {
          "type": "string"
        },
        "license": {
          "type": "string"
        },
        "summary": {
          "type": "string"
        },
        "version": {
          "type": "string"
        },
        "owner": {
          "type": "string"
        },
        "supports": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Support"
          }
        },
        "depends": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Dependency"
          }
        },
        "sha256": {
          "type": "string"
        },
        "groups": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Group"
          }
        },
        "controls": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Control"
          }
        },
        "attributes": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Attribute"
          }
        },
        "latest_version": {
          "type": "string"
        }
      }
    },
    "v1ProfileData": {
      "type": "object",
      "properties": {
        "owner": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "version": {
          "type": "string"
        },
        "data": {
          "type": "string",
          "format": "byte"
        }
      }
    },
    "v1Profiles": {
      "type": "object",
      "properties": {
        "profiles": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1Profile"
          }
        },
        "total": {
          "type": "integer",
          "format": "int32"
        }
      }
    },
    "v1Query": {
      "type": "object",
      "properties": {
        "filters": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/v1ListFilter"
          }
        },
        "order": {
          "$ref": "#/definitions/QueryOrderType"
        },
        "sort": {
          "type": "string"
        },
        "page": {
          "type": "integer",
          "format": "int32"
        },
        "per_page": {
          "type": "integer",
          "format": "int32"
        },
        "owner": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "version": {
          "type": "string"
        }
      }
    },
    "v1Ref": {
      "type": "object",
      "properties": {
        "url": {
          "type": "string"
        },
        "ref": {
          "type": "string"
        }
      }
    },
    "v1Result": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string"
        },
        "code_desc": {
          "type": "string"
        },
        "run_time": {
          "type": "number",
          "format": "float"
        },
        "start_time": {
          "type": "string"
        },
        "message": {
          "type": "string"
        },
        "skip_message": {
          "type": "string"
        }
      }
    },
    "v1ResultSummary": {
      "type": "object",
      "properties": {
        "valid": {
          "type": "boolean",
          "format": "boolean"
        },
        "timestamp": {
          "type": "string"
        },
        "location": {
          "type": "string"
        },
        "controls": {
          "type": "integer",
          "format": "int32"
        }
      }
    },
    "v1SourceLocation": {
      "type": "object",
      "properties": {
        "ref": {
          "type": "string"
        },
        "line": {
          "type": "integer",
          "format": "int32"
        }
      }
    },
    "v1Support": {
      "type": "object",
      "properties": {
        "os_name": {
          "type": "string"
        },
        "os_family": {
          "type": "string"
        },
        "release": {
          "type": "string"
        },
        "inspec_version": {
          "type": "string"
        },
        "platform": {
          "type": "string"
        }
      }
    }
  },
  "x-stream-definitions": {
    "v1ProfileData": {
      "type": "object",
      "properties": {
        "result": {
          "$ref": "#/definitions/v1ProfileData"
        },
        "error": {
          "$ref": "#/definitions/runtimeStreamError"
        }
      },
      "title": "Stream result of v1ProfileData"
    }
  }
}
`)
}
