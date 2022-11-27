package top.wecoding.iam.server;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.util.Properties;

@Slf4j
@SpringBootApplication(scanBasePackages = "top.wecoding.iam")
public class WeCodingIamApplication {

  public static void main(String[] args) {
    Properties properties = System.getProperties();
    properties.setProperty("spring.cloud.nacos.config.server-addr", "http://localhost:8848");
    SpringApplication.run(WeCodingIamApplication.class, args);
    log.info("{} started successfully.", WeCodingIamApplication.class.getSimpleName());
  }
}
