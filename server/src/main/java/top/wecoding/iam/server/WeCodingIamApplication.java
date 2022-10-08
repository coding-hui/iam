package top.wecoding.iam.server;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import top.wecoding.feign.annotation.EnableWeCodingFeignClients;
import top.wecoding.iam.sdk.EnableIAMResourceServer;

@SpringBootApplication
@EnableIAMResourceServer
@EnableWeCodingFeignClients
public class WeCodingIamApplication {

  public static void main(String[] args) {
    SpringApplication.run(WeCodingIamApplication.class, args);
  }
}
