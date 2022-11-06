package top.wecoding.iam.server.util;

import cn.hutool.extra.spring.SpringUtil;
import top.wecoding.core.util.WebUtils;
import top.wecoding.iam.server.mapper.UserMapper;

/**
 * @author liuyuhui
 * @date 2022/11/5
 * @qq 1515418211
 */
public class LogUtil {

  private static final UserMapper userMapper = SpringUtil.getBean(UserMapper.class);

  public static void successLogin(String userId) {
    userMapper.flushLastLoginInfo(userId, WebUtils.getIP());
  }
}
