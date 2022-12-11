package top.wecoding.iam.framework.security.handler;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationListener;
import org.springframework.security.authentication.event.LogoutSuccessEvent;
import org.springframework.security.core.Authentication;
import org.springframework.security.web.authentication.preauth.PreAuthenticatedAuthenticationToken;
import org.springframework.stereotype.Component;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
@Component
public class WeCodingLogoutSuccessEventHandler implements ApplicationListener<LogoutSuccessEvent> {

  @Override
  public void onApplicationEvent(LogoutSuccessEvent event) {
    Authentication authentication = (Authentication) event.getSource();
    if (authentication instanceof PreAuthenticatedAuthenticationToken) {
      handle(authentication);
    }
  }

  /**
   * 处理退出成功方法
   *
   * <p>获取到登录的authentication 对象
   *
   * @param authentication 登录对象
   */
  public void handle(Authentication authentication) {
    // TODO 发送退出日志
    log.info("user: {} logout success", authentication.getPrincipal());
  }
}
