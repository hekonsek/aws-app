package cmd

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file8 := &embedded.EmbeddedFile{
		Filename:    "Dockerfile",
		FileModTime: time.Unix(1552547852, 0),

		Content: string("FROM openjdk:11-jre\n\nRUN mkdir /app\nADD target/*.jar /app\n\nCMD [\"/usr/bin/java\", \"-jar\", \"/app/*.jar\"]\n"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "assume_service_role_template.json",
		FileModTime: time.Unix(1551990782, 0),

		Content: string("{\n  \"Version\":\"2012-10-17\",\n  \"Statement\":[\n    {\n      \"Effect\":\"Allow\",\n      \"Principal\":{\n        \"Service\":\"{{ .Service }}\"\n      },\n      \"Action\":\"sts:AssumeRole\"\n    }\n  ]\n}"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "buildspec-build.yml",
		FileModTime: time.Unix(1552337535, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - mvn test package\nartifacts:\n  files:\n    - \"**/*\""),
	}
	fileb := &embedded.EmbeddedFile{
		Filename:    "buildspec-dockerize.yml",
		FileModTime: time.Unix(1552549826, 0),

		Content: string("version: 0.2\n\nphases:\n  pre_build:\n    commands:\n      - $(aws ecr get-login --region $AWS_REGION --no-include-email)\n      - REPOSITORY=`./awsom step ecr`\n      - VERSION=`./awsom step version-current`\n  build:\n    commands:\n      - docker build -t $REPOSITORY:$VERSION .\n      - docker push $REPOSITORY:$VERSION\nartifacts:\n  files:\n    - imageDetail.json"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "buildspec-version.yml",
		FileModTime: time.Unix(1552563244, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - cp /usr/bin/awsom awsom\n      - ./awsom app configure\n      - ./awsom step clone\n      - ./awsom step version\nartifacts:\n  files:\n    - \"**/*\""),
	}

	// define dirs
	dir7 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1552563244, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file8, // "Dockerfile"
			file9, // "assume_service_role_template.json"
			filea, // "buildspec-build.yml"
			fileb, // "buildspec-dockerize.yml"
			filec, // "buildspec-version.yml"

		},
	}

	// link ChildDirs
	dir7.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`../rice`, &embedded.EmbeddedBox{
		Name: `../rice`,
		Time: time.Unix(1552563244, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir7,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"Dockerfile":                        file8,
			"assume_service_role_template.json": file9,
			"buildspec-build.yml":               filea,
			"buildspec-dockerize.yml":           fileb,
			"buildspec-version.yml":             filec,
		},
	})
}
