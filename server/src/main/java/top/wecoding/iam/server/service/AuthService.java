package top.wecoding.iam.server.service;

import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.common.model.response.AuthInfoResponse;
import top.wecoding.iam.common.model.response.CommonLoginResponse;
import top.wecoding.iam.common.model.response.UserInfoResponse;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface AuthService {

  AuthInfoResponse authInfo(AuthInfo authInfo);

  UserInfoResponse userInfo(AuthInfo authInfo);

  void logout(AuthInfo authInfo);

  CommonLoginResponse login(LoginRequest loginRequest);
}
