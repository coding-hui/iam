package top.wecoding.iam.model.response;

import com.fasterxml.jackson.annotation.JsonUnwrapped;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import top.wecoding.iam.model.GroupInfo;
import top.wecoding.iam.model.UserInfo;

import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class UserInfoResponse {

  @JsonUnwrapped private UserInfo userInfo;

  @JsonUnwrapped private List<GroupInfo> groupInfoList;
}
