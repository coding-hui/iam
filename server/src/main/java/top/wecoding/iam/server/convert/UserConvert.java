package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.pojo.UserInfo;
import top.wecoding.iam.server.pojo.User;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
@Mapper
public interface UserConvert {

  UserConvert INSTANCE = Mappers.getMapper(UserConvert.class);

  default UserInfoResponse toUserInfoResponse(User user, String tenantName) {
    if (user == null && tenantName == null) {
      return null;
    }

    UserInfoResponse.UserInfoResponseBuilder userInfoResponse = UserInfoResponse.builder();

    userInfoResponse.userInfo(toUserInfo(user, tenantName));

    return userInfoResponse.build();
  }

  @Mappings(@Mapping(source = "user", target = "userInfo"))
  UserInfoResponse toUserInfoResponse(User user);

  UserInfo toUserInfo(User user);

  UserInfo toUserInfo(User user, String tenantName);

  User toUser(CreateUserRequest createUserRequest);

  User toUser(UpdateUserRequest updateUserRequest);
}
