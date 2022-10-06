package top.wecoding.iam.common.model.request;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public interface PasswordRequest {

  boolean reset();

  String getOldPwd();

  String getNewPwd();
}
