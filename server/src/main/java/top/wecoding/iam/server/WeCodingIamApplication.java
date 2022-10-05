package top.wecoding.iam.server;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.openfeign.EnableFeignClients;
import top.wecoding.iam.sdk.EnableIAMResourceServer;

@EnableFeignClients(basePackages = "top.wecoding")
@SpringBootApplication
@EnableIAMResourceServer
public class WeCodingIamApplication {

  public static void main(String[] args) {
    SpringApplication.run(WeCodingIamApplication.class, args);
  }
}
