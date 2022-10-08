FROM openjdk:8-jre
EXPOSE 80
MAINTAINER wecoding@yeah.net
WORKDIR /wecoding
ADD ../target/wecoding-iam.jar ./wecoding-iam.jar
COPY ../target/jars ./jars/
ENTRYPOINT ["java", "-Dloader.path=/app/jars", "-Djava.security.egd=file:/dev/./urandom", "-jar", "wecoding-iam.jar"]
CMD ["--spring.profiles.active=prod"]
