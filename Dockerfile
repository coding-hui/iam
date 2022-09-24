FROM java:8
EXPOSE 80
MAINTAINER wecoding@yeah.net
WORKDIR /wecoding
ADD ./jar/wecoding.jar ./wecoding.jar
ENTRYPOINT ["java", "-Djava.security.egd=file:/dev/./urandom", "-jar", "wecoding.jar"]
CMD ["--spring.profiles.active=prod"]
