package top.wecoding.iam.pojo;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.*;
import top.wecoding.mybatis.base.LogicDeletedBaseEntity;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_user_group")
@EqualsAndHashCode(callSuper = true)
public class UserGroup extends LogicDeletedBaseEntity {

  @TableId private Long id;

  private String tenantId;

  private String userId;

  private String groupId;
}
