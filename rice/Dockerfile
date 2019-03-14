FROM openjdk:11-jre

RUN mkdir /app
ADD target/*.jar /app

CMD ["/usr/bin/java", "-jar", "/app/*.jar"]
