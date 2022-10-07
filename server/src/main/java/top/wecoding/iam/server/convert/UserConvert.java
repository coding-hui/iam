package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.UserInfo;
import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.server.pojo.User;

import java.util.List;
import java.util.Set;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
 */
@Mapper
public interface UserConvert {

  UserConvert INSTANCE = Mappers.getMapper(UserConvert.class);

  @Mappings(@Mapping(source = "user", target = "userInfo"))
  UserInfoResponse toUserInfoResponse(User user);

  @Mappings({@Mapping(source = "user", target = "userInfo")})
  UserInfoResponse toUserInfoResponse(User user, String tenantName);

  @Mappings({
    @Mapping(source = "user", target = "userInfo"),
    @Mapping(source = "groupInfoList", target = "groupInfoList"),
    @Mapping(source = "permissions", target = "permissions"),
    @Mapping(source = "roles", target = "roles")
  })
  UserInfoResponse toUserInfoResponse(
      User user, List<GroupInfo> groupInfoList, Set<String> permissions, Set<String> roles);

  UserInfo toUserInfo(User user);

  UserInfo toUserInfo(User user, String tenantName);

  User toUser(CreateUserRequest createUserRequest);

  User toUser(UpdateUserRequest updateUserRequest);
}
