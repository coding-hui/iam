package top.wecoding.iam.server.convert;

import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;
import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.UpdateUserRequest;
import top.wecoding.iam.server.entity.UserProfile;

/**
 * @author liuyuhui
 * @date 2022/11/13
 * @qq 1515418211
 */
@Mapper
public interface UserProfileConvert {

  UserProfileConvert INSTANCE = Mappers.getMapper(UserProfileConvert.class);

  UserProfile toUserProfile(CreateUserRequest createUserRequest);

  UserProfile toUserProfile(UpdateUserRequest updateUserRequest);
}
