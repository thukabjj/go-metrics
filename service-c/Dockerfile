FROM maven:3.8.3-openjdk-17 AS build
COPY  . /home/maven/src
WORKDIR /home/maven/src
RUN mvn clean install


FROM openjdk:17.0.2-slim

EXPOSE 8080

RUN mkdir /app

WORKDIR /app

RUN ls

COPY --from=build /home/maven/src/target/service-c-0.0.1-SNAPSHOT.jar ./spring-boot-application.jar
RUN ls

ENTRYPOINT ["java","-jar","spring-boot-application.jar"]