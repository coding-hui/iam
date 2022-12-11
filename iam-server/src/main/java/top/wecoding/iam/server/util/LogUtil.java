package top.wecoding.iam.server.util;

import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.stereotype.Component;
import top.wecoding.commons.core.util.IpUtil;
import top.wecoding.iam.server.mapper.UserMapper;

/**
 * @author liuyuhui
 * @date 2022/11/5
 */
@Component
public class LogUtil implements ApplicationContextAware {

  private static UserMapper userMapper;

  private static ApplicationContext context;

  public static void successLogin(String userId) {
    getUserMapper().flushLastLoginInfo(userId, IpUtil.getIp());
  }

  public static UserMapper getUserMapper() {
    if (userMapper == null) {
      userMapper = context.getBean(UserMapper.class);
    }
    return userMapper;
  }

  @Override
  public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
    LogUtil.context = applicationContext;
  }
}
