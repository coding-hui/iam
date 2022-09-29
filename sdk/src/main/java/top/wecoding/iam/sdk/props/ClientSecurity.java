package top.wecoding.iam.sdk.props;

import lombok.Data;

import java.util.ArrayList;
import java.util.List;

/**
 * 客户端令牌认证信息
 *
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
@Data
public class ClientSecurity {

  private final List<String> pathPatterns = new ArrayList<>();

  private String clientId;
}
