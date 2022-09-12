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
@TableName("iam_group")
@EqualsAndHashCode(callSuper = true)
public class Group extends LogicDeletedBaseEntity {

  @TableId private Long id;

  private String groupId;

  private String tenantId;

  private String groupName;
}
