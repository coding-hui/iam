package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import top.wecoding.mybatis.base.BaseEntity;

/**
 * @author liuyuhui
 * @date 2022/9/12
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_user_group")
@EqualsAndHashCode(callSuper = true)
public class UserGroup extends BaseEntity {

  @TableId private Long id;

  private String tenantId;

  private String userId;

  private String groupId;
}
