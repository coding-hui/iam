package top.wecoding.iam.sdk.props;

import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;

import java.util.ArrayList;
import java.util.List;

/**
 * Security 配置
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@Data
@ConfigurationProperties("wecoding.security")
public class SecurityProperties {

  private final List<ClientSecurity> client = new ArrayList<>();

  private List<String> whites = new ArrayList<>();
}
