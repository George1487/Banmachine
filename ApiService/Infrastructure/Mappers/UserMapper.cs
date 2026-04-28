using Domain.Users;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class UserMapper
{
    public static User ToDomain(UserEntity entity)
    {
        return new User(
            entity.Id,
            entity.Email,
            entity.FullName,
            entity.Password,
            entity.GroupName,
            entity.Role == "teacher" ? UserRole.Teacher : UserRole.Student);
    }

    public static UserEntity ToEntity(User domain)
    {
        return new UserEntity
        {
            Id = domain.UserId,
            Email = domain.Email,
            FullName = domain.FullName,
            Password = domain.Password,
            GroupName = domain.GroupName,
            Role = domain.Role == UserRole.Teacher ? "teacher" : "student"
        };
    }
}
