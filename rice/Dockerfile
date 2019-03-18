FROM openjdk:11-jre

RUN mkdir /app
ADD target/*.jar /app/app.jar

CMD ["/usr/bin/java", "-jar", "/app/app.jar"]
