package top.wecoding.iam.server;

import java.util.Properties;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import top.wecoding.feign.annotation.EnableWeCodingFeignClients;
import top.wecoding.iam.sdk.EnableIAMResourceServer;

@SpringBootApplication
@EnableIAMResourceServer
@EnableWeCodingFeignClients
public class WeCodingIamApplication {

  public static void main(String[] args) {
    Properties properties = System.getProperties();
    properties.setProperty("spring.cloud.nacos.config.server-addr", "http://wecoding.top:8848");
    SpringApplication.run(WeCodingIamApplication.class, args);
  }
}
