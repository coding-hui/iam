FROM openjdk:20-ea-17-jdk
EXPOSE 80
MAINTAINER wecoding@yeah.net
WORKDIR /wecoding
ADD ./target/iam-server.jar ./iam-server.jar
COPY ./target/jars ./jars/
ENTRYPOINT ["java", "-Dloader.path=/wecoding/jars", "-Dspring.config.additional-location=/etc/wecoding-iam/config.yaml", "-Djava.security.egd=file:/dev/./urandom", "-jar", "iam-server.jar"]
CMD ["--spring.profiles.active=prod"]
