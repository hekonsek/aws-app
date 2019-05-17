package main

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file9 := &embedded.EmbeddedFile{
		Filename:    "Dockerfile",
		FileModTime: time.Unix(1552921724, 0),

		Content: string("FROM openjdk:11-jre\n\nRUN mkdir /app\nADD target/*.jar /app/app.jar\n\nCMD [\"/usr/bin/java\", \"-jar\", \"/app/app.jar\"]\n"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "assume_service_role_template.json",
		FileModTime: time.Unix(1551990782, 0),

		Content: string("{\n  \"Version\":\"2012-10-17\",\n  \"Statement\":[\n    {\n      \"Effect\":\"Allow\",\n      \"Principal\":{\n        \"Service\":\"{{ .Service }}\"\n      },\n      \"Action\":\"sts:AssumeRole\"\n    }\n  ]\n}"),
	}
	fileb := &embedded.EmbeddedFile{
		Filename:    "buildspec-build.yml",
		FileModTime: time.Unix(1552337535, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - mvn test package\nartifacts:\n  files:\n    - \"**/*\""),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "buildspec-configure.yml",
		FileModTime: time.Unix(1552569626, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - cp /usr/bin/awsom awsom\n      - ./awsom app configure\nartifacts:\n  files:\n    - \"**/*\""),
	}
	filed := &embedded.EmbeddedFile{
		Filename:    "buildspec-dockerize.yml",
		FileModTime: time.Unix(1552549826, 0),

		Content: string("version: 0.2\n\nphases:\n  pre_build:\n    commands:\n      - $(aws ecr get-login --region $AWS_REGION --no-include-email)\n      - REPOSITORY=`./awsom step ecr`\n      - VERSION=`./awsom step version-current`\n  build:\n    commands:\n      - docker build -t $REPOSITORY:$VERSION .\n      - docker push $REPOSITORY:$VERSION\nartifacts:\n  files:\n    - imageDetail.json"),
	}
	filee := &embedded.EmbeddedFile{
		Filename:    "buildspec-version.yml",
		FileModTime: time.Unix(1552569626, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - ./awsom step clone\n      - ./awsom step version\nartifacts:\n  files:\n    - \"**/*\""),
	}

	// define dirs
	dir8 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1552921725, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file9, // "Dockerfile"
			filea, // "assume_service_role_template.json"
			fileb, // "buildspec-build.yml"
			filec, // "buildspec-configure.yml"
			filed, // "buildspec-dockerize.yml"
			filee, // "buildspec-version.yml"

		},
	}

	// link ChildDirs
	dir8.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`../rice`, &embedded.EmbeddedBox{
		Name: `../rice`,
		Time: time.Unix(1552921725, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir8,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"Dockerfile":                        file9,
			"assume_service_role_template.json": filea,
			"buildspec-build.yml":               fileb,
			"buildspec-configure.yml":           filec,
			"buildspec-dockerize.yml":           filed,
			"buildspec-version.yml":             filee,
		},
	})
}
