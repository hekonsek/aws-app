package awsom

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "Dockerfile",
		FileModTime: time.Unix(1552921724, 0),

		Content: string("FROM openjdk:11-jre\n\nRUN mkdir /app\nADD target/*.jar /app/app.jar\n\nCMD [\"/usr/bin/java\", \"-jar\", \"/app/app.jar\"]\n"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "assume_service_role_template.json",
		FileModTime: time.Unix(1551990782, 0),

		Content: string("{\n  \"Version\":\"2012-10-17\",\n  \"Statement\":[\n    {\n      \"Effect\":\"Allow\",\n      \"Principal\":{\n        \"Service\":\"{{ .Service }}\"\n      },\n      \"Action\":\"sts:AssumeRole\"\n    }\n  ]\n}"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "buildspec-build.yml",
		FileModTime: time.Unix(1552337535, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - mvn test package\nartifacts:\n  files:\n    - \"**/*\""),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "buildspec-configure.yml",
		FileModTime: time.Unix(1552569626, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - cp /usr/bin/awsom awsom\n      - ./awsom app configure\nartifacts:\n  files:\n    - \"**/*\""),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "buildspec-dockerize.yml",
		FileModTime: time.Unix(1552549826, 0),

		Content: string("version: 0.2\n\nphases:\n  pre_build:\n    commands:\n      - $(aws ecr get-login --region $AWS_REGION --no-include-email)\n      - REPOSITORY=`./awsom step ecr`\n      - VERSION=`./awsom step version-current`\n  build:\n    commands:\n      - docker build -t $REPOSITORY:$VERSION .\n      - docker push $REPOSITORY:$VERSION\nartifacts:\n  files:\n    - imageDetail.json"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "buildspec-version.yml",
		FileModTime: time.Unix(1552569626, 0),

		Content: string("version: 0.2\n\nphases:\n  build:\n    commands:\n      - ./awsom step clone\n      - ./awsom step version\nartifacts:\n  files:\n    - \"**/*\""),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1552921725, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "Dockerfile"
			file3, // "assume_service_role_template.json"
			file4, // "buildspec-build.yml"
			file5, // "buildspec-configure.yml"
			file6, // "buildspec-dockerize.yml"
			file7, // "buildspec-version.yml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`rice`, &embedded.EmbeddedBox{
		Name: `rice`,
		Time: time.Unix(1552921725, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"Dockerfile":                        file2,
			"assume_service_role_template.json": file3,
			"buildspec-build.yml":               file4,
			"buildspec-configure.yml":           file5,
			"buildspec-dockerize.yml":           file6,
			"buildspec-version.yml":             file7,
		},
	})
}
