package top.wecoding.iam.framework.feign;

import cn.hutool.core.collection.CollUtil;
import cn.hutool.core.util.StrUtil;
import feign.RequestInterceptor;
import feign.RequestTemplate;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpHeaders;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;
import top.wecoding.core.util.WebUtils;
import top.wecoding.iam.common.constant.SecurityConstants;

import javax.servlet.http.HttpServletRequest;
import java.util.Collection;

/**
 * Feign 请求拦截器
 *
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
@RequiredArgsConstructor
public class IAMFeignRequestInterceptor implements RequestInterceptor {

  private final BearerTokenResolver tokenResolver;

  @Override
  public void apply(RequestTemplate template) {
    Collection<String> fromHeader = template.headers().get(SecurityConstants.FROM);
    if (CollUtil.isNotEmpty(fromHeader) && fromHeader.contains(SecurityConstants.INNER)) {
      return;
    }

    HttpServletRequest request = WebUtils.getRequest();
    if (null == request) {
      return;
    }
    // 传递用户信息请求头，防止丢失
    String token = tokenResolver.resolve(request);
    if (StrUtil.isBlank(token)) {
      return;
    }
    template.header(
        HttpHeaders.AUTHORIZATION,
        String.format("%s %s", OAuth2AccessToken.TokenType.BEARER.getValue(), token));
  }
}
