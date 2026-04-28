using Domain.Users;
using Ports.InBound.DTO;

namespace Ports.OutBound.Contracts;

public interface IUserRepository
{
    
    UserResult AddUser(User user);
    
    UserResult GetUser(Guid userId);
    
    UserResult GetUserByEmail(string email);
}