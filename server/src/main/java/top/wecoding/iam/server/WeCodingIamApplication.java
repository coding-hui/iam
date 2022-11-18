package top.wecoding.iam.server;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import top.wecoding.feign.annotation.EnableWeCodingFeignClients;

@EnableWeCodingFeignClients
@SpringBootApplication(scanBasePackages = "top.wecoding.iam")
public class WeCodingIamApplication {

  public static void main(String[] args) {
    SpringApplication.run(WeCodingIamApplication.class, args);
  }
}
