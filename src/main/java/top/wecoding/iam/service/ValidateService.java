package top.wecoding.iam.service;

import top.wecoding.core.result.R;

/**
 * @author liuyuhui
 * @qq 1515418211
 */
public interface ValidateService {

  /**
   * 获取验证码
   *
   * @return Result
   */
  R<?> createCode();

  /**
   * 获取短信验证码
   *
   * @param mobile 手机号码
   * @return Result
   */
  R<?> createSmsCode(String mobile);

  /**
   * 校验验证码
   *
   * @param key 　key
   * @param code 验证码
   */
  void check(String key, String code);
}
