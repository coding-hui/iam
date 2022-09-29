package top.wecoding.iam.sdk.filter;

import cn.hutool.core.util.StrUtil;
import cn.hutool.json.JSONUtil;
import lombok.extern.slf4j.Slf4j;
import lombok.var;
import org.springframework.web.filter.OncePerRequestFilter;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.iam.sdk.context.AuthContextHolder;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Slf4j
public class AuthenticationProcessingFilter extends OncePerRequestFilter {

  @Override
  protected void doFilterInternal(
      HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
      throws ServletException, IOException {

    try {
      var authInfo = this.getAuthInfo(request);
      if (authInfo != null) {
        AuthContextHolder.setContext(authInfo);
      }
      filterChain.doFilter(request, response);
    } finally {
      AuthContextHolder.clearContext();
    }
  }

  private AuthInfo getAuthInfo(HttpServletRequest request) {
    String authorization = request.getHeader(TokenConstant.AUTHENTICATION);
    if (StrUtil.isNotBlank(authorization)) {
      log.debug("header authorization not exist");
      return null;
    }
    try {
      authorization =
          new String(
              Base64.getDecoder().decode(authorization.getBytes(StandardCharsets.UTF_8)),
              StandardCharsets.UTF_8);
      return JSONUtil.toBean(authorization, AuthInfo.class);
    } catch (Exception e) {
      return null;
    }
  }
}
